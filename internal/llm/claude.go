// Package llm wraps the Anthropic Messages API. We call the REST endpoint
// directly to avoid pulling in an SDK for what amounts to two endpoints
// (parse and dossier).
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	apiBase    = "https://api.anthropic.com"
	apiVersion = "2023-06-01"

	ModelHaiku  = "claude-haiku-4-5-20251001"
	ModelSonnet = "claude-sonnet-4-6"
)

// Client is safe for concurrent use.
type Client struct {
	apiKey string
	http   *http.Client
}

// New returns a Client, or nil if apiKey is empty (so callers can keep the
// LLM-dependent endpoints behind a "not configured" 503 in dev).
func New(apiKey string) *Client {
	if strings.TrimSpace(apiKey) == "" {
		return nil
	}
	return &Client{
		apiKey: apiKey,
		// Upper-bound timeout. Per-call ctx tightens it (parse=30s,
		// dossier=150s) — dossier with web search legitimately takes a minute.
		http: &http.Client{Timeout: 180 * time.Second},
	}
}

// Message.Content is `any` to accept either a string (the common case) or a
// []ContentBlock (when sending an image alongside text). Anthropic's API
// accepts both shapes on the wire.
type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

// ContentBlock is one element of a multi-part message body. Use TextBlock /
// ImageBlock helpers to build these — the JSON tags here match Anthropic's
// expected wire format.
type ContentBlock struct {
	Type   string       `json:"type"`
	Text   string       `json:"text,omitempty"`
	Source *ImageSource `json:"source,omitempty"`
}

type ImageSource struct {
	Type      string `json:"type"`       // "base64"
	MediaType string `json:"media_type"` // image/png, image/jpeg, image/gif, image/webp
	Data      string `json:"data"`       // base64-encoded, no data: prefix
}

func TextBlock(s string) ContentBlock {
	return ContentBlock{Type: "text", Text: s}
}

func ImageBlock(mediaType, base64Data string) ContentBlock {
	return ContentBlock{
		Type:   "image",
		Source: &ImageSource{Type: "base64", MediaType: mediaType, Data: base64Data},
	}
}

// Tool is a (subset of) the Anthropic Messages API tool spec. We only need
// the server-side tools (web_search), which Anthropic executes for us — the
// final assistant turn already contains the synthesized text.
type Tool struct {
	Type    string `json:"type"`
	Name    string `json:"name,omitempty"`
	MaxUses int    `json:"max_uses,omitempty"`
}

func WebSearchTool(maxUses int) Tool {
	return Tool{Type: "web_search_20250305", Name: "web_search", MaxUses: maxUses}
}

type messagesRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system,omitempty"`
	Messages  []Message `json:"messages"`
	Tools     []Tool    `json:"tools,omitempty"`
}

type messagesResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
	Usage      struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// CreateMessage hits POST /v1/messages and returns the concatenated text
// content across all text blocks in the assistant response. Pass tools=nil for
// a vanilla call; pass a tool list (e.g. WebSearchTool) for server-side tools.
func (c *Client) CreateMessage(ctx context.Context, model, system string, messages []Message, maxTokens int, tools ...Tool) (string, error) {
	body, err := json.Marshal(messagesRequest{
		Model:     model,
		MaxTokens: maxTokens,
		System:    system,
		Messages:  messages,
		Tools:     tools,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", apiBase+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", apiVersion)
	req.Header.Set("content-type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic %d: %s", resp.StatusCode, string(raw))
	}
	var out messagesResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", fmt.Errorf("parse anthropic response: %w", err)
	}
	if len(out.Content) == 0 {
		return "", errors.New("anthropic returned empty content")
	}
	var text strings.Builder
	for _, b := range out.Content {
		if b.Type == "text" {
			text.WriteString(b.Text)
		}
	}
	return text.String(), nil
}

// ParsedJob is the structured result we expect Claude to extract from arbitrary
// pasted text. All fields are best-effort; omitted fields mean "not found".
type ParsedJob struct {
	Company    string `json:"company,omitempty"`
	Role       string `json:"role,omitempty"`
	Location   string `json:"location,omitempty"`
	Seniority  string `json:"seniority,omitempty"`
	JDURL      string `json:"jd_url,omitempty"`
	Source     string `json:"source,omitempty"`
	SalaryNote string `json:"salary_note,omitempty"`
}

