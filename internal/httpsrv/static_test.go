package httpsrv

import (
	"strings"
	"testing"

	"github.com/backyonatan-alt/jobsearch/internal/config"
)

const sampleHTML = `<!doctype html><html><head><title>Pursuit</title></head><body>x</body></html>`

func TestInjectGA_NoID_Untouched(t *testing.T) {
	s := &Server{Cfg: &config.Config{GA4MeasurementID: ""}}
	got := string(s.injectGA([]byte(sampleHTML)))
	if got != sampleHTML {
		t.Fatalf("expected untouched HTML with no measurement id, got:\n%s", got)
	}
	if strings.Contains(got, "googletagmanager") {
		t.Fatalf("gtag snippet leaked when no id configured")
	}
}

func TestInjectGA_WithID_InjectsBeforeHeadClose(t *testing.T) {
	s := &Server{Cfg: &config.Config{GA4MeasurementID: "G-TEST123"}}
	got := string(s.injectGA([]byte(sampleHTML)))

	if !strings.Contains(got, `src="https://www.googletagmanager.com/gtag/js?id=G-TEST123"`) {
		t.Fatalf("loader script missing or wrong id:\n%s", got)
	}
	if !strings.Contains(got, `gtag('config','G-TEST123',{send_page_view:false})`) {
		t.Fatalf("config call missing or send_page_view not disabled:\n%s", got)
	}
	// Snippet must sit inside <head>, immediately before </head>.
	if !strings.Contains(got, `</script></head>`) {
		t.Fatalf("snippet not injected directly before </head>:\n%s", got)
	}
	if strings.Count(got, "</head>") != 1 {
		t.Fatalf("expected exactly one </head>, got %d", strings.Count(got, "</head>"))
	}
}
