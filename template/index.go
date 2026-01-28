package template

import (
	"bytes"
	"html/template"

	"git.lewoof.xyz/gitbrowse/config"
)

type IndexPage struct {
	Repos []string;
	Config *config.PageConfig;
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
	return "<!DOCTYPE html><html>" + CommonHead(p.Config) + p.Body() + "</html>"
}
