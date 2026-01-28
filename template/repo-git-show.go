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
	bodyBuffer.WriteString(CommonHeader(p.Config, "Git Show"))
	gitShow := GitShow(p.Config.RootDir, p.Commit)
	body = bodyBuffer.String() + gitShow + "</article></main></body>"
	return
}

func (p RepoGitShowPage) FullPage() string {
	return "<!DOCTYPE html><html>" + CommonHead(p.Config) + p.Body() + "</html>"
}
