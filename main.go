package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"ascii-art-web/ascii"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}, statuscode int) {
	w.WriteHeader(statuscode)
	t := template.Must(template.ParseFiles(tmpl))
	t.Execute(w, data)
}

func Restrict(next http.HandlerFunc) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		restrictedPaths := []string{"/static", "/images"}
		for _, path := range restrictedPaths {
			if r.URL.Path == path || r.URL.Path == path+"/" {
				renderTemplate(w, "templates/403.html", nil, 403)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderTemplate(w, "templates/400.html", nil, 400)
		return
	}
	if r.URL.Path != "/" {
		renderTemplate(w, "templates/404.html", nil, 404)
		return
	}
	renderTemplate(w, "templates/index.html", nil, 200)
}

func Ascii(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		renderTemplate(w, "templates/400.html", nil, 400)
		return
	}

	input := r.FormValue("text")
	style := r.FormValue("Style")
	if style != "standard" && style != "thinkertoy" && style != "shadow" {
		renderTemplate(w, "templates/400_invalidEntry.html", nil, 400)
		return
	}

	PrintArt, unprintable := ascii.PrintArt(input, style)
	data := struct {
		Output             string
		UnprintableWarning bool
	}{
		Output:             PrintArt,
		UnprintableWarning: unprintable,
	}
	renderTemplate(w, "templates/index.html", data, 200)
}

func downloadText(w http.ResponseWriter, r *http.Request) {
	output := r.FormValue("output")
	if output == "" {
		renderTemplate(w, "templates/500_NoContent.html", nil, 500)
		return

	}
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_art.txt")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(output)))
	_, err := io.WriteString(w, output)
	if err != nil {
		renderTemplate(w, "templates/500.html", nil, 500)
		return

	}
}

func downloadHTML(w http.ResponseWriter, r *http.Request) {
	output := r.FormValue("output")
	if output == "" {
		renderTemplate(w, "templates/500_NoContent.html", nil, 500)
		return
	}

	htmlContent := fmt.Sprintf("<pre>%s</pre>", output)
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_art.html")
	w.Header().Set("Content-Length", strconv.Itoa(len(htmlContent)))
	w.Header().Set("Content-Type", "text/html")
	_, err := io.WriteString(w, htmlContent)
	if err != nil {
		renderTemplate(w, "templates/500.html", nil, 500)
		return

	}
}

func About(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/About.html", nil, 200)
}

func readME(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/readme.html", nil, 200)
}

func main() {
	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	imageFs := http.FileServer(http.Dir(filepath.Join("templates", "images")))
	http.Handle("/images/", http.StripPrefix("/images/", imageFs))

	http.HandleFunc("/", Home)
	http.HandleFunc("/ascii", Ascii)
	http.HandleFunc("/download/txt", downloadText)
	http.HandleFunc("/download/html", downloadHTML)
	http.HandleFunc("/about", About)
	http.HandleFunc("/readme", readME)

	fmt.Println("local host running : http://localhost:8080")
	http.ListenAndServe(":8080", Restrict(http.DefaultServeMux.ServeHTTP))
}
