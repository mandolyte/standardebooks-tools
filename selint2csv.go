package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// ConvertUnicodeTableToCSV converts the UTF-8 box-drawn table into a 5-column CSV:
// [Code, Severity, File, Line, Message] with empty cell propagation.
func ConvertUnicodeTableToCSV(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	var blockLines []string
	isHeader := true

	// State carried forward from previous rows
	var lastCode, lastSeverity, lastFile string

	flushBlock := func(lines []string) error {
		if len(lines) == 0 {
			return nil
		}

		// Split line into cells along all vertical box boundaries and vertical joints
		splitCols := func(line string) []string {
			return strings.FieldsFunc(line, func(r rune) bool {
				switch r {
				case '│', '┃', '|', '├', '┤', '┏', '┓', '┗', '┛', '┬', '┴', '┼', '┡', '╇', '┩', '┢', '┪':
					return true
				default:
					return false
				}
			})
		}

		firstParts := splitCols(lines[0])
		numCols := len(firstParts)
		if numCols < 4 {
			return nil
		}

		// Accumulate lines for each column
		colLines := make([][]string, numCols)
		for _, line := range lines {
			parts := splitCols(line)
			for col := 0; col < numCols; col++ {
				if col < len(parts) {
					colLines[col] = append(colLines[col], parts[col])
				} else {
					colLines[col] = append(colLines[col], "")
				}
			}
		}

		// --- Process Header Row ---
		if isHeader {
			isHeader = false
			header := []string{"Code", "Severity", "File", "Line", "Message"}
			return csvWriter.Write(header)
		}

		// --- Process Data Rows ---
		// 1. Code & Severity (strips any accidental border runes before checking empty)
		code := cleanCellContent(joinLines(colLines[0], false))
		severity := cleanCellContent(joinLines(colLines[1], false))

		if code == "" {
			code = lastCode
		} else {
			lastCode = code
		}

		if severity == "" {
			severity = lastSeverity
		} else {
			lastSeverity = severity
		}

		// 2. Separate Column 3 into File and Line
		fileName, lineNum := extractFileAndLine(colLines[2])
		fileName = cleanCellContent(fileName)
		lineNum = cleanCellContent(lineNum)

		if fileName == "" {
			fileName = lastFile
		} else {
			lastFile = fileName
		}

		// 3. Message (Preserves inner indentation for XML/code blocks)
		message := cleanCellContent(joinLines(colLines[3], true))

		// Discard purely empty rows
		if code == "" && severity == "" && fileName == "" && lineNum == "" && message == "" {
			return nil
		}

		return csvWriter.Write([]string{code, severity, fileName, lineNum, message})
	}

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// If this is a horizontal divider/boundary line, flush any accumulated row block
		if isHorizontalBorder(trimmed) {
			if len(blockLines) > 0 {
				if err := flushBlock(blockLines); err != nil {
					return err
				}
				blockLines = nil
			}
			continue
		}

		blockLines = append(blockLines, line)
	}

	if len(blockLines) > 0 {
		if err := flushBlock(blockLines); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// cleanCellContent trims whitespace and removes residual box-drawing characters.
func cleanCellContent(s string) string {
	borderRunes := "┏┳┓┗┻┛├┼┤┡╇┩━─-+|│┃└┴┘┌┬┐"
	cleaned := strings.TrimFunc(s, func(r rune) bool {
		return strings.ContainsRune(borderRunes, r) || r == ' ' || r == '\t' || r == '\r' || r == '\n'
	})
	return cleaned
}

// extractFileAndLine separates left-aligned filenames from right-aligned "Line X" indicators.
func extractFileAndLine(rawLines []string) (file string, line string) {
	for _, raw := range rawLines {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}

		// Check if right-aligned: significant leading whitespace before "Line ..."
		leadingSpaces := len(raw) - len(strings.TrimLeft(raw, " "))
		if leadingSpaces > 5 && strings.HasPrefix(trimmed, "Line") {
			line = trimmed
		} else {
			file = trimmed
		}
	}
	return file, line
}

// joinLines cleans and joins multiple lines within a single cell block.
func joinLines(lines []string, detectIndent bool) string {
	hasIndentation := false
	if detectIndent {
		for _, l := range lines {
			trimmedRight := strings.TrimRight(l, " \t\r\n")
			if strings.HasPrefix(trimmedRight, "  ") || strings.HasPrefix(trimmedRight, "\t") {
				hasIndentation = true
				break
			}
		}
	}

	var cleaned []string
	for _, l := range lines {
		if hasIndentation {
			l = strings.TrimRight(l, " \t\r\n")
		} else {
			l = strings.TrimSpace(l)
		}

		if l != "" {
			cleaned = append(cleaned, l)
		}
	}

	if hasIndentation {
		return strings.Join(cleaned, "\n")
	}
	return strings.Join(cleaned, " ")
}

// isHorizontalBorder checks if a line consists only of border/divider runes.
func isHorizontalBorder(s string) bool {
	borderRunes := "┏┳┓┗┻┛├┼┤┡╇┩━─-+│┃└┴┘┌┬┐"
	nonBorderCount := 0
	for _, r := range s {
		if !strings.ContainsRune(borderRunes, r) && r != ' ' {
			nonBorderCount++
		}
	}
	return nonBorderCount == 0 && utf8.RuneCountInString(s) > 1
}

func main() {
	inFile, err := os.Open("lint.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input: %v\n", err)
		return
	}
	defer inFile.Close()

	outFile, err := os.Create("lint.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output: %v\n", err)
		return
	}
	defer outFile.Close()

	if err := ConvertUnicodeTableToCSV(inFile, outFile); err != nil {
		fmt.Fprintf(os.Stderr, "Conversion failed: %v\n", err)
	}
}