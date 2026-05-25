// Package ics is a minimal iCalendar (RFC 5545) parser scoped to what we need
// for interview scheduling: VEVENT blocks with UID, SUMMARY, LOCATION,
// DESCRIPTION, DTSTART, DTEND, ORGANIZER, ATTENDEE. RRULE/recurrence,
// VTIMEZONE bodies, alarms, and attachments are deliberately ignored.
package ics

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

type Event struct {
	UID         string
	Summary     string
	Location    string
	Description string
	StartsAt    time.Time
	EndsAt      time.Time
	AllDay      bool
	Organizer   Person
	Attendees   []Person
}

type Person struct {
	Name  string
	Email string
}

// Parse reads an iCalendar payload and returns every VEVENT found, in source
// order. Lines that fail to parse individually are skipped rather than failing
// the whole document — calendar exports from the wild routinely include
// unfamiliar properties.
func Parse(r io.Reader) ([]Event, error) {
	lines, err := unfold(r)
	if err != nil {
		return nil, err
	}

	var events []Event
	var cur *Event
	inEvent := false
	depth := 0 // nesting depth so VALARM inside VEVENT doesn't end the event early

	for _, line := range lines {
		name, params, value := splitLine(line)
		upper := strings.ToUpper(name)

		switch upper {
		case "BEGIN":
			v := strings.ToUpper(value)
			if v == "VEVENT" && !inEvent {
				inEvent = true
				cur = &Event{}
				depth = 0
				continue
			}
			if inEvent {
				depth++
			}
		case "END":
			v := strings.ToUpper(value)
			if v == "VEVENT" && inEvent && depth == 0 {
				if cur != nil {
					events = append(events, *cur)
				}
				cur = nil
				inEvent = false
				continue
			}
			if inEvent && depth > 0 {
				depth--
			}
		}

		if !inEvent || cur == nil || depth > 0 {
			continue
		}

		switch upper {
		case "UID":
			cur.UID = unescape(value)
		case "SUMMARY":
			cur.Summary = unescape(value)
		case "LOCATION":
			cur.Location = unescape(value)
		case "DESCRIPTION":
			cur.Description = unescape(value)
		case "DTSTART":
			if t, allDay, err := parseTime(params, value); err == nil {
				cur.StartsAt = t
				cur.AllDay = allDay
			}
		case "DTEND":
			if t, _, err := parseTime(params, value); err == nil {
				cur.EndsAt = t
			}
		case "ORGANIZER":
			cur.Organizer = parsePerson(params, value)
		case "ATTENDEE":
			cur.Attendees = append(cur.Attendees, parsePerson(params, value))
		}
	}

	return events, nil
}

// unfold reads the input and joins RFC-5545 line continuations: any line
// beginning with a space or HTAB is appended to the previous line with the
// leading whitespace stripped.
func unfold(r io.Reader) ([]string, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 4*1024*1024) // tolerate long DESCRIPTION lines
	var out []string
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		if line == "" {
			continue
		}
		if (line[0] == ' ' || line[0] == '\t') && len(out) > 0 {
			out[len(out)-1] += line[1:]
			continue
		}
		out = append(out, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// splitLine breaks a content line into (name, params, value).
// e.g. "DTSTART;TZID=America/New_York:20260525T100000"
//
//	→ name="DTSTART", params={"TZID":"America/New_York"}, value="20260525T100000"
func splitLine(line string) (name string, params map[string]string, value string) {
	colon := strings.IndexByte(line, ':')
	if colon < 0 {
		return "", nil, ""
	}
	left, value := line[:colon], line[colon+1:]
	params = map[string]string{}
	if semi := strings.IndexByte(left, ';'); semi >= 0 {
		name = left[:semi]
		for _, p := range splitParams(left[semi+1:]) {
			if eq := strings.IndexByte(p, '='); eq > 0 {
				params[strings.ToUpper(p[:eq])] = strings.Trim(p[eq+1:], `"`)
			}
		}
	} else {
		name = left
	}
	return name, params, value
}

// splitParams splits on ';' but respects double-quoted regions so that
// parameter values containing ';' (e.g. quoted CNs) survive.
func splitParams(s string) []string {
	var out []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			if depth == 0 {
				depth = 1
			} else {
				depth = 0
			}
		case ';':
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	out = append(out, s[start:])
	return out
}

// parseTime handles the three forms we care about:
//   - "20260525T140000Z"             — UTC
//   - "20260525T140000"  + TZID=...  — wall time in a named zone
//   - "20260525"         + VALUE=DATE — all-day
func parseTime(params map[string]string, value string) (time.Time, bool, error) {
	value = strings.TrimSpace(value)
	if params["VALUE"] == "DATE" || (len(value) == 8 && !strings.ContainsAny(value, "T")) {
		t, err := time.Parse("20060102", value)
		return t, true, err
	}
	if strings.HasSuffix(value, "Z") {
		t, err := time.Parse("20060102T150405Z", value)
		return t, false, err
	}
	if tzid := params["TZID"]; tzid != "" {
		loc, err := time.LoadLocation(tzid)
		if err != nil {
			// Unknown zone — fall back to UTC rather than dropping the event.
			loc = time.UTC
		}
		t, err := time.ParseInLocation("20060102T150405", value, loc)
		return t, false, err
	}
	// "Floating" time — no zone information. Treat as UTC; the caller can
	// surface a warning if it matters.
	t, err := time.ParseInLocation("20060102T150405", value, time.UTC)
	return t, false, err
}

// parsePerson pulls a display name + email from an ORGANIZER/ATTENDEE line.
// Typical shape: "ATTENDEE;CN=Jane Doe;ROLE=REQ-PARTICIPANT:mailto:jane@x.com"
func parsePerson(params map[string]string, value string) Person {
	p := Person{Name: params["CN"]}
	v := strings.TrimSpace(value)
	if strings.HasPrefix(strings.ToLower(v), "mailto:") {
		p.Email = v[len("mailto:"):]
	} else if strings.Contains(v, "@") {
		p.Email = v
	}
	return p
}

// unescape reverses the iCalendar text escaping rules from RFC 5545 §3.3.11.
func unescape(s string) string {
	if !strings.ContainsRune(s, '\\') {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' || i+1 >= len(s) {
			b.WriteByte(c)
			continue
		}
		switch s[i+1] {
		case 'n', 'N':
			b.WriteByte('\n')
		case '\\':
			b.WriteByte('\\')
		case ',':
			b.WriteByte(',')
		case ';':
			b.WriteByte(';')
		default:
			b.WriteByte(c)
			b.WriteByte(s[i+1])
		}
		i++
	}
	return b.String()
}

// Duration returns the event's duration. Returns zero if either end is unset.
func (e Event) Duration() time.Duration {
	if e.StartsAt.IsZero() || e.EndsAt.IsZero() {
		return 0
	}
	return e.EndsAt.Sub(e.StartsAt)
}

// String is a debug helper.
func (e Event) String() string {
	return fmt.Sprintf("Event{UID:%q Summary:%q Start:%s End:%s}", e.UID, e.Summary, e.StartsAt, e.EndsAt)
}
