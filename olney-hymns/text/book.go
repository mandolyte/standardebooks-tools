package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var romanNumerals = map[int]string{
	1:  "I",
	2:  "II",
	3:  "III",
	4:  "IV",
	5:  "V",
	6:  "VI",
	7:  "VII",
	8:  "VIII",
	9:  "IX",
	10: "X",
}

// Define the XML chunk as a constant raw string literal
const xmlHeader string = `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" epub:prefix="z3998: http://www.daisy.org/z3998/2012/vocab/structure/, se: https://standardebooks.org/vocab/1.0" xml:lang="en-GB">
<head>
	<title>$BOOKNUM$</title>
	<link href="../css/core.css" rel="stylesheet" type="text/css"/>
	<link href="../css/local.css" rel="stylesheet" type="text/css"/>
</head>
<body epub:type="bodymatter z3998:fiction">
<section id="$BOOKDASHNUM$" epub:type="part">
`

const xmlFooter string = `</section>
</body>
</html>
`

// AssertStringsMatchCaseInsensitive compares two strings.
// If they don't match (ignoring case), it triggers a fatal error.
func AssertStringsMatchCaseInsensitive(str1, str2 string) {
	// EqualFold handles Unicode case-folding and is more efficient than strings.ToLower
	if !strings.EqualFold(str1, str2) {
		fmt.Println("--- DEBUG INFO: STRING MISMATCH ---")
		fmt.Printf("String 1: [%s]\n", str1)
		fmt.Printf("String 2: [%s]\n", str2)
		fmt.Println("-----------------------------------")

		log.Fatalf("Fatal Error: The provided strings do not match (case-insensitive check).")
	}
}

// SplitByLastEmDash finds the final em-dash (—) in a string and
// returns everything to the left and everything to the right.
func SplitByLastEmDash(line string) (string, string) {
	delimiter := "—"

	// 1. Find the index of the last occurrence of the em-dash
	lastIndex := strings.LastIndex(line, delimiter)
	fmt.Printf("Last Index of em-dash: %d\n", lastIndex)

	// 2. If the delimiter isn't found, LastIndex returns -1
	if lastIndex == -1 {
		return line, ""
		// fmt.Println("--- DEBUG INFO: DELIMITER MISSING ---")
		// fmt.Printf("Line Content: [%s]\n", line)
		// fmt.Println("-------------------------------------")
		// log.Fatal("Fatal Error: The required em-dash (—) was not found in the line.")
	}

	// 3. Extract the substrings based on the index
	// left is everything before the dash
	left := strings.TrimSpace(line[:lastIndex])

	// right is everything after the dash (skip the length of the dash itself)
	// An em-dash is a multi-byte UTF-8 character, so we use len(delimiter)
	right := strings.TrimSpace(line[lastIndex+len(delimiter):])

	return left, right
}

// ReadCSVToRecords takes a filename and returns a 2D slice of strings.
// It treats any error as fatal.
func ReadCSVToRecords(fileName string) [][]string {
	// 1. Open the file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error: Could not open CSV file %s: %v", fileName, err)
	}
	// Ensure the file is closed when the function finishes
	defer file.Close()

	// 2. Initialize the CSV reader
	reader := csv.NewReader(file)

	// 3. Read all records into a 2D slice ([][]string)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error: Failed to parse CSV file %s: %v", fileName, err)
	}

	return records
}

func processAndCount(input string) ([]int, int) {
	// 1. Clean the string
	cleaned := strings.TrimSuffix(input, "a")
	cleaned = strings.TrimSuffix(cleaned, "D")

	// 2. Split into segments
	segments := strings.Split(cleaned, ",")

	var result []int
	for _, s := range segments {
		// 3. Convert and check for errors
		val, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			// This "makes it fatal": prints the error and stops the program
			log.Fatalf("FATAL ERROR: Could not parse %q as an integer. Original input: %q. Error: %v", s, input, err)
		}
		result = append(result, val)
	}

	return result, len(result)
}

// Warning Global Variable
// Global variables defined at the package level
var (
	ofile         *os.File
	oerr          error
	baseFileName  string
	originalMeter string
	bookNum       string
)

