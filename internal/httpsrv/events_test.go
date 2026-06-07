package httpsrv

import "testing"

func TestEventNameValidation(t *testing.T) {
	valid := []string{"paste_parse", "screenshot_parse", "status_change", "first_application", "a", "dossier_open"}
	for _, n := range valid {
		if !eventNameRe.MatchString(n) {
			t.Errorf("expected %q to be a valid event name", n)
		}
	}
	invalid := []string{
		"",                  // empty
		"Paste_Parse",       // uppercase
		"1foo",              // leading digit
		"_foo",              // leading underscore
		"foo-bar",           // hyphen
		"foo bar",           // space
		"email@example.com", // PII-ish / illegal chars
		"this_event_name_is_far_too_long_to_be_allowed_here", // > 40 chars
	}
	for _, n := range invalid {
		if eventNameRe.MatchString(n) {
			t.Errorf("expected %q to be rejected", n)
		}
	}
}
