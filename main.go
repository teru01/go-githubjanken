package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	server := http.Server{
		Addr: getListenPort(),
	}
	http.HandleFunc("/contributions", handler)
	server.ListenAndServe()
}

func getListenPort() string {
	port := os.Getenv("PORT")
	if port != "" {
		return ":" + port
	}
	return ":3000"
}

func handler(writer http.ResponseWriter, request *http.Request) {
	githubID := request.FormValue("github_id")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")
	contri, err := fetchContributions(fmt.Sprintf("https://github.com/users/%v/contributions", githubID))
	if err != nil {
		switch err := err.(type) {
		case *statusError:
			writer.WriteHeader(err.code)
		default:
			writer.WriteHeader(500)
		}
		return
	}
	response := contrib{Contributions: contri}
	json, err := json.MarshalIndent(&response, "", "\t\t")
	if err != nil {
		writer.WriteHeader(500)
		return
	}
	fmt.Fprintln(writer, json)
}

type contrib struct {
	Contributions int `json:"contributions"`
}

type statusError struct {
	code int
}

func (e *statusError) Error() string {
	return string(e.code)
}

func fetchContributions(url string) (contri int, err error) {
	response, err := http.Get(url)
	if err != nil {
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = &statusError{response.StatusCode}
		return
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return
	}
	content, err := doc.Find("h2").Html()
	if err != nil {
		return
	}
	contri, err = strconv.Atoi(strings.Split(strings.TrimSpace(content), " ")[0])
	return
}
