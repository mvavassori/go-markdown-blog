package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// A struct to hold the front matter and content
type BlogPost struct {
	Title    string
	Excerpt  string
	Date     string
	Content  template.HTML
	Filename string
	ParsedDate time.Time
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Read all files in the 'posts' directory
	files, err := os.ReadDir("posts")
	if err != nil {
		log.Printf("Error reading posts directory: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Initialize an empty slice to hold the blog posts
	var posts []BlogPost
	// Loop through the files and extract front matter and content
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			// Read the Markdown file
			markdownBytes, err := os.ReadFile(filepath.Join("posts", file.Name()))
			if err != nil {
				log.Println(err)
				continue // Skip this file
			}

			// Extract front matter
			frontMatter, _ := extractFrontMatter(string(markdownBytes))

			// Parse the date from the front matter
            parsedDate, err := time.Parse("2006-01-02", frontMatter["date"])
            if err != nil {
                log.Printf("Error parsing date: %v", err)
                continue
            }

			// Create a new BlogPost struct and append it to the slice
			posts = append(posts, BlogPost{
				Title:    frontMatter["title"],
				Excerpt:  frontMatter["excerpt"],
				Date:     frontMatter["date"],
				Filename: strings.TrimSuffix(file.Name(), ".md"),
				ParsedDate: parsedDate,
			})
		}
	}

	// Sort the posts by date, newest first
    sort.Slice(posts, func(i, j int) bool {
        return posts[i].ParsedDate.After(posts[j].ParsedDate)
    })

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// Extract post name from url. It's a slice: will give you the part of the path after "/posts/"
	postSlug := r.URL.Path[len("/posts/"):]

	// Append '.md' to get the actual filename
	filename := postSlug + ".md"
	// Read the Markdown file
	markdownBytes, err := os.ReadFile(filepath.Join("posts", filename))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Extract front matter
	frontMatter, content := extractFrontMatter(string(markdownBytes))

	// Convert Markdown content to HTML
	htmlContent := mdToHTML([]byte(content))

	// Parse and execute the template with front matter and content
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := BlogPost{
		Title:   frontMatter["title"],
		Excerpt: frontMatter["excerpt"],
		Date:    frontMatter["date"],
		Content: template.HTML(htmlContent),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func extractFrontMatter(markdownContent string) (map[string]string, string) {
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

func mdToHTML(md []byte) []byte {
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

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// Serve about.html on the /about route
    http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/about.html")
    })
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/posts/", postHandler)
	fmt.Println("Starting server at port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