func fileWriter(line string) {
	_, oerr = ofile.WriteString(line)
	if oerr != nil {
		log.Fatalf("FATAL: Could not write to file book-%v.txt: %v\n", bookNum, oerr)
	}
}

func main() {
	// 1. Check for command line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run book.go <base filename>")
		return
	}
	baseFileName = os.Args[1]
	bookNum = string(baseFileName[len(baseFileName)-1])

	// 1. Open (or create) the file for writing
	// os.Create creates the file if it doesn't exist, or truncates it if it does.
	ofilename := baseFileName + ".txt"
	ofile, oerr = os.Create(ofilename)
	if oerr != nil {
		log.Fatalf("FATAL: Could not create file book-%v.txt: %v", bookNum, oerr)
	}
	defer ofile.Close()

	csvFilename := "../olney-" + baseFileName + ".csv"
	records := ReadCSVToRecords(csvFilename)
	meterCol := 20
	rowIndex := 1
	cowperFlagCol := 15
	const titleCol = 0

	// 2. Output the XML chunk first
	//    a. first update the title for book
	//    b. second, update the section
	rnum := ""
	if bookNum == "1" {
		rnum = "I"
	} else if bookNum == "2" {
		rnum = "II"
	} else if bookNum == "3" {
		rnum = "III"
	} else {
		log.Fatalf("FATAL: bookNum must be 1, 2, or 3; was: %v", bookNum)
	}
	_bookTitle := "book " + rnum
	_bookTitle = strings.ToTitle(_bookTitle)
	// _bookTitle = strings.Replace(_bookTitle, "-", " ", 1)
	_xmlHeader := strings.Replace(xmlHeader, "$BOOKNUM$", _bookTitle, 1)
	_xmlHeader = strings.Replace(_xmlHeader, "$BOOKDASHNUM$", baseFileName, 1)
	fileWriter((_xmlHeader))

	// 3. Open the file
	fileName := baseFileName + ".xhtml"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	// 4. Read the file into an array of strings
	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	// 5. Loop that outputs each element to standard output
	start_hymn := false
	hymn_number := ""
	hymn_title_line := ""
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if line == "<h>" {
			start_hymn = true
			continue
		}
		if line == "</h>" {
			start_hymn = false
			continue
		}
		if start_hymn {
			if hymn_number == "" {
				hymn_number = line
				continue
			}
			if hymn_title_line == "" {
				hymn_title_line = line
				// continue
				start_hymn = false
			}
			// If both hymn_number and hymn_title are set, output them
			title, scriptureRef := SplitByLastEmDash(hymn_title_line)
			// fix straight quotes in title
			title = strings.ReplaceAll(title, "'", "’")
			title = strings.TrimSuffix(title, ".")
			fmt.Printf("Book: %s, %s, %s\n", hymn_number, title, scriptureRef)
			csvTitle := records[rowIndex][titleCol]
			csvTitle = strings.ReplaceAll(csvTitle, "'", "’")
			cowperFlag := records[rowIndex][cowperFlagCol]
			if cowperFlag == "No" || cowperFlag == "Yes" {
				// then all is well
			} else {
				log.Fatalf("FATAL ERROR\n: Cowper flag not No or Yes: %v\n", cowperFlag)
			}
			fmt.Printf("CSV: %s, %s\n", records[rowIndex][meterCol], records[rowIndex][titleCol])
			if baseFileName == "book-3" {
				switch hymn_number {
				case "89.", "90.", "91.", "92.", "94.", "95.", "98.", "102.", "103.", "104.", "105.", "106.", "107.":
					// don't do the title compare
					// Now process the stanza lines
					// Do this in a function using the current index i
					originalMeter = records[rowIndex][meterCol]
					process_stanzas(cowperFlag, scriptureRef, title+".", hymn_number, i, lines, records[rowIndex][meterCol])
				default:
					AssertStringsMatchCaseInsensitive(title, csvTitle)
					// Now process the stanza lines
					// Do this in a function using the current index i
					originalMeter = records[rowIndex][meterCol]
					process_stanzas(cowperFlag, scriptureRef, title+".", hymn_number, i+1, lines, records[rowIndex][meterCol])
				}
			} else {
				AssertStringsMatchCaseInsensitive(title, csvTitle)
				// Now process the stanza lines
				// Do this in a function using the current index i
				originalMeter = records[rowIndex][meterCol]
				process_stanzas(cowperFlag, scriptureRef, title+".", hymn_number, i+1, lines, records[rowIndex][meterCol])
			}
			rowIndex++
			hymn_number = ""
			hymn_title_line = ""
		}
		//fmt.Println(line)
	}

	// 6. Output the XML closing elements
	fileWriter(xmlFooter)
}

