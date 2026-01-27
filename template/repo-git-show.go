package template

import (
	"bytes"
	"git.lewoof.xyz/clone/gitbrowse/config"
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
	t := template.Must(template.New("body").Parse(`
		<body class="repo-git-show">
			<header>
			 	<img src="{{.Config.Thumbnail}}" alt="Thumbnail">
				<div>
				<h1>{{.Config.Title}}</h1>
				<p>Clone URL: <code>{{.Config.CloneURL}}</code></p>
				<table>
					<tr>
					<td><a href="{{.Config.URLRoot}}/">Readme</a></td>
					<td><a href="{{.Config.URLRoot}}/branch/master/tree">Tree</a></td>
					<td><em><a href="{{.Config.URLRoot}}/branch/master/commit">Commits</a></em></td>
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
	gitShow := GitShow(p.Config.RootDir, p.Commit)
	body = bodyBuffer.String() + gitShow + "</article></main></body>"
	return
}

func (p RepoGitShowPage) FullPage() string {
	return "<!DOCTYPE html><html>" + p.Head() + p.Body() + "</html>"
}
