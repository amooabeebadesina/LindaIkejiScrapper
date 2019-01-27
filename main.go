package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Info is a struct for each data row
type Info struct {
	Title  string   `json:"title"`
	URL    string   `json:"url"`
	Images []string `json:"images"`
}

var url string
var data []Info
var fileName string

func getFileName() string {
	var buffer bytes.Buffer
	currentTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	buffer.WriteString("data_")
	buffer.WriteString(currentTime)
	buffer.WriteString(".json")
	return buffer.String()
}

func scrapeLindaIkeji() {
	// Create an empty json file using the current timestamp
	fileName = getFileName()
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	url := "https://www.lindaikejisblog.com"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	readDataToFile(doc)

}

func readDataToFile(doc *goquery.Document) {

	// Find each story block and get data
	doc.Find(".story_block").Each(func(i int, s *goquery.Selection) {
		info := Info{}
		title := s.Find("h1").Text()
		link, hasHref := s.Find("a").Eq(0).Attr("href")
		if !hasHref {
			link = " "
		}
		info.Title = title
		info.URL = link

		s.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			info.Images = append(info.Images, src)
		})
		data = append(data, info)
	})
	jsonData, _ := json.Marshal(data)
	_ = ioutil.WriteFile(fileName, jsonData, 0644)
}

func main() {
	scrapeLindaIkeji()
}
