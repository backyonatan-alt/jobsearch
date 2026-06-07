package httpsrv

import (
	"bytes"
	"io"
	"net/http"
)

// serveIndexHTML serves the SvelteKit SPA shell (index.html) for every
// client-routed path. When a GA4 measurement ID is configured it injects the
// gtag snippet right before </head>; with no ID set (dev/local) the document is
// served untouched, so nothing analytics-related loads off prod.
func (s *Server) serveIndexHTML(w http.ResponseWriter, r *http.Request) {
	f, err := s.Static.Open("/index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "stat", http.StatusInternalServerError)
		return
	}
	data, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, "read", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, "index.html", fi.ModTime(), bytes.NewReader(s.injectGA(data)))
}

func (s *Server) injectGA(html []byte) []byte {
	id := s.Cfg.GA4MeasurementID
	if id == "" {
		return html
	}
	snippet := `<script async src="https://www.googletagmanager.com/gtag/js?id=` + id + `"></script>` +
		`<script>window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}` +
		`gtag('js',new Date());gtag('config','` + id + `',{send_page_view:false});</script>`
	return bytes.Replace(html, []byte("</head>"), []byte(snippet+"</head>"), 1)
}
