package template

import (
	"bytes"
	"html/template"
	"os/exec"
	"strings"

	"git.lewoof.xyz/gitbrowse/config"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type RepoTagsPage struct {
	Repo      *git.Repository
	Config    *config.PageConfig
}

func (p RepoTagsPage) Head() (head string) {
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

func (p RepoTagsPage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	t := template.Must(template.New("body").Parse(`
		<body class="repo-tags">
			<header>
			 	<img src="{{.Config.Thumbnail}}" alt="Thumbnail">
				<div>
				<table>
				<h1>{{.Config.Title}}</h1>
					<tr>
					<td><a href="{{.Config.URLRoot}}/">Readme</a></td>
					<td><a href="{{.Config.URLRoot}}/branch/master/tree">Tree</a></td>
					<td><a href="{{.Config.URLRoot}}/branch/master/commit">Commits</a></td>
					<td><em><a href="{{.Config.URLRoot}}/branch">Branches</a></em></td>
					<td><a href="{{.Config.URLRoot}}/tag">Tags</a></td>
					</tr>
				</table>
				</div>
			</header>
			<main>
			<article>
	`))
	t.Execute(&bodyBuffer, p)

	tags, err := p.Repo.Tags()
	checkErr(err)

	type Row struct {
		URLRoot *string
		Branch  string
		Tag     *object.Tag
		ShortHash string
	}

	rows := []string{}

	tags.ForEach(func(r *plumbing.Reference) error {
		t, err := p.Repo.TagObject(r.Hash())
		if err == plumbing.ErrObjectNotFound {
			return nil
		}
		checkErr(err)
		var rowBuffer bytes.Buffer
		rowTemplate := template.Must(template.New("row").Parse(`<tr>
<td class="tag"><a href="{{.URLRoot}}/tag/{{.Tag.Hash.String}}">{{.Tag.Name}}</a></td>
<td class="taghash"><a href="{{.URLRoot}}/tag/{{.Tag.Hash.String}}">{{.ShortHash}}</a></td>
<td class="tagmessage">{{.Tag.Message}}</td>
<td class="tagger">
	<a href="mailto:{{.Tag.Tagger.Email}}">
	{{.Tag.Tagger.Name}}
	</a>
</td>
<td class="date">{{.Tag.Tagger.When.Format "2006-01-02 15:04:05"}}</td>
</tr>`))
		shortHash, err := exec.Command("git", "rev-parse", "--short", t.Hash.String()).Output()
		rowTemplate.Execute(&rowBuffer, Row{&p.Config.URLRoot, r.Name().Short(), t, string(shortHash)})
		rows = append(rows, rowBuffer.String())
		return nil
	})

	tableHeader := "<tr><th>Name</th><th>Hash</th><th>Message</th><th>Tagger</th><th>Date</th></tr>"

	table := "<table>" + tableHeader + strings.Join(rows, "") + "</table>"

	descTemplate := template.Must(template.New("desc").Parse(`
		<p class="description">
			{{.}}
		</p>
		`))


	descTemplate.Execute(&bodyBuffer, "Showing tags for repository")

	body = bodyBuffer.String() +
		table + "</article></main></body>"

	return body
}

func (p RepoTagsPage) FullPage() string {
	return "<!DOCTYPE html><html>" + p.Head() + p.Body() + "</html>"
}
