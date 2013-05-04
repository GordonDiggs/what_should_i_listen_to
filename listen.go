package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	GOOGLE_URL = "http://ajax.googleapis.com/ajax/services/search/images?v=1.0&rsz=8&start=1&q="
	RANDOM_URL = "https://x-vinyl.herokuapp.com/items/random.json"
	INDEX      = `
    <html>
      <head>
        <title>What should I listen to?</title>
        <style type='text/css'>
          body { background: url({{.Url}}) no-repeat center center fixed;
                 background-size: cover; background-color: black;
                color: white; font-size: 36px; text-align:center; font-family: sans-serif; }
          em { color: yellow; }
          strong { color: red; }
        </style> 
      </head>
      <body>
        <p>
          You should listen to <em>{{.Title}}</em>
          <br>
          by <strong>{{.Artist}}</strong> on {{.Label}}.
          <br>
          It is a {{.Format}} record.
        </p>
      </body>
    </html>
  `
)

type ImageResult struct {
	Unescapedurl string
}

type ImageResults struct {
	Results []ImageResult
}

type ResponseData struct {
	ResponseData ImageResults
}

type Record struct {
	Title  string
	Artist string
	Format string
	Label  string
	Url    string
}

func (r *Record) GetAlbumArt() string {
	var resp ResponseData
	q := url.QueryEscape(fmt.Sprintf("%s %s cover", r.Artist, r.Title))
	req_url := fmt.Sprintf("%s%s", GOOGLE_URL, q)
	response, _ := http.Get(req_url)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &resp)
	return resp.ResponseData.Results[0].Unescapedurl
}

func GetRecord() Record {
	var r Record
	response, _ := http.Get(RANDOM_URL)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)
	return r
}

func handler(w http.ResponseWriter, r *http.Request) {
	record := GetRecord()
	record.Url = record.GetAlbumArt()
	t := template.New("Record Template")
	t, _ = t.Parse(INDEX)
	t.Execute(w, record)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Waiting for requests...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
