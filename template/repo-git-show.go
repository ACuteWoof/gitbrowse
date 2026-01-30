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
	"git.lewoof.xyz/gitbrowse/config"
)

type RepoGitShowPage struct {
	Commit string
	Config *config.PageConfig
}

func (p RepoGitShowPage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	bodyBuffer.WriteString(CommonHeader(p.Config, "Show"))
	gitShow := GitShow(p.Config.RootDir, p.Commit)
	body = bodyBuffer.String() + gitShow + "</article></main></body>"
	return
}

func (p RepoGitShowPage) FullPage() string {
	return "<!DOCTYPE html><html>" + CommonHead(p.Config) + p.Body() + "</html>"
}
