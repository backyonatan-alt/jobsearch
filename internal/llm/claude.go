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

// Company brief — shared across every interview round for an application.
// Researched once; the round briefs reference it for grounding but don't repeat
// the company block.
const companyBriefSystemPrompt = `You are an interview-prep researcher with web search. Given a company and a role, produce a COMPANY briefing the candidate can read in 60 seconds — true for every round of their interview loop, independent of who interviews them.

CRITICAL — identify the RIGHT company first. Company names are often shared by several unrelated organizations. The user message may include a LOCATION, a JOB-POSTING URL, and/or an AUTHORITATIVE company website — when present, treat these as the GROUND TRUTH for which company to research, even if a more famous company shares the name. If you can't tell which is meant, research the one that best matches the provided location / website / job posting — never default to a better-known same-named company. The "identity" block of your answer MUST state exactly which company (canonical name + primary website domain) you researched, so the candidate can verify you got the right one.

You MUST use web search to ground every claim (recent news, raises, launches, leadership, the team's public engineering/design culture, Glassdoor/Blind/levels.fyi on the loop). Do NOT invent numbers, valuations, or interview steps you can't source — omit what you can't verify (empty string/array is fine). For each non-obvious claim, capture the SPECIFIC source page URL (a deep link to the article/review/page — never just the site homepage) in "sources".

Return ONLY a JSON object with this exact shape (no prose, no markdown fences):

{
  "identity": {
    "name":    "Canonical company name you researched",
    "domain":  "primary website domain, e.g. 365scores.com",
    "summary": "One short line — what they are and where they're based — so the candidate can confirm it's the right company."
  },
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
  },
  "sources": [
    {"label": "TechCrunch — $40M Series C (Mar 2026)", "href": "https://techcrunch.com/2026/03/.../"},
    {"label": "Glassdoor — interview reviews", "href": "https://www.glassdoor.com/Interview/..."}
  ]
}

Every "href" in "sources" must be a real, specific page you actually used — a deep link, not a homepage. Omit "sources" entirely rather than inventing or linking to homepages.

If after searching you genuinely can't find enough information, return: {"error": "could not find enough public information about this company"}.

Output JSON only.`

// Interviewer brief — specific to one round / one interviewer. The company is
// passed only for grounding; do NOT emit a company block here (it lives in the
// shared company brief).
const interviewerBriefSystemPrompt = `You are an interview-prep researcher with web search. Given a company, role, and an interviewer's name, produce a briefing about THAT PERSON the candidate can read in 60 seconds before the round.

You MUST use web search to ground every claim. Search for the interviewer's recent essays, talks, podcasts, papers (last 12 months), their professional arc, and how they're known to interview. Do NOT invent quotes, papers, or events — omit what you can't find (empty string/array is fine). For every signal, "source_url" must be the SPECIFIC page the claim came from (a deep link to that essay/talk/post — never the site homepage); "source" is just the short display domain. Omit a signal rather than linking to a homepage.

The user message may include a "This round" line describing the round's FORMAT (e.g. a home-assignment presentation, a system-design session, a behavioral loop, a take-home review, an HR screen). When present, tailor the WHOLE brief to that format, not just to the person: "style" describes how this person runs THAT kind of round; "lands"/"avoid" become format-specific (for a home-assignment presentation: defending choices and tradeoffs, structuring the walkthrough, handling "why not X?" drills — not generic rapport tips); "questions" fit the format. If the format is unclear from the text, prep for it as a conversation with this person and say nothing about format.

Do NOT include a company overview — that's covered separately. Stay focused on the person. The user message may include a company location or authoritative website to disambiguate same-named companies — use it to make sure you're researching this person at the RIGHT company.

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
    {"date": "Apr 2026", "kind": "Essay", "body": "What it was about and why it matters for this interview.", "source": "domain.com", "source_url": "https://domain.com/the-specific-post"}
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
  ]
}

If the interviewer is a generic name you can't find online, still return a useful briefing scoped to their likely role at the company (snapshot/style/lands/avoid/questions), with prior/signals/links empty.

If after searching you genuinely can't find enough information, return: {"error": "could not find enough public information about this person"}.

Output JSON only.`

// companyGroundingPrompt assembles the user message with whatever disambiguation
// signals we have. A linkedin.com job URL is dropped — it identifies the posting,
// not the company, and is the exact source of the "wrong same-named company" bug.
func companyGroundingPrompt(company, role, location, jdURL, companyURL string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Company: %s\nRole: %s\n", company, role)
	if s := strings.TrimSpace(location); s != "" {
		fmt.Fprintf(&b, "Location: %s\n", s)
	}
	if s := strings.TrimSpace(jdURL); s != "" && !isLinkedInURL(s) {
		fmt.Fprintf(&b, "Job-posting URL: %s\n", s)
	}
	if s := strings.TrimSpace(companyURL); s != "" {
		fmt.Fprintf(&b, "Authoritative company website: %s (the candidate confirmed THIS is the company — trust it over any better-known company that shares the name)\n", s)
	}
	return b.String()
}

