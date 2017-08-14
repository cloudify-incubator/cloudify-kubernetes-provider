package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func printBottomLine(columnSizes []int) {
	fmt.Printf("+")
	for _, size := range columnSizes {
		fmt.Print(strings.Repeat("-", size+2))
		fmt.Printf("+")
	}
	fmt.Printf("\n")
}

func printLine(columnSizes []int, lines []string) {
	fmt.Printf("|")
	for col, size := range columnSizes {
		fmt.Print(" " + lines[col] + " ")
		fmt.Print(strings.Repeat(" ", size-utf8.RuneCountInString(lines[col])))
		fmt.Printf("|")
	}
	fmt.Printf("\n")
}

func PrintTable(titles []string, lines [][]string) {
	columnSizes := make([]int, len(titles))

	// column title sizes
	for col, name := range titles {
		if columnSizes[col] < utf8.RuneCountInString(name) {
			columnSizes[col] = utf8.RuneCountInString(name)
		}
	}

	// column value sizes
	for _, values := range lines {
		for col, name := range values {
			if columnSizes[col] < utf8.RuneCountInString(name) {
				columnSizes[col] = utf8.RuneCountInString(name)
			}
		}
	}

	printBottomLine(columnSizes)
	// titles
	printLine(columnSizes, titles)
	printBottomLine(columnSizes)
	// lines
	for _, values := range lines {
		printLine(columnSizes, values)
	}
	printBottomLine(columnSizes)
}
