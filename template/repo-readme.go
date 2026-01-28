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
	"strings"

	"git.lewoof.xyz/gitbrowse/config"
)

type RepoInfoPage struct {
	InfoFiles []RepoInfoFile
	Config *config.PageConfig
}

type RepoInfoFile struct {
	Name string
	RenderedContent string
}

func (p RepoInfoPage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	bodyBuffer.WriteString(CommonHeader(p.Config, "Info"))
	bodyBuffer.WriteString("<div class=\"info\">")
	bodyBuffer.WriteString("<table class=\"info-nav\"><tbody><tr>")
	for _, infoFile := range p.InfoFiles {
		bodyBuffer.WriteString("<td class=\"txt\"><a href=\"#" + infoFile.Name + "\">" + infoFile.Name + "</a></td>")
	}
	bodyBuffer.WriteString("</tr></tbody></table>")
	for _, infoFile := range p.InfoFiles {
		tdClass := "txt"
		if strings.HasSuffix(infoFile.Name, ".md") {
			tdClass = "content"
		}
		bodyBuffer.WriteString("<table class=\"info\" id=\"" + infoFile.Name + "\"><tbody><tr><td>" + infoFile.Name + "</td></tr><tr><td class=\"" + tdClass + "\">")
		bodyBuffer.WriteString(infoFile.RenderedContent)
		bodyBuffer.WriteString("</td></tr></tbody></table>")
	}
	body = bodyBuffer.String() + "</div></article></main></body>"
	return
}

func (p RepoInfoPage) FullPage() string {
	return "<!DOCTYPE html><html>" + CommonHead(p.Config) + p.Body() + "</html>"
}
