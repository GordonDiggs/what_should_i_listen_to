package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
)

const (
	GOOGLE_URL = "http://ajax.googleapis.com/ajax/services/search/images?v=1.0&rsz=8&start=1&q="
	RANDOM_URL = "https://x-vinyl.herokuapp.com/items/random.json"
	INDEX      = `
    <html>
      <head>
        <title>What should I listen to?</title>
        <style type='text/css'>
          body { background: black; color: white; font-size: 36px; text-align:center; font-family: sans-serif; }
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

type Record struct {
	Title  string
	Artist string
	Format string
	Label  string
}

/*
func getAlbumArt(r Record) string {
	q := url.QueryEscape(fmt.Sprintf("%s %s", r.Artist, r.Title))
	req_url := fmt.Sprintf("%s%s", GOOGLE_URL, q)
	return req_url
}
*/

func GetRecord() Record {
	var r Record
	response, _ := http.Get(RANDOM_URL)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)
	return r
}

func handler(w http.ResponseWriter, r *http.Request) {
	record := GetRecord()
	t := template.New("Record Template")
	t, _ = t.Parse(INDEX)
	t.Execute(w, record)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
