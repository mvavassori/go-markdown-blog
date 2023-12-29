# Minimalist Markdown Blog

## Overview

This project presents a minimalist blog platform built using only the Go standard library for HTML and CSS rendering, with a focus on simplicity and efficiency. The blog is designed to cater to minimalists who appreciate a clean, straightforward approach to web design and functionality.

## Features

- **No JavaScript**: The entire blog operates without any JavaScript, emphasizing speed, simplicity, and enhanced security.
- **Go Standard Library**: Built using the powerful and efficient Go standard library, ensuring robust performance and stability.
- **Markdown Support**: Utilizes the lightweight [gomarkdown/markdown](https://github.com/gomarkdown/markdown) package for converting markdown to HTML, streamlining the content creation process.
- **Minimalist Design**: A clean and uncluttered interface that focuses on content, free from distractions.
- **Automatic Dark Mode**: When your system is set to dark mode, the blog will transform to dark mode.
- **"No Technologies Detected"**: When analyzed with tools like Wappalyzer, the site proudly shows "No technologies detected", reflecting its minimalist and straightforward tech stack.

## Getting Started

### Running Locally

1. Clone the repository:
   ```bash
   git clone https://your-repository-url.git
   ```
2. Navigate to the project directory:
   ```bash
   cd minimalist-markdown-blog
   ```
3. Run the Go application:
   ```bash
   go run .
   ```
   The blog will be accessible at http://localhost:8080.

### Building the Application

1. Compile the application:
   ```bash
   go build -o myblog
   ```
2. Run the Executable:
   ```
   ./myblog
   ```
   Access the blog at http://localhost:8080.

## Blog Structure

- `main.go`: The main entry point for the blog application.
- `/templates`: Contains HTML templates for rendering blog posts and pages.
- `/static`: Houses CSS files for styling and static assets.
- `/posts`: Markdown files for blog posts are stored here.

## Acknowledgments

This project uses the [gomarkdown/markdown](https://github.com/gomarkdown/markdown) package, which is distributed under the Simplified BSD License. Full license text is available [here](https://github.com/gomarkdown/markdown/blob/master/LICENSE.txt).

## License

This project is open-sourced under the MIT License.
