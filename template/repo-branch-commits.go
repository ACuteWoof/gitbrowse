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
	"os/exec"
	"strconv"
	"strings"

	"git.lewoof.xyz/gitbrowse/config"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type RepoBranchLogPage struct {
	Format    string
	Repo      *git.Repository
	Branch    string
	BranchRef *plumbing.Reference
	Config    *config.PageConfig
}

func (p RepoBranchLogPage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	bodyBuffer.WriteString(CommonHeader(p.Config, "Commits"))

	commits, err := p.Repo.Log(&git.LogOptions{From: p.BranchRef.Hash(), Order: git.LogOrderCommitterTime})
	checkErr(err)

	type Row struct {
		URLRoot   *string
		Branch    *string
		Commit    *object.Commit
		ShortHash string
	}

	rows := []string{}

	commits.ForEach(func(c *object.Commit) error {
		var rowBuffer bytes.Buffer
		var rowTemplate *template.Template
		if p.Format == "rss" {
			rowTemplate = template.Must(template.New("row").Parse(`<item>
<link>{{.URLRoot}}/show/{{.Commit.Hash.String}}</link>
<title>{{.Commit.Message}}</title>
<description><![CDATA[
<p>{{.ShortHash}}</p>
<p>{{.Commit.Message}}</p>
<hr>
<p>
{{.Commit.Author.Name}}
({{.Commit.Author.Email}})
<p>
]]></description>
<author>
	{{.Commit.Author.Name}}
	({{.Commit.Author.Email}})
</author>
<pubDate>{{.Commit.Author.When.UTC.Format "Mon, 02 Jan 2006 15:04:05 GMT"}}</pubDate>
</item>`))
		} else {
			rowTemplate = template.Must(template.New("row").Parse(`<tr>
<td class="commithash"><a href="{{.URLRoot}}/show/{{.Commit.Hash.String}}">{{.ShortHash}}</a></td>
<td class="commitmessage">{{.Commit.Message}}</td>
<td class="author">
	<a href="mailto:{{.Commit.Author.Email}}">
	{{.Commit.Author.Name}}
	</a>
</td>
<td class="date">{{.Commit.Author.When.UTC.Format "15:04, Jan 2 2006"}}</td>
</tr>`))
		}

		cmd := exec.Command("git", "-c", "safe.directory="+p.Config.RootDir, "rev-parse", "--short", c.Hash.String())
		cmd.Dir = p.Config.RootDir
		shortHash, err := cmd.Output()
		if err != nil {
			c.Hash.Write(shortHash)
		}
		checkErr(err)

		rowTemplate.Execute(&rowBuffer, Row{&p.Config.URLRoot, &p.Branch, c, string(shortHash)})
		rows = append(rows, strings.Replace(rowBuffer.String(), ">&lt;![CDATA[", "><![CDATA[", 1))
		return nil
	})

	var tableHeader, table string
	if p.Format == "rss" {
		tableHeader = "<title>" + p.Config.Title + "</title><link>" + p.Config.URLRoot + "</link><description>Commits on branch " + p.Branch + "</description><language>en-us</language>"
	} else {
		tableHeader = "<tr><th>Commit</th><th>Message</th><th>Author</th><th>Date (UTC)</th></tr>"
	}

	if p.Format == "rss" {
		table = "<channel>" + tableHeader + strings.Join(rows, "") + "</channel>"
	        return table
	} else {
		table = "<table>" + tableHeader + strings.Join(rows, "") + "</table>"
	}

	descTemplate := template.Must(template.New("desc").Parse(`
		<p class="description">
			{{.}}
		</p>
		`))

	descTemplate.Execute(&bodyBuffer, "Showing "+strconv.Itoa(len(rows))+" commits for branch "+p.Branch)

	body = bodyBuffer.String() +
		table + "</article></main></body>"

	return body
}

func (p RepoBranchLogPage) FullPage() string {
	if p.Format == "rss" {
		return `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" 
     xmlns:atom="http://www.w3.org/2005/Atom"
     xmlns:content="http://purl.org/rss/1.0/modules/content/">
` + p.Body() + "</rss>"
	}
	return "<!DOCTYPE html><html>" + CommonHead(p.Config) + p.Body() + "</html>"
}