const parseSystemPrompt = `You are a precise job-listing parser. The user will paste arbitrary text — a LinkedIn URL, a copied job description, an email forward, a recruiter message, or a description in their own words — OR attach a screenshot of a job page — and you must extract structured fields.

When a screenshot is attached, read the visible text in the image (company name, job title, location, posting source, salary range if shown) and treat it as the primary input. The page is often LinkedIn, Greenhouse, Lever, or a company careers page — set source accordingly. If a URL is visible in the screenshot (browser bar or share link), put it in jd_url.

Be PERMISSIVE: short notes count as long as you can extract a company AND a role. Examples of valid input you should parse, not reject:

  "Anthropic / Senior ML Engineer / referred by Sarah"
    → company=Anthropic, role=Senior ML Engineer, source=Referral

  "Stripe — Staff Backend Engineer · referred by Mia · applied 14 May"
    → company=Stripe, role=Staff Backend Engineer, source=Referral

  "talked to Vercel about a Frontend Eng role, JD: https://vercel.com/careers/eng-1234"
    → company=Vercel, role=Frontend Engineer, jd_url=https://vercel.com/careers/eng-1234

Return ONLY a JSON object with these keys (omit any you cannot determine with high confidence; do not invent values):
- company:    the hiring company name (e.g. "Anthropic")
- role:       the role title (e.g. "Senior Software Engineer")
- location:   location or remote status (e.g. "San Francisco", "Remote", "Remote (US)")
- seniority:  level if clear (e.g. "Senior", "Staff", "Principal", "Director")
- jd_url:     the canonical job-description URL if present
- source:     "LinkedIn" if from a LinkedIn URL/page; "Referral" if the input mentions being referred by someone; else infer (e.g. "Greenhouse", "Lever", "Company site", "Email") or omit
- salary_note: comp info if mentioned (e.g. "$220k-$280k base", "$160-200k OTE")

Output the JSON object only — no prose, no markdown fences.

Only return {"error": "not a job listing"} if the input genuinely isn't about a job — e.g. a recipe, a chat about the weather, a code snippet, an empty string, or something with no extractable company or role.`

// ParseImage is the optional screenshot/photo input for ParseJob. MediaType
// must be one of image/png, image/jpeg, image/gif, image/webp. Data is the
// raw base64 string (no `data:` prefix).
type ParseImage struct {
	MediaType string
	Data      string
}

// ParseJob asks Claude Haiku to extract structured job-listing fields from
// ParseError carries a machine-readable reason alongside the user-facing
// message, so analytics can tell *why* a parse failed instead of bucketing
// everything as "parse_failed".
type ParseError struct {
	Reason string // linkedin_url, url_fetch_failed, llm_error, no_json, not_a_job, model_error, bad_json, empty_extraction
	Msg    string
}

func (e *ParseError) Error() string { return e.Msg }

// arbitrary text and/or a screenshot. If text is a bare URL we fetch it
// server-side first (LinkedIn is rejected since they block scraping). When an
// image is supplied we send it as a vision content block alongside the text.
func (c *Client) ParseJob(ctx context.Context, text string, image *ParseImage) (*ParsedJob, error) {
	cctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	ctx = cctx
	text = strings.TrimSpace(text)

	// Bare-URL path: fetch the page server-side so Claude has something to read.
	// Only when there's no image — if the user attached a screenshot, the URL
	// is likely just the page they screenshotted, not the input we should fetch.
	if image == nil {
		if u, ok := isBareURL(text); ok {
			if isLinkedInURL(u) {
				return nil, &ParseError{"linkedin_url", "LinkedIn blocks page fetching — copy the JD body text from the LinkedIn page and paste that instead"}
			}
			body, err := c.fetchURL(ctx, u)
			if err != nil {
				return nil, &ParseError{"url_fetch_failed", fmt.Sprintf("couldn't fetch %s (%v) — try pasting the JD text directly", u, err)}
			}
			text = "Source URL: " + u + "\n\nPage content:\n" + body
		}
	}

	var msg Message
	if image != nil {
		// User caption defaults to a generic instruction when empty, so the model
		// always has at least one text block to anchor the image to.
		caption := text
		if caption == "" {
			caption = "Extract company, role, location, seniority, and any JD URL or source visible in this screenshot."
		} else {
			caption = "Extract from this screenshot. User added: " + caption
		}
		msg = Message{
			Role: "user",
			Content: []ContentBlock{
				ImageBlock(image.MediaType, image.Data),
				TextBlock(caption),
			},
		}
	} else {
		msg = Message{Role: "user", Content: text}
	}

	resp, err := c.CreateMessage(ctx, ModelHaiku, parseSystemPrompt, []Message{msg}, 600)
	if err != nil {
		return nil, &ParseError{"llm_error", err.Error()}
	}

	raw, err := extractFirstJSONObject(resp)
	if err != nil {
		return nil, &ParseError{"no_json", fmt.Sprintf("model didn't return JSON (raw: %.200s)", resp)}
	}

	// Explicit-error reply: surface a friendlier message.
	var errProbe struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(raw, &errProbe) == nil && errProbe.Error != "" {
		if strings.EqualFold(errProbe.Error, "not a job listing") {
			return nil, &ParseError{"not_a_job", "that doesn't look like a job listing — try pasting the JD body text instead"}
		}
		return nil, &ParseError{"model_error", errProbe.Error}
	}

	var job ParsedJob
	if err := json.Unmarshal(raw, &job); err != nil {
		return nil, &ParseError{"bad_json", fmt.Sprintf("parse model JSON: %v", err)}
	}
	if job.Company == "" && job.Role == "" {
		return nil, &ParseError{"empty_extraction", "could not extract company or role from the input"}
	}
	return &job, nil
}

