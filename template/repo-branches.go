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

type RepoBranchesPage struct {
	Repo      *git.Repository
	Config    *config.PageConfig
}

func (p RepoBranchesPage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	bodyBuffer.WriteString(CommonHeader(p.Config, "Branches"))

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
<td class="commithash"><a href="{{.URLRoot}}/show/{{.Commit.Hash.String}}">{{.ShortHash}}</a></td>
<td class="commitmessage">{{.Commit.Message}}</td>
<td class="author">
	<a href="mailto:{{.Commit.Author.Email}}">
	{{.Commit.Author.Name}}
	</a>
</td>
<td class="date">{{.Commit.Author.When.UTC.Format "15:04, Jan 2 2006"}}</td>
<td class="actions"><a href="{{.URLRoot}}/branch/{{.Branch}}/commit">See Commits</a></td>
</tr>`))
		cmd := exec.Command("git", "rev-parse", "--short", c.Hash.String())
		cmd.Dir = p.Config.RootDir
		shortHash, err := cmd.Output()
		rowTemplate.Execute(&rowBuffer, Row{&p.Config.URLRoot, r.Name().Short(), c, string(shortHash)})
		rows = append(rows, rowBuffer.String())
		return nil
	})

	tableHeader := "<tr><th>Branch</th><th>Head</th><th>Head Message</th><th>Head Author</th><th>Head Date (UTC)</th><th></th></tr>"

	table := "<table>" + tableHeader + strings.Join(rows, "") + "</table>"

	descTemplate := template.Must(template.New("desc").Parse(`
		<p class="description">
			{{.}}
		</p>
		`))


	descTemplate.Execute(&bodyBuffer, "Showing " + strconv.Itoa(len(rows)) + " branches for repository")

	body = bodyBuffer.String() +
		table + "</article></main></body>"

	return body
}

func (p RepoBranchesPage) FullPage() string {
	return "<!DOCTYPE html><html>" + CommonHead(p.Config) + p.Body() + "</html>"
}
