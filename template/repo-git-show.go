package template

import (
	"bytes"
	"git.lewoof.xyz/gitbrowse/config"
	"html/template"
)

type RepoGitShowPage struct {
	Commit string
	Config *config.PageConfig
}

func (p RepoGitShowPage) Head() (head string) {
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

func (p RepoGitShowPage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	bodyBuffer.WriteString(CommonHeader(p.Config, "Git Show"))
	gitShow := GitShow(p.Config.RootDir, p.Commit)
	body = bodyBuffer.String() + gitShow + "</article></main></body>"
	return
}

func (p RepoGitShowPage) FullPage() string {
	return "<!DOCTYPE html><html>" + p.Head() + p.Body() + "</html>"
}
