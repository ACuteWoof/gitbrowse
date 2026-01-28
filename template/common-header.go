package template

import (
	"bytes"
	"html/template"

	"git.lewoof.xyz/clone/gitbrowse/config"
)

func CommonHeader(c *config.PageConfig, currentPage string) string {
	type Page struct {
		Path string
		Name string
		Root *string
	}

	var pages []Page = []Page{
		{"/", "Readme", &c.URLRoot},
		{"/branch/master/tree", "Tree", &c.URLRoot},
		{"/branch/master/commit", "Commits", &c.URLRoot},
		{"/branch", "Branches", &c.URLRoot},
		{"/tag", "Tags", &c.URLRoot},
		{"/commit/HEAD", "Git Show", &c.URLRoot},
	}
	var headBuffer bytes.Buffer
	t := template.Must(template.New("head").Parse(`
	<body class="{{$currentPage}}">
	<header>
		<img src="{{.Config.Thumbnail}}" alt="Thumbnail">
		<div>
		<h1>{{.Config.Title}}</h1>
		<p>Clone URL: <code>{{.Config.CloneURL}}</code></p>
		<table>
			<tr>
			{{range .Pages}}
				{{if eq .Name $currentPage}}
					<td><em><a href="{{.Root}}{{.Path}}">{{.Name}}</a></em></td>
				{{else}}
					<td><a href="{{.Root}}{{Path}}">{{.Name}}</a></td>
				{{end}}
			</tr>
		</table>
		</div>
	</header>
	<main>
	<article>
	`))
	type TemplateInfo struct {
		Pages  []Page
		Config *config.PageConfig
	}
	var ti TemplateInfo = TemplateInfo{pages, c}
	t.Execute(&headBuffer, ti)
	return headBuffer.String()
}
