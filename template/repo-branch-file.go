package template

import (
	"bytes"
	"html"
	"html/template"
	"strconv"
	"strings"

	"git.lewoof.xyz/gitbrowse/config"
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
	t := template.Must(template.New("body").Parse(`
		<body class="repo-branch-file">
			<header>
			 	<img src="{{.Config.Thumbnail}}" alt="Thumbnail">
				<div>
				<h1>{{.Config.Title}}</h1>
				<p>Clone URL: <code>{{.Config.CloneURL}}</code></p>
				<table>
					<tr>
					<td><em><a href="{{.Config.URLRoot}}/">Readme</a></em></td>
					<td><a href="{{.Config.URLRoot}}/branch/master/tree">Tree</a></td>
					<td><a href="{{.Config.URLRoot}}/branch/master/commit">Commits</a></td>
					<td><a href="{{.Config.URLRoot}}/branch">Branches</a></td>
					<td><a href="{{.Config.URLRoot}}/tag">Tags</a></td>
					</tr>
				</table>
				</div>
			</header>
			<main>
			<article>
	`))
	t.Execute(&bodyBuffer, p)

	type Crumb struct {
		Name   string
		Root   *string
		Branch *string
	}

	type Breadcrumb struct {
		Crumbs []Crumb
	}

	var breadcrumbsBuffer bytes.Buffer
	breadcrumbTemplate := template.Must(template.New("breadcrumb").Parse(`
	<table class="breadcrumbs">
		<tr>
		{{range .Crumbs}}
			<td><a href="{{.Root}}/branch/{{.Branch}}/tree/{{.Name}}">{{.Name}}</a></td>
		{{end}}
		</tr>
	</table>
	`))
	defaultCrumbs := []Crumb{
		{"/", &p.Config.URLRoot, &p.Branch},
	}
	if p.FilePath != "" {
		for entry := range strings.SplitSeq(p.FilePath, "/") {
			defaultCrumbs = append(defaultCrumbs, Crumb{entry, &p.Config.URLRoot, &p.Branch})
		}
	}
	breadcrumbTemplate.Execute(&breadcrumbsBuffer, Breadcrumb{defaultCrumbs})
	breadcrumbs := breadcrumbsBuffer.String()

	descTemplate := template.Must(template.New("desc").Parse(`
		<p class="description">
			{{.}}
		</p>
		`))
	descTemplate.Execute(&bodyBuffer, "Browsing file on branch "+p.Branch)

	bodyBuffer.WriteString(breadcrumbs)
	bodyBuffer.WriteString("<table class=\"code\">")
	bodyBuffer.WriteString("<tbody>")
	for i, line := range strings.Split(p.Contents, "\n") {
		bodyBuffer.WriteString("<tr id=\"" + strconv.Itoa(i + 1) + "\">")
		bodyBuffer.WriteString("<td class=\"linenumber\">")
		bodyBuffer.WriteString("<a href=\"#" + strconv.Itoa(i + 1) + "\">")
		bodyBuffer.WriteString(strconv.Itoa((i + 1)))
		bodyBuffer.WriteString("</a>")
		bodyBuffer.WriteString("</td>")
		bodyBuffer.WriteString("<td>")
		bodyBuffer.WriteString(html.EscapeString(line))
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
