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
import "bytes"
import "path/filepath"
import (
	"github.com/disintegration/imaging"
	"image"
	sysColor "image/color"
)
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

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

func bgFit(src image.Image, w, h int, bgColor sysColor.Color) image.Image {
	dst := imaging.New(w, h, bgColor)
	tmp := imaging.Fill(src, w, h, imaging.TopLeft, imaging.CatmullRom)
	return imaging.PasteCenter(dst, tmp)
}

// Resize image to fil fill
func resizeImage(imageName string) {
	if _, err := os.Stat(imageName); err != nil {
		return
	}
	src, err := imaging.Open(imageName)
	if err != nil {
		panic(err)
	}

	testTransparent := bgFit(src, 1024, 768, sysColor.Transparent)
	err = imaging.Save(testTransparent, imageName)
	if err != nil {
		panic(err)
	}
}

// Fetch hacker news items
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

func s3ListObjects() {

	sess := session.New(&aws.Config{Region: aws.String("us-west-1")})

	svc := s3.New(sess)

	allFiles, _ := filepath.Glob("./*.pnga")
	for _, fileName := range allFiles {
		file, _ := ioutil.ReadFile(fileName)

		fmt.Printf("Uploading %s\n", fileName)
		_, err := svc.PutObject(&s3.PutObjectInput{
			Bucket:      aws.String("hackernews-screenshots"), // Required
			Key:         aws.String(fileName),
			ACL:         aws.String("public-read"),
			Body:        bytes.NewReader(file),
			ContentType: aws.String("image/png"),
		})
		fmt.Println("Done")

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return
		}

	}

	// upload json files
	allFiles, _ = filepath.Glob("./*.json")
	for _, fileName := range allFiles {
		file, _ := ioutil.ReadFile(fileName)

		fmt.Printf("Uploading %s\n", fileName)
		_, err := svc.PutObject(&s3.PutObjectInput{
			Bucket:      aws.String("hackernews-screenshots"), // Required
			Key:         aws.String(fileName),
			ACL:         aws.String("public-read"),
			Body:        bytes.NewReader(file),
			ContentType: aws.String("application/json"),
		})
		fmt.Println("Done")

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return
		}

	}

}

func main() {
	/*
	 *s3ListObjects()
	 */

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
		resizeImage(item.BufferedURL)
	}

	for _, item := range loadHackerNewsItems("./newest.json") {
		screenshot(item.URL, item.BufferedURL)
		resizeImage(item.BufferedURL)
	}

	for _, item := range loadHackerNewsItems("./show.json") {
		screenshot(item.URL, item.BufferedURL)
		resizeImage(item.BufferedURL)
	}

}
