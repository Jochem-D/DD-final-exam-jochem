package cmd

import (
	"ddsheetfinal/helpers"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const charactersPath = "/characters/"

// ExecuteServe starts a simple static file server for the HTML frontend
func ExecuteServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	root := fs.String("dir", ".", "project root (contains frontend/ and characters/)")
	port := fs.String("port", "8080", "port to listen on")
	_ = fs.Parse(args)

	projectRoot := *root
	frontendDir := filepath.Join(projectRoot, "frontend")
	charDir := filepath.Join(projectRoot, "characters")

	mux := http.NewServeMux()

	// Serve frontend assets at /frontend/
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir(frontendDir))))

	// Dynamic manifest for the frontend that lists all character JSON files
	mux.HandleFunc("/frontend/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		entries, err := os.ReadDir(charDir)
		if err != nil {
			http.Error(w, "cannot read characters directory", http.StatusInternalServerError)
			return
		}
		var out []string
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.HasSuffix(strings.ToLower(name), ".json") {
				// frontend expects paths it can fetch directly; expose under /characters/
				out = append(out, charactersPath+name)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	})

	// Enrich endpoint (fetches dnd5eapi data, writes enriched file)
	mux.HandleFunc("/api/enrich", func(w http.ResponseWriter, r *http.Request) {
		// delegate to helpers.RunEnrich via the existing HTTP handler in helpers
		// to avoid importing helpers here (and possible cycles), forward to the handler
		// by calling the package-level EnrichHandler.
		// Note: we import helpers at top of file when needed.
		helpers.EnrichHandler(w, r)
	})

	// Derive endpoint: compute derived stats server-side
	mux.HandleFunc("/api/derive", func(w http.ResponseWriter, r *http.Request) {
		helpers.DeriveHandler(w, r)
	})

	// Serve enrichments (generated enriched JSON files)
	// Characters GET (serve file) and PUT (save edits)
	mux.HandleFunc(charactersPath, func(w http.ResponseWriter, r *http.Request) {
		rel := strings.TrimPrefix(r.URL.Path, charactersPath)
		rel = strings.TrimPrefix(rel, "/")
		if rel == "" {
			http.Error(w, "missing filename", http.StatusBadRequest)
			return
		}
		// prevent path traversal
		if strings.Contains(rel, "..") || strings.HasPrefix(rel, "/") {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		target := filepath.Join(charDir, filepath.Clean(rel))

		switch r.Method {
		case http.MethodGet:
			http.ServeFile(w, r, target)
			return
		case http.MethodPut:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "read error", http.StatusBadRequest)
				return
			}
			log.Printf("PUT %s (%d bytes)\n", target, len(body))
			// validate JSON
			var tmp interface{}
			if err := json.Unmarshal(body, &tmp); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			// ensure characters dir exists
			if err := os.MkdirAll(charDir, 0755); err != nil {
				http.Error(w, "cannot create characters dir", http.StatusInternalServerError)
				return
			}
			if err := os.WriteFile(target, body, 0644); err != nil {
				http.Error(w, "save failed", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
			return
		case http.MethodDelete:
			// Remove the target file if it exists.
			if err := os.Remove(target); err != nil {
				if os.IsNotExist(err) {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
				http.Error(w, "delete failed", http.StatusInternalServerError)
				return
			}
			log.Printf("DELETE %s OK\n", target)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	// Root: redirect to frontend UI
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/frontend/index.html", http.StatusFound)
			return
		}
		// Fallback: try to serve from project root
		http.FileServer(http.Dir(projectRoot)).ServeHTTP(w, r)
	})

	addr := ":" + *port
	log.Printf("Serving frontend at http://localhost%s/frontend/ (characters from %s)\n", addr, charDir)
	log.Fatal(http.ListenAndServe(addr, mux))
}