var (
	bareURLRe = regexp.MustCompile(`^https?://\S+$`)
	scriptRe  = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	styleRe   = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	tagRe     = regexp.MustCompile(`<[^>]+>`)
	wsRe      = regexp.MustCompile(`\s+`)
)

// ParsedEvent is the structured calendar event we expect Claude to extract
// from a screenshot or pasted email body. Mirrors the .ics parser output so
// the frontend can render the preview through the same code path.
type ParsedEvent struct {
	Summary     string   `json:"summary,omitempty"`
	Location    string   `json:"location,omitempty"`
	Description string   `json:"description,omitempty"`
	StartsAt    string   `json:"starts_at,omitempty"` // ISO 8601
	EndsAt      string   `json:"ends_at,omitempty"`   // ISO 8601
	AllDay      bool     `json:"all_day,omitempty"`
	Attendees   []string `json:"attendees,omitempty"` // free-form names/emails
}

const eventParserSystemPrompt = `You are a precise calendar-event parser. The user pastes either a screenshot (Gmail invite, Google Calendar event detail, Outlook screenshot, etc.) OR raw text (email body, copied invite, plain notes). Extract the event details.

Return ONLY a JSON object with this shape (omit any you can't determine):

{
  "summary":     "Event title (e.g. 'Stripe — Technical screen')",
  "location":    "Conferencing link or physical address (e.g. 'Google Meet', 'Zoom: https://...', '510 Townsend St, SF')",
  "description": "Short description, agenda, or notes if any",
  "starts_at":   "ISO 8601 with timezone (e.g. '2026-05-28T14:00:00-04:00'). If only a date is given, use 00:00 in the user's likely timezone and set all_day: true.",
  "ends_at":     "ISO 8601 with timezone. Omit if not specified.",
  "all_day":     true|false,
  "attendees":   ["Free-form list of names or emails of attendees"]
}

Rules:
- summary is required if any event is described; if you can't find a summary, fall back to the company / interview type (e.g. "Interview").
- Default timezone: if the screenshot/text shows a timezone (PDT, EST, "Israel Time", "America/New_York"), use it. Otherwise use the user's timezone given in the message; only if none is given, fall back to US Eastern.
- Date resolution: resolve dates relative to the "Current date" given in the user message. A year-less or relative date ("next Wednesday", "10 June", "tomorrow") means the SOONEST such date on or after the current date.
- Weekday consistency: if the input names a weekday (e.g. "Wednesday, 10 June"), the resolved date MUST fall on that weekday. A bare month/day can land on different weekdays in different years — pick the year that makes the named weekday correct. Never emit a date whose weekday contradicts the text.
- Times like "2pm" or "14:00": assume the local timezone you inferred.
- If the input is a series (multiple events), pick the SOONEST one — only return a single event.

Output the JSON object only — no prose, no markdown. If the input has no event you can extract, return: {"error": "no event found"}.`

