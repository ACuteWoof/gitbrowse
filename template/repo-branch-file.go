package template

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"

	"git.lewoof.xyz/clone/gitbrowse/config"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

type RepoBranchFilePage struct {
	Contents string
	FilePath string
	Branch   string
	IsBinary bool
	Config   *config.PageConfig
}

func (p RepoBranchFilePage) Head() (head string) {
	var headBuffer bytes.Buffer
	t := template.Must(template.New("head").Parse(`
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<title>{{.Title}}</title>
			<meta name="description" content="{{.Description}}">
			{{range .Styles}}
				<link rel="stylesheet" href="{{.}}">
			{{end}}
			<link rel="icon" href="{{.Favicon}}">
		</head>
	`))
	t.Execute(&headBuffer, *p.Config)
	head = headBuffer.String()
	return
}

func (p RepoBranchFilePage) Body() (body string) {
	var bodyBuffer bytes.Buffer

	bodyBuffer.WriteString(CommonHeader(p.Config, "Tree"))

	type Crumb struct {
		Name        string
		DisplayName string
		Root        *string
		Branch      *string
	}

	type Breadcrumb struct {
		Crumbs []Crumb
	}

	var breadcrumbsBuffer bytes.Buffer
	breadcrumbTemplate := template.Must(template.New("breadcrumb").Parse(`
	<table class="breadcrumbs">
		<tr>
		{{range .Crumbs}}
			<td><a href="{{.Root}}/branch/{{.Branch}}/tree/{{.Name}}">{{.DisplayName}}</a></td>
		{{end}}
		</tr>
	</table>
	`))
	defaultCrumbs := []Crumb{
		{"", "/", &p.Config.URLRoot, &p.Branch},
	}
	if p.FilePath != "" {
		for entry := range strings.SplitSeq(p.FilePath, "/") {
			defaultCrumbs = append(defaultCrumbs, Crumb{entry, entry, &p.Config.URLRoot, &p.Branch})
		}
	}
	breadcrumbTemplate.Execute(&breadcrumbsBuffer, Breadcrumb{defaultCrumbs})
	breadcrumbs := breadcrumbsBuffer.String()

	descTemplate := template.Must(template.New("desc").Parse(`
		<p class="description">
			{{.}}
		</p>
		`))
	descTemplate.Execute(&bodyBuffer, "Viewing file on branch "+p.Branch)

	bodyBuffer.WriteString(breadcrumbs)
	bodyBuffer.WriteString("<table class=\"code\">")
	bodyBuffer.WriteString("<tbody>")
	highlightedContents := getHighlightedHTML(defaultCrumbs[len(defaultCrumbs)-1].Name, p.Contents)
	for i, line := range strings.Split(highlightedContents, "\n") {
		bodyBuffer.WriteString("<tr id=\"" + strconv.Itoa(i+1) + "\">")
		bodyBuffer.WriteString("<td class=\"linenumber\">")
		bodyBuffer.WriteString("<a href=\"#" + strconv.Itoa(i+1) + "\">")
		bodyBuffer.WriteString(strconv.Itoa((i + 1)))
		bodyBuffer.WriteString("</a>")
		bodyBuffer.WriteString("</td>")
		bodyBuffer.WriteString("<td>")
		bodyBuffer.WriteString(line)
		bodyBuffer.WriteString("</td>")
		bodyBuffer.WriteString("</tr>")
	}
	bodyBuffer.WriteString("<tbody>")
	bodyBuffer.WriteString("<table>")

	body = bodyBuffer.String()

	return
}

func (p RepoBranchFilePage) FullPage() string {
	return "<!DOCTYPE html><html>" + p.Head() + p.Body() + "</html>"
}

func getHighlightedHTML(filename string, contents string) string {
	var w bytes.Buffer
	var lexer chroma.Lexer
	lexer = lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Analyse(contents)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}

	lexer = chroma.Coalesce(lexer)
	style := styles.Get("gruvbox")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.Standalone(false), html.WithClasses(false), html.PreventSurroundingPre(true))
	iterator, err := lexer.Tokenise(nil, string(contents))
	checkErr(err)
	err = formatter.Format(&w, style, iterator)
	checkErr(err)
	r := strings.TrimPrefix(strings.TrimSuffix(w.String(), "</code></pre>"), "<pre><code>")
	return r
}
