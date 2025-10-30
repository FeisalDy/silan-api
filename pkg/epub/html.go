package epub

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

// ExtractText extracts plain text from HTML content
func ExtractText(htmlContent []byte) string {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return ""
	}

	var sb strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				sb.WriteString(text)
				sb.WriteString(" ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return strings.TrimSpace(sb.String())
}