// ParseEvent asks Claude to extract a single calendar event from a screenshot
// and/or pasted text. Same Haiku + Vision flow as ParseJob. Returns nil with
// an error if nothing extractable is found.
func (c *Client) ParseEvent(ctx context.Context, text string, image *ParseImage, tz string) (*ParsedEvent, error) {
	cctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	ctx = cctx
	text = strings.TrimSpace(text)

	if text == "" && image == nil {
		return nil, errors.New("paste an event screenshot or text first")
	}

	// Anchor relative/year-less dates and weekday checks to "now", and do it in
	// the user's own timezone so a bare "2:30pm" resolves to their wall clock —
	// not the server's tz, and not a hardcoded US Eastern default. Without this
	// an Israeli user typing "2:30pm" gets it read as US Eastern and shifted.
	loc := time.Local
	if tz != "" {
		if l, err := time.LoadLocation(tz); err == nil {
			loc = l
		}
	}
	now := time.Now().In(loc)
	dateAnchor := "Current date: " + now.Format("Monday, 2006-01-02 -07:00") + "."
	if tz != "" {
		dateAnchor += " User's timezone: " + tz + ". Interpret any time without an explicit timezone in the user's timezone."
	}

	var msg Message
	if image != nil {
		caption := dateAnchor + "\n"
		if text == "" {
			caption += "Extract the calendar event: title, time, location, attendees."
		} else {
			caption += "Extract the calendar event from this screenshot. User added: " + text
		}
		msg = Message{
			Role: "user",
			Content: []ContentBlock{
				ImageBlock(image.MediaType, image.Data),
				TextBlock(caption),
			},
		}
	} else {
		msg = Message{Role: "user", Content: dateAnchor + "\n\n" + text}
	}

	resp, err := c.CreateMessage(ctx, ModelHaiku, eventParserSystemPrompt, []Message{msg}, 600)
	if err != nil {
		return nil, err
	}
	raw, err := extractFirstJSONObject(resp)
	if err != nil {
		return nil, fmt.Errorf("model didn't return JSON (raw: %.300s)", resp)
	}
	var errProbe struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(raw, &errProbe) == nil && errProbe.Error != "" {
		return nil, errors.New(errProbe.Error)
	}
	var ev ParsedEvent
	if err := json.Unmarshal(raw, &ev); err != nil {
		return nil, fmt.Errorf("event JSON missing expected fields: %w", err)
	}
	if ev.Summary == "" && ev.StartsAt == "" {
		return nil, errors.New("could not extract an event from the input")
	}
	return &ev, nil
}

func isBareURL(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if !bareURLRe.MatchString(s) || len(s) > 500 {
		return "", false
	}
	return s, true
}

func isLinkedInURL(u string) bool {
	lower := strings.ToLower(u)
	return strings.Contains(lower, "linkedin.com/")
}

// fetchURL retrieves a public URL and returns visible text content (HTML tags,
// scripts, and styles stripped). Capped to ~30k chars so we don't blow the
// token budget on huge pages.
func (c *Client) fetchURL(ctx context.Context, url string) (string, error) {
	cctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(cctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Pursuit/0.1)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 400_000))
	if err != nil {
		return "", err
	}
	text := stripHTML(string(raw))
	if len(text) > 30_000 {
		text = text[:30_000]
	}
	return text, nil
}