func verify_stanza_data(hymn_number string, expected_line_count int, stanza_line_count int, which_stanza int) {
	// Doubled hymns have twice the number of lines per stanza
	doubled_hymns_book1 := map[string]bool{
		"7.":   true,
		"8.":   true,
		"9.":   true,
		"12.":  true,
		"28.":  true,
		"34.":  true,
		"35.":  true,
		"37.":  true,
		"46.":  true,
		"52.":  true,
		"60.":  true,
		"61.":  true,
		"62.":  true,
		"63.":  true,
		"65.":  true,
		"89.":  true,
		"92.":  true,
		"93.":  true,
		"95.":  true,
		"114.": true,
		"117.": true,
		"123.": true,
		"127.": true,
	}

	doubled_hymns_book2 := map[string]bool{
		"41.": true,
		"67.": true,
		"79.": true,
	}

	doubled_hymns_book3 := map[string]bool{
		"1.":   true,
		"3.":   true,
		"4.":   true,
		"6.":   true,
		"9.":   true,
		"14.":  true,
		"16.":  true,
		"25.":  true,
		"30.":  true,
		"32.":  true,
		"37.":  true,
		"48.":  true,
		"54.":  true,
		"66.":  true,
		"75.":  true,
		"86.":  true,
		"89.":  true,
		"97.":  true,
		"98.":  true,
		"101.": true,
		"102.": true,
	}

	if baseFileName == "book-1" && doubled_hymns_book1[hymn_number] {
		stanza_line_count /= 2
	} else if baseFileName == "book-2" && doubled_hymns_book2[hymn_number] {
		stanza_line_count /= 2
	} else if baseFileName == "book-3" && doubled_hymns_book3[hymn_number] {
		stanza_line_count /= 2
	} else if strings.HasSuffix(originalMeter, "D") {
		stanza_line_count /= 2
	}
	if expected_line_count != stanza_line_count {
		log.Fatalf("FATAL ERROR\n: Stanza %d line count %d does not match expected meter line count %d",
			which_stanza, stanza_line_count, expected_line_count)
	}
}

func printStanzaHeader(hymn_number string, stanza_number int) {
	h := strings.TrimSuffix(hymn_number, ".")
	fileWriter(fmt.Sprintf("<section id=\"stanza-%v-%v-%d>\n", bookNum, h, stanza_number))
	fileWriter(" 	<header>\n")
	fileWriter(fmt.Sprintf(" 		<p>%v</p>\n", romanNumerals[stanza_number]))
	fileWriter(" 	</header>\n")
	fileWriter(" 	<p>\n")
}

func printStanzaFooter() {
	fileWriter(" 	</p>\n")
	fileWriter(" 	</section>\n")
}

