package api

import (
	"io/fs"
	"net/http"
	"strings"
)

// spaHandler returns an http.Handler that serves static files from the
// embedded filesystem and falls back to index.html for client-side routes
// (Single Page Application behaviour).
//
// The embedded FS is expected to have files under "dist/..." so we
// sub-filesystem into "dist" first.
func spaHandler(embeddedFS fs.FS) http.Handler {
	distFS, err := fs.Sub(embeddedFS, "dist")
	if err != nil {
		panic("embedded dist/ sub-filesystem not found: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(distFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clean path
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else {
			path = strings.TrimPrefix(path, "/")
		}

		// Try to open the file in the embedded FS.
		// If it exists, serve it. Otherwise serve index.html (SPA fallback).
		if f, err := distFS.Open(path); err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for client-side routing
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
