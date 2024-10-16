package main

import (
	"encoding/json"
	"fmt"
	"io"
	"main/data"
	"net/http"
	"os"
	"path/filepath"
)

func middlewareValidation(next http.Handler, accMethods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range accMethods {
			if method == r.Method {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Method is not allowed", http.StatusBadRequest)
		return
	})
}
func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Today, you will have a cup of tea :D")
}
func postHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//dst: destination
	dst, err := os.Create(filepath.Join("uploads", handler.Filename))
	if err != nil {
		panic(err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "I have received your file")
}

func jsonMessageH(w http.ResponseWriter, r *http.Request) {
	jsonData := r.FormValue("Person")
	var person data.Person
	err := json.Unmarshal([]byte(jsonData), &person)
	if err != nil {
		panic(err)
	}
	fmt.Println("Received json data: ", person)

	//response
	fmt.Fprintln(w, "I got your JSON")
}
func main() {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/get", middlewareValidation(http.HandlerFunc(getHandler), "GET"))
	mux.Handle("/upload", middlewareValidation(http.HandlerFunc(postHandler), "POST"))
	mux.Handle("/json", middlewareValidation(http.HandlerFunc(jsonMessageH), "POST"))

	server := http.Server{
		Addr:    "localhost:5678",
		Handler: mux,
	}
	server.ListenAndServe()
}
