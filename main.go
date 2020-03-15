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
	fmt.Println("request from ", request.RemoteAddr)
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
		fmt.Fprintf(writer, err.Error())
		return
	}
	response := contrib{Contributions: contri}
	encoder := json.NewEncoder(writer)
	err = encoder.Encode(&response)
	if err != nil {
		writer.WriteHeader(500)
		return
	}
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

func fetchContributions(url string) (int, error) {
	response, err := http.Get(url)
	if err != nil {
		return -1, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = &statusError{response.StatusCode}
		return -1, err
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return -1, err
	}
	content, err := doc.Find("h2").Html()
	if err != nil {
		return -1, err
	}
	contributions := strings.Split(strings.TrimSpace(content), " ")[0]
	return strconv.Atoi(strings.Join(strings.Split(contributions, ","), ""))
}
