package httpsrv

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/backyonatan-alt/jobsearch/internal/config"
)

func TestRequireOpsToken(t *testing.T) {
	called := false
	ok := func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	cases := []struct {
		name       string
		opsToken   string
		authHeader string
		wantStatus int
		wantCalled bool
	}{
		{"disabled when unset", "", "Bearer anything", http.StatusNotFound, false},
		{"missing header", "secret", "", http.StatusUnauthorized, false},
		{"wrong token", "secret", "Bearer nope", http.StatusUnauthorized, false},
		{"bare token without prefix accepted", "secret", "secret", http.StatusOK, true},
		{"valid token", "secret", "Bearer secret", http.StatusOK, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called = false
			s := &Server{Cfg: &config.Config{OpsToken: tc.opsToken}}
			req := httptest.NewRequest(http.MethodGet, "/api/ops/beta-interest/pending", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			rec := httptest.NewRecorder()
			s.requireOpsToken(ok)(rec, req)
			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if called != tc.wantCalled {
				t.Errorf("handler called = %v, want %v", called, tc.wantCalled)
			}
		})
	}
}
