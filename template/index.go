package template

import (
	"bytes"
	"html/template"

	"git.lewoof.xyz/clone/gitbrowse/config"
)

type IndexPage struct {
	Repos []string;
	Config *config.PageConfig;
}

func (p IndexPage) Head() (head string) {
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

func (p IndexPage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	t := template.Must(template.New("body").Parse(`
		<body class="user">
			<main>
			<h1>{{.Config.Title}}</h1>
			<ul>
				{{range .Repos}}
					<li><a href="{{.}}">{{.}}</a></li>
				{{end}}
			</ul>
			</main>
		</body>
	`))
	t.Execute(&bodyBuffer, p)
	body = bodyBuffer.String()
	return
}

func (p IndexPage) FullPage() string {
	return "<!DOCTYPE html><html>" + p.Head() + p.Body() + "</html>"
}
