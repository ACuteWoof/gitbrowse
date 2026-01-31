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

type RepoTagsPage struct {
	Repo   *git.Repository
	Config *config.PageConfig
}

func (p RepoTagsPage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	bodyBuffer.WriteString(CommonHeader(p.Config, "Tags"))
	tags, err := p.Repo.Tags()
	checkErr(err)

	type Row struct {
		URLRoot   *string
		Branch    string
		Tag       *object.Tag
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
<td class="tag"><a href="{{.URLRoot}}/show/{{.Tag.Hash.String}}">{{.Tag.Name}}</a></td>
<td class="taghash"><a href="{{.URLRoot}}/show/{{.Tag.Hash.String}}">{{.ShortHash}}</a></td>
<td class="tagmessage">{{.Tag.Message}}</td>
<td class="tagger">
	<a href="mailto:{{.Tag.Tagger.Email}}">
	{{.Tag.Tagger.Name}}
	</a>
</td>
<td class="date">{{.Tag.Tagger.When.UTC.Format "15:04 Jan 2 2006"}}</td>
<td class="download">
<table class="subtable">
<tr>
<td>	<a href="{{.URLRoot}}/tag/{{.Tag.Name}}/{{.Tag.Name}}.zip">zip</a></td>
<td>	<a href="{{.URLRoot}}/tag/{{.Tag.Name}}/{{.Tag.Name}}.tar.gz">tar.gz</a></td>
<td>	<a href="{{.URLRoot}}/tag/{{.Tag.Name}}/{{.Tag.Name}}.tar.gz">tar</a></td>
<td>	<a href="{{.URLRoot}}/tag/{{.Tag.Name}}/{{.Tag.Name}}.tgz">tgz</a></td>
	</tr>
	</table>
</td>
</tr>`))
		cmd := exec.Command("git", "-c", "safe.directory="+p.Config.RootDir, "rev-parse", "--short", t.Hash.String())
		cmd.Dir = p.Config.RootDir
		shortHash, err := cmd.Output()
		if err != nil {
			t.Hash.Write(shortHash)
		}
		rowTemplate.Execute(&rowBuffer, Row{&p.Config.URLRoot, r.Name().Short(), t, string(shortHash)})
		rows = append(rows, rowBuffer.String())
		return nil
	})

	tableHeader := "<tr><th>Name</th><th>Hash</th><th>Message</th><th>Tagger</th><th>Date</th><th>Download</th></tr>"

	table := "<table>" + tableHeader + strings.Join(rows, "") + "</table>"

	descTemplate := template.Must(template.New("desc").Parse(`
		<p class="description">
			{{.}}
		</p>
		`))

	descTemplate.Execute(&bodyBuffer, "Showing "+strconv.Itoa(len(rows))+" tags for repository")

	body = bodyBuffer.String() +
		table + "</article></main></body>"

	return body
}

func (p RepoTagsPage) FullPage() string {
	return "<!DOCTYPE html><html>" + CommonHead(p.Config) + p.Body() + "</html>"
}
