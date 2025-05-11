package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func fetchDefinition(word string) (string, error) {
	url := fmt.Sprintf("https://merriam-webster.com/dictionary/%s", word)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func extractDefinitions(html string) []string {
	var definitions []string
	dtTextRegex := regexp.MustCompile(`<span class="dtText">(.*?)</span>`)
	matches := dtTextRegex.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(match) > 1 {
			strongRegex := regexp.MustCompile(`.*</strong>(.*)`)
			submatches := strongRegex.FindStringSubmatch(match[1])

			if len(submatches) > 1 {
				definition := strings.TrimSpace(submatches[1])
				definition = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(definition, "")
				definitions = append(definitions, fmt.Sprintf("-- %s", definition))
			} else {

				uppercaseRegex := regexp.MustCompile(`.*<span class="text-uppercase">([^<]*)</span>.*`)
				upperMatches := uppercaseRegex.FindStringSubmatch(match[1])
				if len(upperMatches) > 1 {
					definitions = append(definitions, fmt.Sprintf("-* %s", upperMatches[1]))
				}
			}
		}
	}

	return definitions
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dict <word>")
		os.Exit(1)
	}

	word := os.Args[1]
	html, err := fetchDefinition(word)
	if err != nil {
		fmt.Printf("Error fetching definition: %v\n", err)
		os.Exit(1)
	}

	definitions := extractDefinitions(html)
	for _, def := range definitions {
		fmt.Println(def)
	}
}
