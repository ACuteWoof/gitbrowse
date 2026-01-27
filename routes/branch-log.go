package routes

import (
	"fmt"
	"net/http"

	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/template"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type RepoBranchLogRoute struct {
	ConfigGetter func(repo string) config.PageConfig
}

func (route RepoBranchLogRoute) Handler(w http.ResponseWriter, req *http.Request) {
	repo := req.PathValue("repo")
	branch := req.PathValue("branch")

	config := route.ConfigGetter(repo)

	r, err := git.PlainOpen(config.RootDir)
	errCheck(err)

	exists := false
	bs, err := r.Branches()
	errCheck(err)
	bs.ForEach(func(b *plumbing.Reference) error {
		if b.Name().Short() == branch {
			exists = true
			return nil
		}
		return nil
	})
	if !exists {
		http.Redirect(w, req, config.URLRoot + "/branch", http.StatusTemporaryRedirect)
		return
	}

	refName := plumbing.NewBranchReferenceName(branch)
	ref, err := r.Reference(refName, true)
	fmt.Fprintf(w, template.RepoBranchLogPage{Repo: r, Branch: branch, BranchRef: ref, Config: &config}.FullPage())
}
