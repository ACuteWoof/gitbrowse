package template

import (
	"bytes"
	"html/template"
	"strings"

	"git.lewoof.xyz/gitbrowse/config"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type RepoBranchLogPage struct {
	Repo      *git.Repository
	Branch    string
	BranchRef *plumbing.Reference
	Config    *config.PageConfig
}

func (p RepoBranchLogPage) Head() (head string) {
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

func (p RepoBranchLogPage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	t := template.Must(template.New("body").Parse(`
		<body class="repo-branch-commits">
			<header>
			 	<img src="{{.Config.Thumbnail}}" alt="Thumbnail">
				<div>
				<h1>{{.Config.Title}}</h1>
				<nav>
					<a href="{{.Config.URLRoot}}/">Readme</a>
					<a href="{{.Config.URLRoot}}/branch/master/tree">Tree</a>
					<em><a href="{{.Config.URLRoot}}/branch/master/commit">Commits</a></em>
					<a href="{{.Config.URLRoot}}/branch">Branches</a>
					<a href="{{.Config.URLRoot}}/tag">Tags</a>
				</nav>
				</div>
			</header>
			<main>
			<article>
	`))
	t.Execute(&bodyBuffer, p)

	commits, err := p.Repo.Log(&git.LogOptions{From: p.BranchRef.Hash(), Order: git.LogOrderCommitterTime})
	checkErr(err)

	type Row struct {
		URLRoot *string
		Branch  *string
		Commit  *object.Commit
		ShortHash string
	}

	rows := []string{}

	commits.ForEach(func(c *object.Commit) error {
		var rowBuffer bytes.Buffer
		rowTemplate := template.Must(template.New("row").Parse(`<tr>
<td class="commithash"><a href="{{.URLRoot}}/commit/{{.Commit.Hash.String}}">{{.ShortHash}}</a></td>
<td class="commitmessage">{{.Commit.Message}}</td>
<td class="author">
	<a href="mailto:{{.Commit.Author.Email}}">
	{{.Commit.Author.Name}}
	</a>
</td>
<td class="date">{{.Commit.Author.When.Format "2006-01-02 15:04:05"}}</td>
</tr>`))
		shortHash := c.Hash.String()[:7]
		rowTemplate.Execute(&rowBuffer, Row{&p.Config.URLRoot, &p.Branch, c, shortHash})
		rows = append(rows, rowBuffer.String())
		return nil
	})

	table := "<table>" + strings.Join(rows, "") + "</table>"

	descTemplate := template.Must(template.New("desc").Parse(`
		<p class="description">
			{{.}}
		</p>
		`))


	descTemplate.Execute(&bodyBuffer, "Showing commits for branch "+p.Branch)

	body = bodyBuffer.String() +
		table + "</article></main></body>"

	return body
}

func (p RepoBranchLogPage) FullPage() string {
	return "<!DOCTYPE html><html>" + p.Head() + p.Body() + "</html>"
}
