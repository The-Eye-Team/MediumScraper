package main

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	nurl "net/url"
	"os"
	"strings"
	"time"

	readability "github.com/RadhiFadlillah/go-readability"
	"github.com/kennygrant/sanitize"
)

//var checkPre = color.Yellow("[") + color.Green("✓") + color.Yellow("]")
//var crossPre = color.Yellow("[") + color.Red("✗") + color.Yellow("]")

var client = http.Client{}

func init() {
	// Disable HTTP/2: Empty TLSNextProto map
	client.Transport = http.DefaultTransport
	client.Transport.(*http.Transport).TLSNextProto =
		make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)
}

func downloadCover(url, path string, client *http.Client) error {
	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		downloadCover(url, path, client)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func scrapeArticle(url string) error {
	// Create URL
	parsedURL, _ := nurl.Parse(url)

	// Fetch readable content
	article, err := readability.FromURL(parsedURL.String(), 5*time.Second)
	if err != nil {
		return err
	}

	// Show results
	title := strings.Split(article.Title, " –")

	text := []byte("Title: " + sanitize.Path(title[0]) + "\n")
	text = append(text, "Image: "+article.Image+"\n"...)
	text = append(text, "Author: "+article.Byline+"\n"...)
	text = append(text, "Excerpt: "+article.Excerpt+"\n"...)
	text = append(text, "\nContent: \n\n"+article.Content+"\n"...)

	author := strings.Replace(article.Byline, " ", "_", -1)

	// Create destination pth
	os.MkdirAll(author+"/"+sanitize.Path(title[0])+"/", os.ModePerm)

	// Write the article to file
	err = ioutil.WriteFile(author+"/"+sanitize.Path(title[0])+"/"+sanitize.Path(title[0])+".txt", text, 0644)
	if err != nil {
		return err
	}

	// Download cover
	err = downloadCover(article.Image, author+"/"+sanitize.Path(title[0])+"/"+sanitize.Path(title[0])+".jpg", &client)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := scrapeArticle(os.Args[1])
	if err != nil {
		scrapeArticle(os.Args[1])
	}
}
