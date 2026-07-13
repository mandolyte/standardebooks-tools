package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// wordRegex matches sequence of Unicode letters, apostrophes, and hyphens
var wordRegex = regexp.MustCompile(`[\pL'\-]+`)

// Replacer to normalize curly quotes to straight quotes
var quoteReplacer = strings.NewReplacer(
	"“", "\"",
	"”", "\"",
	"‘", "'",
	"’", "'",
)

// cleanParagraph extracts bare words from an HTML node tree and reunites them with spaces.
func cleanParagraph(n *html.Node) string {
	var words []string

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			// Convert to lowercase and normalize curly quotes to straight quotes
			cleanedData := quoteReplacer.Replace(strings.ToLower(node.Data))
			found := wordRegex.FindAllString(cleanedData, -1)
			words = append(words, found...)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return strings.Join(words, " ")
}

// extractParagraphs returns a slice of strings, where each element is the "bare words" version of a <p>
func extractParagraphs(xhtmlContent string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(xhtmlContent))
	if err != nil {
		return nil, err
	}

	var paragraphs []string

	var findParagraphs func(*html.Node)
	findParagraphs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			cleanedText := cleanParagraph(n)
			paragraphs = append(paragraphs, cleanedText)
			return 
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findParagraphs(c)
		}
	}

	findParagraphs(doc)
	return paragraphs, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <file1.xhtml> <file2.xhtml>")
		os.Exit(1)
	}

	file1Name := os.Args[1]
	file2Name := os.Args[2]

	bytes1, err := os.ReadFile(file1Name)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file1Name, err)
		os.Exit(1)
	}

	bytes2, err := os.ReadFile(file2Name)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file2Name, err)
		os.Exit(1)
	}

	paragraphs1, err := extractParagraphs(string(bytes1))
	if err != nil {
		fmt.Printf("Error parsing XHTML from %s: %v\n", file1Name, err)
		os.Exit(1)
	}

	paragraphs2, err := extractParagraphs(string(bytes2))
	if err != nil {
		fmt.Printf("Error parsing XHTML from %s: %v\n", file2Name, err)
		os.Exit(1)
	}

	len1 := len(paragraphs1)
	len2 := len(paragraphs2)

	if len1 != len2 {
		fmt.Println("Error: The files do not have the same number of paragraphs.")
		fmt.Printf("%s contains %d paragraph(s).\n", file1Name, len1)
		fmt.Printf("%s contains %d paragraph(s).\n", file2Name, len2)
		os.Exit(1)
	}

	for i := 0; i < len1; i++ {
		p1Text := paragraphs1[i]
		p2Text := paragraphs2[i]

		words1 := strings.Fields(p1Text)
		words2 := strings.Fields(p2Text)

		p1WordsCount := len(words1)
		p2WordsCount := len(words2)

		if p1WordsCount != p2WordsCount {
			// Find the index of the first differing word
			mismatchIdx := 0
			minLen := p1WordsCount
			if p2WordsCount < minLen {
				minLen = p2WordsCount
			}

			for j := 0; j < minLen; j++ {
				if words1[j] != words2[j] {
					mismatchIdx = j
					break
				}
				// If they match all the way up to the end of the shorter paragraph,
				// the mismatch is the very next word in the longer paragraph.
				if j == minLen-1 {
					mismatchIdx = minLen
				}
			}

			// Slice up to the mismatch (inclusive of the differing word)
			limit1 := mismatchIdx + 1
			if limit1 > p1WordsCount {
				limit1 = p1WordsCount
			}
			limit2 := mismatchIdx + 1
			if limit2 > p2WordsCount {
				limit2 = p2WordsCount
			}

			snippet1 := strings.Join(words1[:limit1], " ")
			snippet2 := strings.Join(words2[:limit2], " ")

			fmt.Printf("Mismatch found at paragraph index %d!\n", i)
			fmt.Println(strings.Repeat("-", 60))
			fmt.Printf("Filename: %s (Index: %d)\n", file1Name, i)
			fmt.Printf("Up to mismatch: \"%s\"\n\n", snippet1)
			fmt.Printf("Filename: %s (Index: %d)\n", file2Name, i)
			fmt.Printf("Up to mismatch: \"%s\"\n", snippet2)
			fmt.Println(strings.Repeat("-", 60))
			
			os.Exit(1)
		}
	}

	fmt.Println("Success! Both files have the same number of paragraphs, and each corresponding paragraph has the exact same word count.")
}