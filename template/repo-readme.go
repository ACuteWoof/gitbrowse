package template

import (
	"bytes"
	"html/template"
	"git.lewoof.xyz/gitbrowse/config"
)

type RepoReadmePage struct {
	Readme string;
	Config *config.PageConfig;
}

func (p RepoReadmePage) Head() (head string) {
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

func (p RepoReadmePage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	t := template.Must(template.New("body").Parse(`
		<body class="user">
			<main>
			<h1>{{.Config.Title}}</h1>
			<article>
	`))
	t.Execute(&bodyBuffer, p)
	body = bodyBuffer.String() + p.Readme + "</article></main></body>"
	return
}

func (p RepoReadmePage) FullPage() string {
	return "<!DOCTYPE html><html>" + p.Head() + p.Body() + "</html>"
}
