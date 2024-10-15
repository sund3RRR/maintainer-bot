package fetcher

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sund3RRR/maintainer-bot/internal/adapters/db"
)

func formatTitle(repo *db.Repo, newTagName string) string {
	return fmt.Sprintf(
		"<b>%s</b> / <b>%s</b> <code>%s</code> -> <code>%s</code>",
		repo.Owner,
		repo.Repo,
		repo.LastTag,
		newTagName,
	)
}
func formatReleaseBody(text string) string {
	stripped := strings.Trim(text, "\n")
	splitted := strings.Split(stripped, "\n")

	var result string

	for _, line := range splitted {
		is_heading := false
		for len(line) > 0 && line[0] == '#' {
			line = strings.TrimPrefix(line, "#")
			line = strings.TrimSpace(line)
			is_heading = true
		}
		if is_heading {
			result += fmt.Sprintf("<b><u>%s</u></b>\n", line)
		} else {
			result += line + "\n"
		}

	}

	result = replaceMdToHtml(result, "`", "code")
	result = replaceMdToHtml(result, `\*\*`, "b")
	result = replaceMdToHtml(result, `\*`, "i")

	return result
}

func replaceMdToHtml(text, mdSymbol, htmlTag string) string {
	// Define a regular expression to find words enclosed in backticks
	pattern := fmt.Sprintf(`%s([^`+"`"+`]+)%s`, mdSymbol, mdSymbol)

	// Replace all occurrences of the pattern with <code>word</code>
	re := regexp.MustCompile(pattern)
	result := re.ReplaceAllString(text, fmt.Sprintf(`<%s>$1</%s>`, htmlTag, htmlTag))

	return result
}
