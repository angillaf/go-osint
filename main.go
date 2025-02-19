package main

import (
	"github.com/fatih/color"
	"go-osint/keywords"
	"go-osint/username"
	"os"
)

func showUsage() {
	color.Red("Usage: go-osint [type] [term] (platform)")
	color.Yellow("- [type]: \"username\" or \"keywords\"")
	color.Yellow("- [term]: The [type] to search for")
	color.Yellow("- (platform): The platform to search for the username on (optional)")
}

func main() {
	if len(os.Args) < 3 {
		showUsage()
		return
	}

	searchType := os.Args[1]
	searchTerm := os.Args[2]
	var outputFile string

	var urls = []string{}

	// Check for output argument
	for i, arg := range os.Args {
		if (arg == "--output" || arg == "-o") && i+1 < len(os.Args) {
			outputFile = os.Args[i+1]
		}
	}

	switch searchType {
	case "username":
		if len(os.Args) == 4 {
			platform := os.Args[3]
			ptr := username.Search(searchTerm, platform)
			if ptr != nil {
				urls = append(urls, *ptr)
			}
		} else {
			ptr := username.SearchAll(searchTerm)
			if ptr != nil {
				urls = append(urls, *ptr...)
			}
		}
	case "keywords":
		ptr := keywords.Search(searchTerm)
		if ptr != nil {
			urls = append(urls, *ptr...)
			for i := range urls {
				color.Green("[✔] Found URL: %s", urls[i])
			}
		}

	default:
		showUsage()
	}

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			color.Red("Failed to create output file: %v", err)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				color.Red("Failed to close output file: %v", err)
			}
		}(file)

		for _, url := range urls {
			_, _ = file.WriteString(url + "\n")
		}
		color.Green("[✔] Wrote %d (found) URLs to %s", len(urls), outputFile)
	}
}
