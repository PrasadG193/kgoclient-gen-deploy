package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/PrasadG193/kgoclient-gen/pkg/generator"
)

const API_VERSION = "v1"
const PORT = "8080"

func main() {
	log.Printf("server started accepting requests on port=%s..\n", PORT)
	http.HandleFunc(fmt.Sprintf("/%s/convert", API_VERSION), HandleConvert)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil))
}

func HandleConvert(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	urlPQ, _ := url.ParseQuery(r.URL.RawQuery)
	method := generator.KubeMethod(urlPQ.Get("method"))
	if len(method) == 0 {
		method = generator.MethodCreate
	}
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(body), method, err)
	gen := generator.New(body, method)
	code, err := gen.Generate()
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Bad Request. Error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	io.WriteString(w, code)
}
