package template

import (
	"bytes"
	"git.lewoof.xyz/gitbrowse/config"
)

type RepoReadmePage struct {
	Readme string
	Config *config.PageConfig
}

func (p RepoReadmePage) Body() (body string) {
	var bodyBuffer bytes.Buffer
	bodyBuffer.WriteString(CommonHeader(p.Config, "Readme"))
	body = bodyBuffer.String() + p.Readme + "</article></main></body>"
	return
}

func (p RepoReadmePage) FullPage() string {
	return "<!DOCTYPE html><html>" + CommonHead(p.Config) + p.Body() + "</html>"
}
