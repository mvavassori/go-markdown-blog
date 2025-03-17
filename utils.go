package main

import (
	"fmt"

	"os"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func ExtractFrontMatter(markdownContent string) (map[string]string, string) {
	frontMatterPattern := regexp.MustCompile(`(?s)---(.*?)---(.*)`)
	matches := frontMatterPattern.FindStringSubmatch(markdownContent)

	frontMatter := make(map[string]string)
	var content string

	if len(matches) == 3 {
		for _, line := range strings.Split(matches[1], "\n") {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				// Remove surrounding quotes
				value = strings.Trim(value, `"`) // Trims just double quotes
				frontMatter[key] = value
			}
		}
		content = matches[2]
	} else {
		content = markdownContent // No front matter, entire content is Markdown
	}

	return frontMatter, content
}

func MdToHTML(md []byte) []byte {
	// Create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	var printAst = false
	// For debugging, not necessary in production
	if printAst {
		fmt.Print("--- AST tree:\n")
		ast.Print(os.Stdout, doc)
		fmt.Print("\n")
	}

	// Create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
