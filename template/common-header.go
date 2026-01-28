package template

import (
	"bytes"
	"html/template"
	"log"

	"git.lewoof.xyz/gitbrowse/config"
	"github.com/go-git/go-git/v6"
)

func CommonHeader(c *config.PageConfig, currentPage string) string {
	r, err := git.PlainOpen(c.RootDir)
	checkErr(err)
	headBranch := getHeadBranch(r)

	type Page struct {
		Path string
		Name string
		Root *string
	}

	var pages []Page = []Page{
		{"/", "Readme", &c.URLRoot},
		{"/branch/" + headBranch + "/tree", "Tree", &c.URLRoot},
		{"/branch/" + headBranch + "/commit", "Commits", &c.URLRoot},
		{"/branch", "Branches", &c.URLRoot},
		{"/tag", "Tags", &c.URLRoot},
		{"/show/HEAD", "Git Show", &c.URLRoot},
	}
	var headBuffer bytes.Buffer
	t := template.Must(template.New("head").Parse(`
	<body class="{{.CurrentPage}}">
	<header>
		<img src="{{.Config.Thumbnail}}" alt="Thumbnail">
		<div>
		<h1>{{.Config.Title}}</h1>
		{{if .Config.CloneURL}}
			<p>Clone URL: <code>{{.Config.CloneURL}}</code></p>
		{{end}}
		<table>
			<tr>
			{{range .Pages}}
				{{if eq .Name $.CurrentPage}}
					<td><em><a href="{{.Root}}{{.Path}}">{{.Name}}</a></em></td>
				{{else}}
					<td><a href="{{.Root}}{{.Path}}">{{.Name}}</a></td>
				{{end}}
			{{end}}
			</tr>
		</table>
		</div>
	</header>
	<main>
	<article>
	`))
	type TemplateInfo struct {
		Pages       []Page
		Config      *config.PageConfig
		CurrentPage string
	}
	var ti TemplateInfo = TemplateInfo{pages, c, currentPage}
	t.Execute(&headBuffer, ti)
	return headBuffer.String()
}

func getHeadBranch(repo *git.Repository) string {
	head, err := repo.Head()
	if err != nil {
		log.Fatal(err)
		return "master"
	}

	if head.Name().IsBranch() {
		branchName := head.Name().Short()
		return branchName
	}

	return "master"
}