// GenerateCompanyBrief researches the company once for the whole application.
// location, jdURL, and companyURL are grounding signals that disambiguate
// same-named companies — any may be empty.
func (c *Client) GenerateCompanyBrief(ctx context.Context, company, role, location, jdURL, companyURL string) (json.RawMessage, error) {
	cctx, cancel := context.WithTimeout(ctx, 150*time.Second)
	defer cancel()

	user := companyGroundingPrompt(company, role, location, jdURL, companyURL)
	resp, err := c.CreateMessage(cctx, ModelSonnet, companyBriefSystemPrompt, []Message{
		{Role: "user", Content: user},
	}, 4000, WebSearchTool(5))
	if err != nil {
		return nil, err
	}
	raw, err := extractFirstJSONObject(resp)
	if err != nil {
		return nil, fmt.Errorf("model didn't return JSON (raw: %.300s)", resp)
	}
	if msg := jsonErrorField(raw); msg != "" {
		return nil, errors.New(msg)
	}
	var shape struct {
		Company struct {
			Blurb string `json:"blurb"`
		} `json:"company"`
	}
	if err := json.Unmarshal(raw, &shape); err != nil {
		return nil, fmt.Errorf("company brief JSON missing expected fields: %w", err)
	}
	if shape.Company.Blurb == "" {
		return nil, errors.New("company brief came back empty")
	}
	return raw, nil
}

// buildInterviewerPromptUser assembles the user message for a round's brief.
// Split out so tests can pin what the model actually sees.
func buildInterviewerPromptUser(company, role, interviewerName, location, companyURL, roundContext, priorDebriefs string) string {
	var user strings.Builder
	fmt.Fprintf(&user, "Company: %s\nRole: %s\n", company, role)
	if s := strings.TrimSpace(location); s != "" {
		fmt.Fprintf(&user, "Company location: %s\n", s)
	}
	if s := strings.TrimSpace(companyURL); s != "" {
		fmt.Fprintf(&user, "Authoritative company website: %s\n", s)
	}
	if strings.TrimSpace(interviewerName) != "" {
		fmt.Fprintf(&user, "Interviewer: %s\n", strings.TrimSpace(interviewerName))
	} else {
		user.WriteString("Interviewer: (not specified — infer the likely interviewer for this round from the role)\n")
	}
	if s := strings.TrimSpace(roundContext); s != "" {
		fmt.Fprintf(&user, "This round: %s\n", s)
	}
	if s := strings.TrimSpace(priorDebriefs); s != "" {
		fmt.Fprintf(&user, "\nEarlier rounds in THIS process already happened — the candidate's own debriefs:\n%s\nUse these to tailor this round: build on what landed, shore up what felt shaky, anticipate follow-ups to what already came up, and don't re-tread ground already covered. Weave this into snapshot/lands/avoid/questions where it helps — do not invent a separate section.\n", s)
	}
	return user.String()
}

// GenerateInterviewerBrief researches one interviewer for a single round.
// location and companyURL ground the research to the right (same-named) company.
// roundContext (may be empty) is the round's title/description — e.g. "Home
// Assignment Presentation" — so the prep is tailored to the round's FORMAT.
// priorDebriefs (may be empty) summarises how earlier rounds went so this round's
// prep can build on what already happened — the debrief feed-forward loop.
func (c *Client) GenerateInterviewerBrief(ctx context.Context, company, role, interviewerName, location, companyURL, roundContext, priorDebriefs string) (json.RawMessage, error) {
	cctx, cancel := context.WithTimeout(ctx, 150*time.Second)
	defer cancel()

	user := buildInterviewerPromptUser(company, role, interviewerName, location, companyURL, roundContext, priorDebriefs)

	resp, err := c.CreateMessage(cctx, ModelSonnet, interviewerBriefSystemPrompt, []Message{
		{Role: "user", Content: user},
	}, 4000, WebSearchTool(5))
	if err != nil {
		return nil, err
	}
	raw, err := extractFirstJSONObject(resp)
	if err != nil {
		return nil, fmt.Errorf("model didn't return JSON (raw: %.300s)", resp)
	}
	if msg := jsonErrorField(raw); msg != "" {
		return nil, errors.New(msg)
	}
	var shape struct {
		Interviewer struct {
			Name string `json:"name"`
		} `json:"interviewer"`
		Snapshot string `json:"snapshot"`
	}
	if err := json.Unmarshal(raw, &shape); err != nil {
		return nil, fmt.Errorf("interviewer brief JSON missing expected fields: %w", err)
	}
	if shape.Interviewer.Name == "" && shape.Snapshot == "" {
		return nil, errors.New("interviewer brief came back empty")
	}
	return raw, nil
}

func jsonErrorField(raw json.RawMessage) string {
	var errProbe struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(raw, &errProbe) == nil {
		return errProbe.Error
	}
	return ""
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
