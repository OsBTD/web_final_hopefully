package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"ascii-art-web/ascii"
)

// renderTemplate is a helper function that parses and executes,
// it takes 4 parameters, the responsewriter (interface), tmpl (string) which'll represent the file to parse,
// data(interface) which is whatever we wanna execute on the templates, and statuscode (int) which is just the status code
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}, statuscode int) {
	w.WriteHeader(statuscode)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Fprintf(w, "error parsing files")
		return
	}
	err2 := t.Execute(w, data)
	if err2 != nil {
		renderTemplate(w, "templates/500.html", nil, http.StatusInternalServerError)
		return
	}
}

// restrict is a middleware that restricts access to specific paths, /static and /images in this case
// it takes a next(handlerfunc) and returns an http handler function that checks if our path is one of the restricted ones if so the file to parse and execute would be the 403 template and status is 403 forbidden
// if the path doesn't figure in our restricted ones the handlerfunc is returned the usual way and the file to be parsed and executed would be determined
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
	// the data struct contains the boolean and the ascii art which is what would be executed
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

// this custom file sever allows to customize the errors in file serving
// for example if a file we're trying to serve doesn't exist or if we don't have the necessary permissions
// otherwise if a file doesn't exist for example a standard 404 error would be displayed
// it takes the root as parameter and returns a handler
// it uses os.Stat which returns meta data about a file on our os system
// and returns an error which could be of two type
// either the file doesn't exist or permissions denied
// we're using serve file instead of fileserver cause it only serves one file instead of a whole directory
func customFileServer(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(root, r.URL.Path)
		if _, err := os.Stat(path); err != nil {
			renderTemplate(w, "templates/404.html", nil, http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, path)
	})
}

func main() {
	// initialize the servemux multiplexer with NewServeMux so we can wrapp it in the restrict function afterwards
	mux := http.NewServeMux()
	// handle the handler functions
	// we're using mux.HandleFunc because http.HandleFunc would automatically use the DefaultServeMux
	// HandleFunc implecitely converts these functions to handler functions
	// it's as if we wrote it like mux.HandleFunc("path", http.HandlerFunc(Function))
	// the handlerfunc part is where the ServeHttp method is applied
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/ascii", Ascii)
	mux.HandleFunc("/download/txt", downloadText)
	mux.HandleFunc("/download/html", downloadHTML)
	mux.HandleFunc("/about", About)
	mux.HandleFunc("/readme", readME)
	// these are the file handlers for images and the css file
	// they use the custome handler which checks if a file exists
	// if not it returns a custom error page
	// we give it the root as a parameter (templates)
	// we strip the static prefix, it doesn't exist in our filesystem it's just how we server the static files on the html
	// we use mux.Handle here cause customeServer returns a Handler

	mux.Handle("/static/", http.StripPrefix("/static/", customFileServer("templates")))
	mux.Handle("/images/", customFileServer("templates"))
	// we give the mux.ServeHTTP as a parameter to the restrict function
	// it would check if a path is restricted before going to treating it with the mux
	fmt.Println("local host running : http://localhost:8088")
	http.ListenAndServe(":8088", Restrict(mux.ServeHTTP))
}
