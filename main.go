package main

import (
	"fmt"
	"os"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/labstack/gommon/color"
)

var arguments = struct {
	Input    string
	RandomUA bool
}{}

// Article struct hold data scraped from an article
type Article struct {
	// Basic informations
	Title   string
	Summary string

	// Body
	Body []string

	// Author
	Author Author
}

// Author struct hold data about an author
type Author struct {
	Name string
}

func scrapeArticle(articleLink string) (Article, error) {
	// Create an article structure
	var article Article
	// Create collector
	c := colly.NewCollector(
		colly.AllowedDomains("medium.com"),
	)

	// Randomize user agent on every request
	if arguments.RandomUA == true {
		extensions.RandomUserAgent(c)
	}

	// Scrape article's title and summary
	c.OnHTML("div.elevateCover", func(e *colly.HTMLElement) {
		// Scrape username
		article.Title = e.ChildText("h1.elevate-h1")

		// Scrape summary
		article.Summary = e.ChildText("p.elevate-summary")
	})

	// Scrape author informations
	c.OnHTML("div.u-flexEnd", func(e *colly.HTMLElement) {
		// Scrape author name
		article.Author.Name = e.ChildText("a.postMetaInline--author")
	})

	// Scrape text
	c.OnHTML("div.postArticle-content", func(e *colly.HTMLElement) {
		e.ForEach("section.section--body", func(_ int, el *colly.HTMLElement) {
			el.ForEach("p.graf--p", func(_ int, em *colly.HTMLElement) {
				article.Body = append(article.Body, em.DOM.Text()+"\n")
				em.ForEach("li.graf--li", func(_ int, ej *colly.HTMLElement) {
					article.Body = append(article.Body, "• "+ej.DOM.Text())
				})
			})
			el.ForEach("li.graf--li", func(_ int, em *colly.HTMLElement) {
				article.Body = append(article.Body, "• "+em.DOM.Text())
				em.ForEach("p.graf--p", func(_ int, ej *colly.HTMLElement) {
					article.Body = append(article.Body, "• "+ej.DOM.Text())
				})
			})
		})
	})

	// Visit page and fill collector
	c.Visit(articleLink)

	return article, nil
}

func main() {
	// Parse arguments and fill the arguments structure
	parseArgs(os.Args)

	// Scrape the article
	article, err := scrapeArticle(arguments.Input)
	if err != nil {
		fmt.Println(color.Red("Error while scraping the article: ") + err.Error())
	}

	fmt.Println("Scraping: " + arguments.Input + "\n")
	fmt.Println("Title:   " + article.Title)
	fmt.Println("Summary: " + article.Summary + "\n")
	fmt.Println("Author name: " + article.Author.Name + "\n")
	fmt.Println("Body: ")
	for index := 0; index < len(article.Body); index++ {
		fmt.Println(article.Body[index])
	}

}
