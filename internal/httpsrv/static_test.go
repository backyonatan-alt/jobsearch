package httpsrv

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/backyonatan-alt/jobsearch/internal/config"
)

const sampleHTML = `<!doctype html><html><head><title>Pursuit</title></head><body>x</body></html>`

// Every client-routed page must fall back to the SPA shell — a route missing
// from the allowlist in Routes() returns a hard 404 on direct load (that's how
// /privacy shipped broken). This guards the whole list, static assets included.
func TestSPARoutesServeShell(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(sampleHTML), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "app.js"), []byte("console.log(1)"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := &Server{
		Cfg:    &config.Config{},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Static: http.Dir(dir),
	}
	h := s.Routes()

	for _, path := range []string{"/", "/app", "/app/1", "/admin", "/preview", "/privacy"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("%s: got %d, want 200 (SPA shell)", path, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "Pursuit") {
			t.Errorf("%s: did not serve the SPA shell", path)
		}
	}

	// A real static asset is still served straight off disk, not the shell.
	req := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "console.log") {
		t.Errorf("static asset /app.js: got %d body=%q", rec.Code, rec.Body.String())
	}
}

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
