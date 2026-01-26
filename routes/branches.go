package routes

import (
	"fmt"
	"net/http"

	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/template"

	"github.com/go-git/go-git/v6"
)

type RepoBranchesRoute struct {
	ConfigGetter func(repo string) config.PageConfig
}

func (route RepoBranchesRoute) Handler(w http.ResponseWriter, req *http.Request) {
	repo := req.PathValue("repo")

	config := route.ConfigGetter(repo)

	r, err := git.PlainOpen(config.RootDir)
	errCheck(err)
	fmt.Fprintf(w, template.RepoBranchesPage{Repo: r, Config: &config}.FullPage())
}
