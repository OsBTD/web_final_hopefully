package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"ascii-art-web/ascii"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}, statuscode int) {
	w.WriteHeader(statuscode)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Fprintf(w, "error parsing files")
	}
	err2 := t.Execute(w, data)
	if err2 != nil {
		renderTemplate(w, "templates/500.html", nil, http.StatusInternalServerError)
	}
}

func Restrict(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		restrictedPaths := []string{"/static", "/images", "/static/images"}
		for _, path := range restrictedPaths {
			if r.URL.Path == path || r.URL.Path == path+"/" {
				renderTemplate(w, "templates/403.html", nil, http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderTemplate(w, "templates/400.html", nil, http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		renderTemplate(w, "templates/404.html", nil, http.StatusNotFound)
		return
	}
	renderTemplate(w, "templates/index.html", nil, http.StatusOK)
}

func Ascii(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		renderTemplate(w, "templates/400.html", nil, http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/ascii" {
		renderTemplate(w, "templates/404.html", nil, http.StatusNotFound)
		return
	}

	input := r.FormValue("text")
	style := r.FormValue("Style")
	if style != "standard" && style != "thinkertoy" && style != "shadow" {
		renderTemplate(w, "templates/400_invalidEntry.html", nil, http.StatusBadRequest)
		return
	}

	PrintArt, unprintable, x := ascii.PrintArt(input, style)
	// print art is the ascii art
	// unprintable is a boolean that's true when input contains unprintable chars
	// x has the value of 420 when there's an error reading the files on the ascii code
	if x == 420 {
		renderTemplate(w, "templates/500.html", nil, http.StatusInternalServerError)
		return
	}

	data := struct {
		Output             string
		UnprintableWarning bool
	}{
		Output:             PrintArt,
		UnprintableWarning: unprintable,
	}
	renderTemplate(w, "templates/index.html", data, http.StatusOK)
}

func downloadText(w http.ResponseWriter, r *http.Request) {
	output := r.FormValue("output")
	if output == "" {
		renderTemplate(w, "templates/500_NoContent.html", nil, http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_art.txt")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(output)))
	_, err := io.WriteString(w, output)
	if err != nil {
		renderTemplate(w, "templates/500.html", nil, http.StatusInternalServerError)
		return

	}
}

func downloadHTML(w http.ResponseWriter, r *http.Request) {
	output := r.FormValue("output")
	if output == "" {
		renderTemplate(w, "templates/500_NoContent.html", nil, http.StatusInternalServerError)
		return
	}

	htmlContent := fmt.Sprintf("<pre>%s</pre>", output)
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_art.html")
	w.Header().Set("Content-Length", strconv.Itoa(len(htmlContent)))
	w.Header().Set("Content-Type", "text/html")
	_, err := io.WriteString(w, htmlContent)
	if err != nil {
		renderTemplate(w, "templates/500.html", nil, http.StatusInternalServerError)
		return

	}
}

func About(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/About.html", nil, http.StatusOK)
}

func readME(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/readme.html", nil, http.StatusOK)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Home)
	mux.HandleFunc("/ascii", Ascii)
	mux.HandleFunc("/download/txt", downloadText)
	mux.HandleFunc("/download/html", downloadHTML)
	mux.HandleFunc("/about", About)
	mux.HandleFunc("/readme", readME)

	staticHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("templates", r.URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			renderTemplate(w, "templates/404.html", nil, http.StatusNotFound)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/static/") {
			http.StripPrefix("/static/", http.FileServer(http.Dir("templates"))).ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/images/") {
			http.StripPrefix("/images/", http.FileServer(http.Dir(filepath.Join("templates", "images")))).ServeHTTP(w, r)
			return
		}
	})

	mux.Handle("/static/", staticHandler)
	mux.Handle("/images/", staticHandler)

	fmt.Println("local host running : http://localhost:8080")
	http.ListenAndServe(":8080", Restrict(mux.ServeHTTP))
}
