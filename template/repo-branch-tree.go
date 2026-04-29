// Gitbrowse: a simple web server for git.
// Copyright (C) 2026 Vithushan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package template

import (
	"bytes"
	"html/template"
	"slices"
	"strconv"
	"strings"

	"git.lewoof.xyz/gitbrowse/config"
	"github.com/go-git/go-git/v6"

	// "github.com/go-git/go-git/v6/plumbing"
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

func (p RepoBranchTreePage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	bodyBuffer.WriteString(CommonHeader(p.Config, "Tree"))

	rows := []string{}

	infofiles := []RepoInfoFile{}

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

			if slices.Contains(config.GetPossibleInfoFiles(), file.Name) {
				contents, err := file.Contents()
				checkErr(err)
				var renderedContent string
				if file.Name[len(file.Name)-3:] == ".md" {
					renderedContent = MarkdownToHtml(contents)
				} else {
					renderedContent = template.HTMLEscapeString(contents)
				}
				infofiles = append(infofiles, RepoInfoFile{Name: file.Name, RenderedContent: renderedContent})
			}

			rowTemplate := template.Must(template.New("row").Parse(`<tr>
<td class="isbinary" id="{{.Type}}">{{.Type}}</td>
<td class="filename"><a href="{{.URLRoot}}/branch/{{.Branch}}/tree/{{.FilePath}}{{.Entry.Name}}/">{{.Entry.Name}}</a></td>
<td class="commitmessage" data-fulltext="{{.LastCommit.Message}}">
	<a href="{{.URLRoot}}/show/{{.LastCommit.Hash.String}}">{{.LastCommit.Message}}</a>
</td>
<td class="author">
	<a href="mailto:{{.LastCommit.Author.Email}}">
	{{.LastCommit.Author.Name}}
	</a>
</td>
<td class="lastupdated">{{.LastCommit.Author.When.UTC.Format "15:04, Jan 2 2006"}}</td>
<td class="filesize">{{.FileSize}}</td>
</tr>`))
			lastCommit := GetLog1File(p.Repo, p.Config.RootDir, p.FilePath+"/"+entry.Name)

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
			rowTemplate := template.Must(template.New("row").Parse(`<tr><td class="isbinary" id="{{.Type}}">{{.Type}}</td><td class="tree-dir filename"><a href="{{.URLRoot}}/branch/{{.Branch}}/tree/{{.FilePath}}{{.Entry.Name}}/">{{.Entry.Name}}</a></td><td></td><td></td><td></td><td></td></tr>`))
			if entry.Mode == filemode.Submodule {
				fileType = "sub"
				rowTemplate = template.Must(template.New("row").Parse(`<tr><td class="isbinary" id="{{.Type}}">{{.Type}}</td><td class="filename">{{.Entry.Name}}</td><td></td><td></td><td></td><td></td></tr>`))
			}
			rowTemplate.Execute(&rowBuffer, Row{&p.Config.URLRoot, &p.Branch, &formattedFilePath, &entry, nil, "", nil, fileType})
		}
		rows = append(rows, rowBuffer.String())
	}

	tableHeader := "<tr><th>Type</th><th>File</th><th>Last Commit</th><th>Author</th><th>Commit Date (UTC)</th><th>Size</th></tr>"

	table := `<table class="tree">` + tableHeader + strings.Join(rows, "") + "</table>"

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
		cumulativeName := ""
		for entry := range strings.SplitSeq(p.FilePath, "/") {
			if cumulativeName == "" {
				cumulativeName = entry
			} else {
				cumulativeName += "/" + entry
			}
			defaultCrumbs = append(defaultCrumbs, Crumb{cumulativeName, entry, &p.Config.URLRoot, &p.Branch})
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

	var renderingBuffer bytes.Buffer

	renderingBuffer.WriteString(`<div class="tree-info">`)
	renderingBuffer.WriteString("<table class=\"nav\"><tbody><tr>")
	for _, infoFile := range infofiles {
		renderingBuffer.WriteString("<td class=\"txt\"><a href=\"#" + infoFile.Name + "\">" + infoFile.Name + "</a></td>")
	}
	renderingBuffer.WriteString("</tr></tbody></table>")
	for _, infoFile := range infofiles {
		tdClass := "txt"
		if strings.HasSuffix(infoFile.Name, ".md") {
			tdClass = "content"
		}
		renderingBuffer.WriteString("<table class=\"info\" id=\"" + infoFile.Name + "\"><tbody><tr><td>" + infoFile.Name + "</td></tr><tr><td class=\"" + tdClass + "\">")
		renderingBuffer.WriteString(infoFile.RenderedContent)
		renderingBuffer.WriteString("</td></tr></tbody></table>")
	}
	renderingBuffer.WriteString(`</div>`)

	body = bodyBuffer.String() + breadcrumbs + table + renderingBuffer.String() + "</article></main></body>"
	return
}

func (p RepoBranchTreePage) FullPage() string {
	return "<!DOCTYPE html><html>" + CommonHead(p.Config) + p.Body() + "</html>"
}
