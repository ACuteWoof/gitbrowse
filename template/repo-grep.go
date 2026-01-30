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
	"fmt"
	"html"
	"html/template"
	"regexp"

	"git.lewoof.xyz/gitbrowse/config"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type RepoGrepPage struct {
	Regex  string
	Branch string
	Config *config.PageConfig
}

func (p RepoGrepPage) Body() string {
	var bodyBuffer bytes.Buffer
	bodyBuffer.WriteString(CommonHeader(p.Config, "Grep"))

	p.SearchBar(&bodyBuffer)

	if p.Regex == "" {
		bodyBuffer.WriteString("<div class=\"error\">No query provided</div>")
		return bodyBuffer.String()
	}

	repo, err := git.PlainOpen(p.Config.RootDir)
	if err != nil {
		bodyBuffer.WriteString("<div class=\"error\">Invalid repository</div>")
		return bodyBuffer.String()
	}

	regex, err := regexp.Compile(p.Regex)
	if err != nil {
		bodyBuffer.WriteString("<div class=\"error\">Invalid regex in query</div>")
		return bodyBuffer.String()
	}

	grepOptions := git.GrepOptions{
		Patterns: []*regexp.Regexp{regex},
	}
	if p.Branch != "" {
		grepOptions.ReferenceName = plumbing.NewBranchReferenceName(p.Branch)
	}

	results, err := repo.Grep(&grepOptions)
	if err != nil {
		fmt.Println(err)
		bodyBuffer.WriteString("<div class=\"error\">Error searching</div>")	
		return bodyBuffer.String()
	}

	if len(results) == 0 {
		bodyBuffer.WriteString("<div class=\"error\">No results found for expression <code>" + html.EscapeString(p.Regex) + "</code></div>")
		return bodyBuffer.String()
	}

	bodyBuffer.WriteString("<table class=\"search-results\">")
	bodyBuffer.WriteString("<tr><th>File</th><th>Line No.</th><th>Line</th></tr>")

	t := template.Must(template.New("grep-result").Parse(`<tr>
	{{if .NeedFileName }}
		<td><a href="{{.Config.URLRoot}}/branch/{{.Branch}}/tree/{{.FileName}}#{{.LineNumber}}">{{.FileName}}</a></td>
		{{else}}
		<td></td>
		{{end}}
		<td><a href="{{.Config.URLRoot}}/branch/{{.Branch}}/tree/{{.FileName}}#{{.LineNumber}}">{{.LineNumber}}</a></td>
		<td class="search-content">
		`))

		type TemplateResult struct {
			Branch string
			FileName string
			LineNumber int
			Content string
			Config *config.PageConfig
			NeedFileName bool
		}

	lastFileName := ""
	for _, result := range results {
		defaultBranch := getHeadBranch(repo)
		r := TemplateResult{
			Branch: defaultBranch,
			FileName: result.FileName,
			LineNumber: result.LineNumber,
			Config: p.Config,
			NeedFileName: result.FileName != lastFileName,
		}
		if p.Branch != "" {
			r.Branch = p.Branch
		}
		t.Execute(&bodyBuffer, r)
		bodyBuffer.WriteString(getHighlightedHTML(result.FileName, result.Content))
		bodyBuffer.WriteString("</td></tr>")
		lastFileName = result.FileName
	}
	bodyBuffer.WriteString("</table>")

	return bodyBuffer.String()
}

func (p RepoGrepPage) FullPage() string {
	return "<!DOCTYPE html><html>" + CommonHead(p.Config) + p.Body() + "</html>"
}

func (p RepoGrepPage) SearchBar(bodyBuffer *bytes.Buffer) {
	bodyBuffer.WriteString("<div class=\"search\">")
	bodyBuffer.WriteString("<form class=\"search-form\" action=\"" + p.Config.URLRoot + "/grep\" method=\"get\">")
	bodyBuffer.WriteString("<input type=\"text\" name=\"q\" placeholder=\"Regex\" value=\"" + p.Regex + "\">")
	if p.Branch != "" {
		bodyBuffer.WriteString("<input type=\"hidden\" name=\"branch\" value=\"" + p.Branch + "\">")
	}
	bodyBuffer.WriteString("<button type=\"submit\">Search</button>")
	bodyBuffer.WriteString("</form>")
	bodyBuffer.WriteString("</div>")
}