func stripHTML(s string) string {
	s = scriptRe.ReplaceAllString(s, " ")
	s = styleRe.ReplaceAllString(s, " ")
	s = tagRe.ReplaceAllString(s, " ")
	s = wsRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// ─────────────────────────────────────────────────────────────────────────
// Dossier generation (Claude Sonnet + web search)
// ─────────────────────────────────────────────────────────────────────────

const dossierSystemPrompt = `You are an interview-prep researcher with web search. Given a company, role, and (optionally) an interviewer's name, produce a structured briefing the candidate can read in 60 seconds before the interview.

You MUST use web search to ground every claim. Search for things like:
- The interviewer's recent essays, talks, podcasts, papers (last 12 months)
- The company's recent news, product launches, leadership statements
- The team's public engineering or design culture
- The role's likely scope given company stage and recent direction

Do NOT invent quotes, papers, or events. If you can't find something concrete, omit that field — empty arrays and empty strings are fine.

Return ONLY a JSON object with this exact shape (no prose, no markdown fences):

{
  "interviewer": {
    "initials": "DA",
    "name": "Dario Amodei",
    "role": "CEO & Co-founder, Anthropic",
    "prior": ["VP Research, OpenAI", "Princeton PhD (computational neuroscience)"],
    "links": [
      {"label": "LinkedIn", "href": "https://linkedin.com/..."},
      {"label": "Essay: Machines of Loving Grace", "href": "https://..."}
    ]
  },
  "snapshot": "A 2-sentence read of who this person is and how they interview. Use <em>...</em> tags around 2-3 key phrases for emphasis.",
  "background": "60–80 words. Their professional arc, framed so the candidate can engage thoughtfully.",
  "signals": [
    {"date": "Apr 2026", "kind": "Essay", "body": "What it was about and why it matters for this interview.", "source": "domain.com"}
  ],
  "style": {
    "lead": "1–2 sentences on how they interview. What are they testing for?",
    "tells": [
      {"lbl": "Format", "val": "..."},
      {"lbl": "Tempo",  "val": "..."},
      {"lbl": "Length", "val": "..."}
    ]
  },
  "lands":     ["4 short, specific things that work with this person."],
  "avoid":     ["4 short, specific things that don't."],
  "questions": [
    {"q": "Question for the candidate to ask back.", "why": "Why it lands."}
  ],
  "company": {
    "blurb":     "One sentence on what the company does. Plain English, no marketing.",
    "direction": "Two sentences on where they're going right now. Recent news, last raise, product launches, leadership shifts. Last 12 months only.",
    "hq":        "City",
    "employees": "~N,NNN (range or approximate is fine)",
    "stage":     "Late stage · $X valuation  /  Series B  /  Public  /  Bootstrapped — pick what fits",
    "founded":   "YYYY",
    "process": [
      {"kind": "Recruiter screen",   "detail": "30 min · culture + role fit"},
      {"kind": "Hiring-manager call","detail": "45 min · architecture, prior wins"},
      {"kind": "Technical screen",   "detail": "60 min · live coding (language)"},
      {"kind": "Onsite loop",        "detail": "4 panels · sys design, IC code, behavioral, bar-raiser"},
      {"kind": "Team match",         "detail": "Two 30-min chats with prospective teams"}
    ],
    "watch_fors": [
      "Five short sentences specific to THIS company's loop.",
      "Drawn from Glassdoor / Blind / levels.fyi / engineering blog content.",
      "Not generic advice — name what THIS team grades for.",
      "Mention rollout / observability / failure modes if relevant.",
      "Mention the company's current strategic bet if it matters."
    ]
  }
}

The "company" block is REQUIRED — populate it for every dossier regardless of whether an interviewer is named. Use web search to ground every claim. Omit fields you can't verify (empty string is fine, but DON'T invent numbers, valuations, or interview steps you can't source).

If the interviewer name is empty, ALSO write a COMPANY-LEVEL interviewer-section:
- interviewer.name = "Hiring team"
- interviewer.role = the company name
- interviewer.initials = the first letter of the company
- interviewer.prior = company highlights or recent moves
- snapshot/background/signals scoped to the company itself
- style/lands/avoid/questions reflect the company's interview culture

If after searching you genuinely can't find enough information, return: {"error": "could not find enough public information about this person/company"}.

Output JSON only.`

// GenerateDossier asks Claude Sonnet (with web search) to research the
// interviewer/company and produce a structured briefing. Returns the JSON
// bytes verbatim so the frontend can render exactly what the model decided
// to ship — no field-level translation in Go.
func (c *Client) GenerateDossier(ctx context.Context, company, role, interviewerName string) (json.RawMessage, error) {
	cctx, cancel := context.WithTimeout(ctx, 150*time.Second)
	defer cancel()

	var user strings.Builder
	fmt.Fprintf(&user, "Company: %s\nRole: %s\n", company, role)
	if strings.TrimSpace(interviewerName) != "" {
		fmt.Fprintf(&user, "Interviewer: %s\n", strings.TrimSpace(interviewerName))
	} else {
		user.WriteString("Interviewer: (not specified — produce a company-level briefing)\n")
	}

	resp, err := c.CreateMessage(cctx, ModelSonnet, dossierSystemPrompt, []Message{
		{Role: "user", Content: user.String()},
	}, 4000, WebSearchTool(5))
	if err != nil {
		return nil, err
	}

	raw, err := extractFirstJSONObject(resp)
	if err != nil {
		return nil, fmt.Errorf("model didn't return JSON (raw: %.300s)", resp)
	}

	var errProbe struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(raw, &errProbe) == nil && errProbe.Error != "" {
		return nil, errors.New(errProbe.Error)
	}

	// Light sanity check — confirm we got the shape we expect.
	var shape struct {
		Interviewer struct {
			Name string `json:"name"`
		} `json:"interviewer"`
		Snapshot string `json:"snapshot"`
	}
	if err := json.Unmarshal(raw, &shape); err != nil {
		return nil, fmt.Errorf("dossier JSON missing expected fields: %w", err)
	}
	if shape.Interviewer.Name == "" && shape.Snapshot == "" {
		return nil, errors.New("dossier came back empty")
	}
	return raw, nil
}

// extractFirstJSONObject finds the first top-level JSON object in s and returns
// it as bytes, ignoring any prose, markdown fences, or trailing text the model
// may have wrapped it in.
func extractFirstJSONObject(s string) ([]byte, error) {
	start := strings.Index(s, "{")
	if start == -1 {
		return nil, errors.New("no JSON object found")
	}
	dec := json.NewDecoder(strings.NewReader(s[start:]))
	var raw json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}
