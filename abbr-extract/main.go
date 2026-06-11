package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/net/html"
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.xhtml"))
	if err != nil {
		log.Fatalf("Error finding files: %v", err)
	}

	for _, file := range files {
		fmt.Printf("\n--- Processing: %s ---\n", file)
		processFile(file)
	}
}

func processFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer f.Close()

	tokenizer := html.NewTokenizer(f)
	inAbbr := false

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			// End of file
			return
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "abbr" {
				inAbbr = true
				fmt.Print("<abbr")
				for _, attr := range token.Attr {
					fmt.Printf(` %s="%s"`, attr.Key, attr.Val)
				}
				fmt.Print(">")
			}
		case html.EndTagToken:
			token := tokenizer.Token()
			if token.Data == "abbr" {
				fmt.Println("</abbr>")
				inAbbr = false
			}
		case html.TextToken:
			if inAbbr {
				// tokenizer.Text() returns the raw bytes of the current token
				fmt.Print(string(tokenizer.Text()))
			}
		}
	}
}