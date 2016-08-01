package main

import "os/exec"
import "fmt"
import "io/ioutil"
import "encoding/json"

const phantomjs = "phantomjs"

func screenshot(url string, outputFile string) {
	fmt.Printf("Capturing screenshot for %s \n", url)

	_, findPhantomJSErr := exec.LookPath(phantomjs)
	if findPhantomJSErr != nil {
		panic("Please make sure phantomjs exec in PATH.")
	}

	args := []string{"webshot-phantomjs.js", url, outputFile}

	execErr := exec.Command(phantomjs, args...).Run()

	if execErr != nil {
		fmt.Printf("Can not download url : %s \n ", url)
		/*
		 *panic(execErr)
		 */
	}
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

	hnItems := loadHackerNewsItems("./feeds.json")
	/*
	 *fmt.Printf("Results: %v\n", hackerNewsItems)
	 */
	for _, item := range hnItems {
		screenshot(item.URL, item.BufferedURL)
	}
}
