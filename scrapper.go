package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// defining a data structure to store the scraped data
type PokemonProduct struct {
	url, image, name, price string
}

// definy output file
const FILE_OUTPUT = "output.csv"

// define targer URL
const TARGET_URL = "https://scrapeme.live/shop/page/1/"

// it verifies if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {

	// reset output.csv if exists
	_, err := os.Stat(FILE_OUTPUT)

	if os.IsExist(err) {
		os.Remove(FILE_OUTPUT)
	}

	// initializing the slice of structs that will contain the scraped data
	var pokemonProducts []PokemonProduct

	// initializing the list of pages to scrape with an empty slice
	var pagesToScrape []string

	// the first pagination URL to scrape
	pageToScrape := TARGET_URL

	// initializing the list of pages discovered with a pageToScrape
	pagesDiscovered := []string{pageToScrape}

	// current iteration
	i := 1
	// max pages to scrape
	limit := 5

	// initializing a Colly instance
	c := colly.NewCollector()
	// setting a valid User-Agent header
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15"

	// iterating over the list of pagination links to implement the crawling logic
	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		// discovering a new page
		newPaginationLink := e.Attr("href")

		// if the page discovered is new
		if !contains(pagesToScrape, newPaginationLink) {
			// if the page discovered should be scraped
			if !contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
	})

	// scraping the product data
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		pokemonProduct := PokemonProduct{}

		pokemonProduct.url = e.ChildAttr("a", "href")
		pokemonProduct.image = e.ChildAttr("img", "src")
		pokemonProduct.name = e.ChildText("h2")
		pokemonProduct.price = e.ChildText(".price")

		pokemonProducts = append(pokemonProducts, pokemonProduct)
	})

	c.OnScraped(func(response *colly.Response) {
		// until there is still a page to scrape
		if len(pagesToScrape) != 0 && i < limit {
			// getting the current page to scrape and removing it from the list
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			// incrementing the iteration counter
			i++

			// visiting a new page
			c.Visit(pageToScrape)
		}
	})

	// visiting the first page
	c.Visit(pageToScrape)

	// opening the CSV file
	file, err := os.Create(FILE_OUTPUT)
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// initializing a file writer
	writer := csv.NewWriter(file)

	// defining the CSV headers
	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}
	// writing the column headers
	writer.Write(headers)

	// adding each Pokemon product to the CSV output file
	for _, pokemonProduct := range pokemonProducts {
		// converting a PokemonProduct to an array of strings
		record := []string{
			pokemonProduct.url,
			pokemonProduct.image,
			pokemonProduct.name,
			pokemonProduct.price,
		}

		// writing a new CSV record
		writer.Write(record)
	}
	defer writer.Flush()
}
