package api

import (
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"

	embeddedwebui "github.com/Resinat/Resin/webui"
)

func registerEmbeddedWebUI(mux *http.ServeMux) {
	distFS, err := embeddedwebui.DistFS()
	if err != nil {
		log.Printf("WebUI embed disabled: %v", err)
		return
	}
	mux.Handle("/", newRootRedirectHandler())
	mux.Handle("/ui", newUIRootRedirectHandler())
	mux.Handle("GET /llms.txt", publicMarkdownHandler(distFS, "llms.txt"))
	mux.Handle("GET /docs.md", publicGuideHandler(distFS))
	mux.Handle("/ui/", newWebUIHandler(distFS))
}

func publicGuideHandler(distFS fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := "user-guide.zh-CN.md"
		if strings.HasPrefix(strings.ToLower(r.URL.Query().Get("lang")), "en") {
			name = "user-guide.en.md"
		}
		publicMarkdownHandler(distFS, name).ServeHTTP(w, r)
	})
}

func publicMarkdownHandler(distFS fs.FS, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}
		data, err := fs.ReadFile(distFS, name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if name == "llms.txt" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		} else {
			w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		}
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})
}

func newWebUIHandler(distFS fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/ui/") {
			http.NotFound(w, r)
			return
		}

		assetPath := strings.TrimPrefix(path.Clean("/"+strings.TrimPrefix(r.URL.Path, "/ui/")), "/")
		if assetPath == "" || assetPath == "." {
			assetPath = "index.html"
		}

		if info, err := fs.Stat(distFS, assetPath); err == nil && !info.IsDir() {
			http.ServeFileFS(w, r, distFS, assetPath)
			return
		}

		// Missing requests with file-like paths should remain 404.
		if path.Ext(assetPath) != "" {
			http.NotFound(w, r)
			return
		}

		http.ServeFileFS(w, r, distFS, "index.html")
	})
}

func newRootRedirectHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" || (r.Method != http.MethodGet && r.Method != http.MethodHead) {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})
}

func newUIRootRedirectHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ui" || (r.Method != http.MethodGet && r.Method != http.MethodHead) {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})
}
