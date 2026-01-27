package template

import (
	"bytes"
	"html/template"
	"strings"
	"os/exec"

	"git.lewoof.xyz/gitbrowse/config"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type RepoBranchesPage struct {
	Repo      *git.Repository
	Config    *config.PageConfig
}

func (p RepoBranchesPage) Head() (head string) {
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

func (p RepoBranchesPage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	t := template.Must(template.New("body").Parse(`
		<body class="repo-branches">
			<header>
			 	<img src="{{.Config.Thumbnail}}" alt="Thumbnail">
				<div>
				<h1>{{.Config.Title}}</h1>
				<p>Clone URL: <code>{{.Config.CloneURL}}</code></p>
				<table>
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

	branches, err := p.Repo.Branches()
	checkErr(err)

	type Row struct {
		URLRoot *string
		Branch  string
		Commit  *object.Commit
		ShortHash string
	}

	rows := []string{}

	branches.ForEach(func(r *plumbing.Reference) error {
		c, err := p.Repo.CommitObject(r.Hash())
		checkErr(err)
		var rowBuffer bytes.Buffer
		rowTemplate := template.Must(template.New("row").Parse(`<tr>
<td class="branchname"><a href="{{.URLRoot}}/branch/{{.Branch}}/tree">{{.Branch}}</a></td>
<td class="commithash"><a href="{{.URLRoot}}/commit/{{.Commit.Hash.String}}">{{.ShortHash}}</a></td>
<td class="commitmessage">{{.Commit.Message}}</td>
<td class="author">
	<a href="mailto:{{.Commit.Author.Email}}">
	{{.Commit.Author.Name}}
	</a>
</td>
<td class="date">{{.Commit.Author.When.Format "2006-01-02 15:04:05"}}</td>
<td class="actions"><a href="{{.URLRoot}}/branch/{{.Branch}}/commit">See Commits</a></td>
</tr>`))
		cmd := exec.Command("git", "rev-parse", "--short", c.Hash.String())
		cmd.Dir = p.Config.RootDir
		shortHash, err := cmd.Output()
		rowTemplate.Execute(&rowBuffer, Row{&p.Config.URLRoot, r.Name().Short(), c, string(shortHash)})
		rows = append(rows, rowBuffer.String())
		return nil
	})

	tableHeader := "<tr><th>Branch</th><th>Head</th><th>Head Message</th><th>Head Author</th><th>Head Date</th><th>Actions</th></tr>"

	table := "<table>" + tableHeader + strings.Join(rows, "") + "</table>"

	descTemplate := template.Must(template.New("desc").Parse(`
		<p class="description">
			{{.}}
		</p>
		`))


	descTemplate.Execute(&bodyBuffer, "Showing branches for repository")

	body = bodyBuffer.String() +
		table + "</article></main></body>"

	return body
}

func (p RepoBranchesPage) FullPage() string {
	return "<!DOCTYPE html><html>" + p.Head() + p.Body() + "</html>"
}
