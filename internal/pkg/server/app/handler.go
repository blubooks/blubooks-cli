package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mainApp "github.com/blubooks/blubooks-cli/internal/app"
)

func (serverApp *App) HanlderHealth(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Write the status code using w.WriteHeader
	res.WriteHeader(http.StatusOK)

	// Write the body text using w.Write
	res.Write([]byte("OK"))
}

func (serverApp *App) HandleIndex(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	err := mainApp.Build(true)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Fehler: %v", err)
		fmt.Fprintf(res, `{"error.message": "%v"}`, err)
		return
	}
	/*
		if warn != nil {
			fmt.Fprintf(res, `{"warn": %v}`, warn)

		}
	*/

	/*
		err, warn := Build("http://localhost:3020/public/", "", "", "")
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Printf("Fehler: %v", err)
			fmt.Fprintf(res, `{"error.message": "%v"}`, appErr)
			return
		}

		if warn != nil {
			fmt.Fprintf(res, `{"warn": %v}`, warn)

		}
	*/

}

type (
	// FallbackResponseWriter wraps an http.Requesthandler and surpresses
	// a 404 status code. In such case a given local file will be served.
	FallbackResponseWriter struct {
		WrappedResponseWriter http.ResponseWriter
		FileNotFound          bool
	}
)

func (frw *FallbackResponseWriter) WriteHeader(statusCode int) {
	//log.Printf("INFO: WriteHeader called with code %d\n", statusCode)
	if statusCode == http.StatusNotFound {
		//log.Printf("INFO: Setting FileNotFound flag\n")
		frw.FileNotFound = true
		return
	}
	frw.WrappedResponseWriter.WriteHeader(statusCode)
}

// Header returns the header of the wrapped response writer
func (frw *FallbackResponseWriter) Header() http.Header {
	return frw.WrappedResponseWriter.Header()
}

// Write sends bytes to wrapped response writer, in case of FileNotFound
// It surpresses further writes (concealing the fact though)
func (frw *FallbackResponseWriter) Write(b []byte) (int, error) {
	if frw.FileNotFound {
		return len(b), nil
	}
	return frw.WrappedResponseWriter.Write(b)
}

func (serverApp *App) HandleClient(w http.ResponseWriter, r *http.Request) {

	frw := FallbackResponseWriter{
		WrappedResponseWriter: w,
		FileNotFound:          false,
	}

	/*
		http.FileServer(http.Dir("public")).ServeHTTP(&frw, r)

		if frw.FileNotFound {
			b, _ := os.ReadFile("public/index.html")
			w.Header().Set("Content-Type", "text/html")
			w.Write(b)
			return
		}
	*/

	fs := http.FileServer(http.Dir("client/default"))

	http.StripPrefix("/public", fs).ServeHTTP(&frw, r)

	if frw.FileNotFound {
		b, _ := os.ReadFile("client/default/index.html")
		w.Header().Set("Content-Type", "text/html")
		w.Write(b)
		return
	}

}