func process_stanzas(cowperFlag, reference, hymn_title, hymn_number string, startIndex int, lines []string, meter string) {
	meterArray, meterCount := processAndCount(meter)
	fmt.Printf("Meter %s has %v: %v\n", meter, meterCount, meterArray)
	author := ""
	if cowperFlag == "Yes" {
		author = "Cowper"
	} else {
		author = "Newton"
	}
	// Print poem section heading lines
	fileWriter(fmt.Sprintf("<section id=\"hymn-%v-%v\" epub:type=\"z3998:hymn\">\n", bookNum, hymn_number))
	fileWriter("  <header>\n")
	fileWriter("  <hgroup>\n")
	fileWriter(fmt.Sprintf("    <h3 epub:type=\"ordinal\">%v</h3>\n", hymn_number))
	fileWriter(fmt.Sprintf("    <p epub:type=\"title\">%v</p>\n", hymn_title))
	fileWriter("  </hgroup>\n")
	fileWriter(fmt.Sprintf("  <p epub:type=\"z3998:contributors\">By %v</p>\n", author))
	fileWriter(fmt.Sprintf("  <p epub:type=\"bridgehead\">%v</p>\n", reference))
	fileWriter("  </header>\n")

	// stanza header for first stanza
	printStanzaHeader(hymn_number, 1)

	stanza_count := 0
	stanza_line_count := 0
	for i := startIndex + 1; i < len(lines); i++ {
		// debug
		fmt.Printf("Processing line %d: %s\n", i, lines[i])

		if lines[i] == "</h>" {
			// output laststanza data
			stanza_count++
			verify_stanza_data(hymn_number, meterCount, stanza_line_count, stanza_count)
			printStanzaFooter()
			stanza_line_count = -1
			break
		}
		if lines[i] == "" {
			// output stanza data
			stanza_count++
			verify_stanza_data(hymn_number, meterCount, stanza_line_count, stanza_count)
			printStanzaFooter()
			printStanzaHeader(hymn_number, stanza_count+1)
			stanza_line_count = 0
		} else {
			stanza_line_count++
			printStanzaLine(stanza_line_count, lines[i], meter)
		}
	}
	fmt.Printf("Number of stanzas: %v\n", stanza_count)
}

func printStanzaLine(stanza_line_count int, line string, meter string) {
	smeter := strings.TrimSuffix(meter, "a")
	smeter = strings.TrimSuffix(smeter, "D")
	switch smeter {
	case "8,6,8,6":
		isEven := stanza_line_count%2 == 0
		if isEven {
			line_1_indent(line)
		} else {
			line_no_indent(line)
		}
	case "8,8,8,8":
		line_no_indent(line)
	case "7,7,7,7":
		line_no_indent(line)
	case "10,10,11,11":
		line_no_indent(line)
	case "6,6,6,6,8,8":
		if stanza_line_count <= 4 {
			line_1_indent(line)
		} else {
			line_no_indent(line)
		}
	case "6,6,8,6":
		if stanza_line_count == 3 {
			line_no_indent(line)
		} else {
			line_1_indent(line)
		}
	case "7,6,7,6":
		isEven := stanza_line_count%2 == 0
		if isEven {
			line_1_indent(line)
		} else {
			line_no_indent(line)
		}
	case "7,7,7,7,7,7":
		line_no_indent(line)
	case "8,7,8,7,7,7":
		if stanza_line_count == 1 || stanza_line_count == 3 {
			line_no_indent(line)
		} else {
			line_1_indent(line)
		}
	case "8,7,8,7":
		line_no_indent(line)
	case "8,7,8,7,11":
		if stanza_line_count == 1 || stanza_line_count == 3 {
			line_1_indent(line)
		} else if stanza_line_count == 2 || stanza_line_count == 4 {
			line_2_indent(line)
		} else {
			line_no_indent(line)
		}
	case "7,6,7,6,7,7,7,6":
		switch stanza_line_count {
		case 1, 3, 5, 6, 7:
			line_no_indent(line)
		case 2, 4, 8:
			line_1_indent(line)
		}
	case "8,8,6,8,8,6":
		switch stanza_line_count {
		case 1, 2, 4, 5:
			line_no_indent(line)
		case 3, 6:
			line_1_indent(line)
		}
	case "8,8,8,8,8,8":
		line_no_indent(line)
	case "10,10,10,10":
		line_no_indent(line)
	case "8,8,8":
		line_no_indent(line)
	case "7,6,7,6,7,7":
		line_no_indent(line)
	case "6,6,6,6,7,7":
		switch stanza_line_count {
		case 1, 2, 3, 4:
			line_1_indent(line)
		case 6, 7:
			line_no_indent(line)
		}
	default:
		// This stops the program and prints the offending value
		log.Fatalf("FATAL: Unsupported meter format received: %q. Check your input data.", meter)
	}
}

func line_no_indent(line string) {
	fileWriter(fmt.Sprintf("    <span>%v</span>\n", line))
}

func line_1_indent(line string) {
	fileWriter(fmt.Sprintf("      <span class=\"i1\">%v</span>\n", line))
}

func line_2_indent(line string) {
	fileWriter(fmt.Sprintf("      <span class=\"i2\">%v</span>\n", line))
}
