package main

import "os/exec"
import "fmt"
import "io/ioutil"
import "encoding/json"
import "github.com/PuerkitoBio/goquery"
import "log"
import "crypto/md5"
import "github.com/fatih/color"
import "os"

const phantomjs = "phantomjs"

func screenshot(url string, outputFile string) {
	if _, err := os.Stat(outputFile); err == nil {
		color.Blue("Screenshot for %s already exists, skip.", url)
		return
	}

	fmt.Printf("Capturing screenshot for %s \n", url)

	_, findPhantomJSErr := exec.LookPath(phantomjs)
	if findPhantomJSErr != nil {
		panic("Please make sure phantomjs exec in PATH.")
	}

	args := []string{"webshot-phantomjs.js", url, outputFile}

	execErr := exec.Command(phantomjs, args...).Run()

	if execErr != nil {
		color.Red("Can not download url : %s \n ", url)
	}
}

func fetchHackerNewsItems(where string) []HackerNewsItem {
	doc, err := goquery.NewDocument(fmt.Sprintf("https://news.ycombinator.com/%s", where))
	if err != nil {
		log.Fatal(err)
	}

	items := []HackerNewsItem{}

	// Find the review items
	doc.Find("tr.athing").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		id, _ := s.Attr("id")
		title := s.Find(".title a").Text()
		url, _ := s.Find(".title a").Attr("href")
		site := s.Find(".sitebit.comhead .sitestr").Text()
		bufferedURL := fmt.Sprintf("%x.png", md5.Sum([]byte(url)))

		item := HackerNewsItem{Title: title, URL: url, CommentURL: fmt.Sprintf("https://news.ycombinator.com/item?id=%s", id), Site: site, BufferedURL: bufferedURL}
		items = append(items, item)
	})

	return items
}

// HackerNewsItem is a representation of hacker news items
type HackerNewsItem struct {
	Title       string
	URL         string
	Site        string
	BufferedURL string
	CommentURL  string
}

func loadHackerNewsItems(filename string) []HackerNewsItem {
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		panic("File error")
	}

	var hackerNewsItems []HackerNewsItem
	json.Unmarshal(file, &hackerNewsItems)
	return hackerNewsItems
}

func main() {

	items := fetchHackerNewsItems("")
	bolB, _ := json.Marshal(items)
	ioutil.WriteFile("./feeds.json", bolB, 0644)

	items = fetchHackerNewsItems("show")
	bolB, _ = json.Marshal(items)
	ioutil.WriteFile("./show.json", bolB, 0644)

	items = fetchHackerNewsItems("newest")
	bolB, _ = json.Marshal(items)
	ioutil.WriteFile("./newest.json", bolB, 0644)

	for _, item := range loadHackerNewsItems("./feeds.json") {
		screenshot(item.URL, item.BufferedURL)
	}

	for _, item := range loadHackerNewsItems("./newest.json") {
		screenshot(item.URL, item.BufferedURL)
	}

	for _, item := range loadHackerNewsItems("./show.json") {
		screenshot(item.URL, item.BufferedURL)
	}

}
