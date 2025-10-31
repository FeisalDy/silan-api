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

// ExtractBodyContent extracts only the body content from XHTML and removes all class attributes
func ExtractBodyContent(htmlContent []byte) string {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return string(htmlContent) // Return original if parsing fails
	}

	// Find the body node
	var bodyNode *html.Node
	var findBody func(*html.Node)
	findBody = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			bodyNode = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findBody(c)
			if bodyNode != nil {
				return
			}
		}
	}
	findBody(doc)

	if bodyNode == nil {
		return string(htmlContent) // Return original if no body found
	}

	// Remove all class attributes from all nodes in the body
	var removeClasses func(*html.Node)
	removeClasses = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Remove class attribute
			newAttrs := []html.Attribute{}
			for _, attr := range n.Attr {
				if attr.Key != "class" {
					newAttrs = append(newAttrs, attr)
				}
			}
			n.Attr = newAttrs
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			removeClasses(c)
		}
	}
	removeClasses(bodyNode)

	// Render the body node's children (inner HTML of body)
	var buf bytes.Buffer
	for c := bodyNode.FirstChild; c != nil; c = c.NextSibling {
		if err := html.Render(&buf, c); err != nil {
			return string(htmlContent)
		}
	}

	return buf.String()
}
