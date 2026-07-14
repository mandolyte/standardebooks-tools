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

	// Determine the maximum bounds to check
	maxLen := len1
	if len2 > maxLen {
		maxLen = len2
	}

	for i := 0; i < maxLen; i++ {
		// Case A: Both files have a paragraph at this index -> Compare their content
		if i < len1 && i < len2 {
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

				fmt.Printf("❌ Mismatch found at paragraph index %d!\n", i)
				fmt.Println(strings.Repeat("-", 60))
				fmt.Printf("Filename: %s (Index: %d)\n", file1Name, i)
				fmt.Printf("Up to mismatch: \"%s\" (%d words)\n\n", snippet1, p1WordsCount)
				fmt.Printf("Filename: %s (Index: %d)\n", file2Name, i)
				fmt.Printf("Up to mismatch: \"%s\" (%d words)\n", snippet2, p2WordsCount)
				fmt.Println(strings.Repeat("-", 60))
				os.Exit(1)
			}
		} else if i >= len1 {
			// Case B: File 2 has an extra paragraph that File 1 does not have
			fmt.Printf("❌ Mismatch found at paragraph index %d! (Paragraph count mismatch)\n", i)
			fmt.Println(strings.Repeat("-", 60))
			fmt.Printf("Filename: %s (Index: %d)\n", file1Name, i)
			fmt.Println("[No corresponding paragraph exists — File 1 ended early]")
			fmt.Println()
			fmt.Printf("Filename: %s (Index: %d)\n", file2Name, i)
			fmt.Printf("Extra Paragraph Content: \"%s\"\n", paragraphs2[i])
			fmt.Println(strings.Repeat("-", 60))
			os.Exit(1)
		} else if i >= len2 {
			// Case C: File 1 has an extra paragraph that File 2 does not have
			fmt.Printf("❌ Mismatch found at paragraph index %d! (Paragraph count mismatch)\n", i)
			fmt.Println(strings.Repeat("-", 60))
			fmt.Printf("Filename: %s (Index: %d)\n", file1Name, i)
			fmt.Printf("Extra Paragraph Content: \"%s\"\n", paragraphs1[i])
			fmt.Println()
			fmt.Printf("Filename: %s (Index: %d)\n", file2Name, i)
			fmt.Println("[No corresponding paragraph exists — File 2 ended early]")
			fmt.Println(strings.Repeat("-", 60))
			os.Exit(1)
		}
	}

	fmt.Println("✅ Success! Both files have the same number of paragraphs, and each corresponding paragraph has the exact same word count.")
}