/*
Copyright (c) 2017 GigaSpaces Technologies Ltd. All rights reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

// return clean list of arguments and options
func CliArgumentsList(osArgs []string) (arguments []string, options []string) {
	for pos, str := range osArgs {
		if str[:1] == "-" {
			return osArgs[:pos], osArgs[pos:]
		}
	}
	return osArgs, []string{}
}
