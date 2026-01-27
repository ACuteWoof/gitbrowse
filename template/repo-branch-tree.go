package template

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"

	"git.lewoof.xyz/clone/gitbrowse/config"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/filemode"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type RepoBranchTreePage struct {
	Repo         *git.Repository
	Branch       string
	BranchCommit *object.Commit
	Tree         *object.Tree
	FilePath     string
	Config       *config.PageConfig
}

func (p RepoBranchTreePage) Head() (head string) {
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

func (p RepoBranchTreePage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	t := template.Must(template.New("body").Parse(`
		<body class="repo-branch-tree">
			<header>
			 	<img src="{{.Config.Thumbnail}}" alt="Thumbnail">
				<div>
				<h1>{{.Config.Title}}</h1>
				<p>Clone URL: <code>{{.Config.CloneURL}}</code></p>
				<table>
					<tr>
					<td><a href="{{.Config.URLRoot}}/">Readme</a></td>
					<td><em><a href="{{.Config.URLRoot}}/branch/master/tree">Tree</a></em></td>
					<td><a href="{{.Config.URLRoot}}/branch/master/commit">Commits</a></td>
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

	rows := []string{}

	type Row struct {
		URLRoot    *string
		Branch     *string
		FilePath   *string
		Entry      *object.TreeEntry
		File       *object.File
		FileSize   string
		LastCommit *object.Commit
		Type       string
	}

	for _, entry := range p.Tree.Entries {
		var rowBuffer bytes.Buffer
		if entry.Mode.IsFile() {
			file, err := p.Tree.File(entry.Name)
			checkErr(err)

			rowTemplate := template.Must(template.New("row").Parse(`<tr>
<td class="isbinary" id="{{.Type}}">{{.Type}}</td>
<td class="filename"><a href="{{.URLRoot}}/branch/{{.Branch}}/tree/{{.FilePath}}{{.Entry.Name}}/">{{.Entry.Name}}</a></td>
<td class="commitmessage">
	<a href="{{.URLRoot}}/commit/{{.LastCommit.Hash.String}}">{{.LastCommit.Message}}</a>
</td>
<td class="author">
	<a href="mailto:{{.LastCommit.Author.Email}}">
	{{.LastCommit.Author.Name}}
	</a>
</td>
<td class="lastupdated">{{.LastCommit.Author.When.UTC.Format "2006-01-02 15:04:05"}} UTC</td>
<td class="filesize">{{.FileSize}}</td>
<td class="filemode">{{.File.Mode.ToOSFileMode}}</td>
</tr>`))
			fileRef, err := p.Repo.Reference(plumbing.NewBranchReferenceName(p.Branch), true)
			checkErr(err)
			log, err := p.Repo.Log(&git.LogOptions{From: fileRef.Hash(), Order: git.LogOrderCommitterTime, PathFilter: func(path string) bool {
				if p.FilePath == "" {
					return path == entry.Name
				}
				return path == p.FilePath+"/"+entry.Name
			}})
			checkErr(err)
			var lastCommit *object.Commit
			log.ForEach(func(c *object.Commit) error {
				lastCommit = c
				return nil
			})
			var fileType string
			isBinary, err := file.IsBinary()
			if isBinary {
				fileType = "bin"
			} else {
				fileType = "txt"
			}
			formattedFilePath := strings.TrimPrefix(p.FilePath+"/", "")
			if !strings.HasSuffix(formattedFilePath, "/") {
				formattedFilePath += "/"
			}
			if formattedFilePath == "/" {
				formattedFilePath = ""
			}

			rowTemplate.Execute(&rowBuffer, Row{&p.Config.URLRoot, &p.Branch, &formattedFilePath, &entry, file, getFormattedSize(float64(file.Size)), lastCommit, fileType})
		} else {
			var fileType string = "dir"
			formattedFilePath := strings.TrimPrefix(p.FilePath+"/", "")
			if !strings.HasSuffix(formattedFilePath, "/") {
				formattedFilePath += "/"
			}
			if formattedFilePath == "/" {
				formattedFilePath = ""
			}
			rowTemplate := template.Must(template.New("row").Parse(`<tr><td class="isbinary" id="{{.Type}}">{{.Type}}</td><td class="tree-dir filename"><a href="{{.URLRoot}}/branch/{{.Branch}}/tree/{{.FilePath}}{{.Entry.Name}}/">{{.Entry.Name}}</a></td><td></td><td></td><td></td><td></td><td class="filemode">{{.Entry.Mode.ToOSFileMode}}</td></tr>`))
			if entry.Mode == filemode.Submodule {
				fileType = "sub"
				rowTemplate = template.Must(template.New("row").Parse(`<tr><td class="isbinary" id="{{.Type}}">{{.Type}}</td><td class="filename">{{.Entry.Name}}</td><td></td><td></td><td></td><td></td><td class="filemode">{{.Entry.Mode.ToOSFileMode}}</td></tr>`))
			}
			rowTemplate.Execute(&rowBuffer, Row{&p.Config.URLRoot, &p.Branch, &formattedFilePath, &entry, nil, "", nil, fileType})
		}
		rows = append(rows, rowBuffer.String())
	}

	tableHeader := "<tr><th>Type</th><th>File</th><th>Commit Message</th><th>Author</th><th>Commit Date</th><th>Size</th><th>Mode</th></tr>"

	table := "<table>" + tableHeader + strings.Join(rows, "") + "</table>"

	type Crumb struct {
		Name   string
		DisplayName string
		Root   *string
		Branch *string
	}

	type Breadcrumb struct {
		Crumbs []Crumb
	}

	var breadcrumbsBuffer bytes.Buffer
	breadcrumbTemplate := template.Must(template.New("breadcrumb").Parse(`
	<table class="breadcrumbs">
		<tr>
		{{range .Crumbs}}
			<td><a href="{{.Root}}/branch/{{.Branch}}/tree/{{.Name}}">{{.DisplayName}}</a></td>
		{{end}}
		</tr>
	</table>
	`))
	defaultCrumbs := []Crumb{
		{"", "/", &p.Config.URLRoot, &p.Branch},
	}
	if p.FilePath != "" {
		for _, entry := range strings.Split(p.FilePath, "/") {
			defaultCrumbs = append(defaultCrumbs, Crumb{entry, entry, &p.Config.URLRoot, &p.Branch})
		}
	}
	breadcrumbTemplate.Execute(&breadcrumbsBuffer, Breadcrumb{defaultCrumbs})
	breadcrumbs := breadcrumbsBuffer.String()

	descTemplate := template.Must(template.New("desc").Parse(`
		<p class="description">
			{{.}}
		</p>
		`))
	descTemplate.Execute(&bodyBuffer, "Browsing tree for branch "+p.Branch+", showing "+strconv.Itoa(len(rows))+" entries")

	body = bodyBuffer.String() + breadcrumbs + table + "</article></main></body>"
	return
}

func (p RepoBranchTreePage) FullPage() string {
	return "<!DOCTYPE html><html>" + p.Head() + p.Body() + "</html>"
}
