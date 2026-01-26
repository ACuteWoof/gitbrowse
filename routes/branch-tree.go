package routes

import (
	"fmt"
	"net/http"
	"strings"

	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/template"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type RepoBranchTreeRoute struct {
	ConfigGetter func(repo string) config.PageConfig
}

func (route RepoBranchTreeRoute) Handler(w http.ResponseWriter, req *http.Request) {
	repo := req.PathValue("repo")
	branch := req.PathValue("branch")
	filePath := req.PathValue("filepath")
	filePath = strings.TrimSuffix(filePath, "/")

	config := route.ConfigGetter(repo)

	r, err := git.PlainOpen(config.RootDir)
	errCheck(err)

	b, err := r.Branch(branch)
	if err != nil || b == nil {
		http.Redirect(w, req, config.URLRoot + "/branch", http.StatusTemporaryRedirect)
		return
	}

	refName := plumbing.NewBranchReferenceName(branch)
	ref, err := r.Reference(refName, true)
	errCheck(err)
	commit, err := r.CommitObject(ref.Hash())
	errCheck(err)
	rootTree, err := commit.Tree()
	errCheck(err)
	if filePath == "" || filePath == "/" {
		fmt.Fprintf(w, template.RepoBranchTreePage{Tree: rootTree, Branch: branch, FilePath: filePath, Config: &config, BranchCommit: commit, Repo: r}.FullPage())
	} else {
		entry, err := rootTree.FindEntry(filePath)
		errCheck(err)
		if entry.Mode.IsFile() {
			return
		}
		tree, err := r.TreeObject(entry.Hash)
		fmt.Fprintf(w, template.RepoBranchTreePage{Tree: tree, Branch: branch, FilePath: filePath, Config: &config, BranchCommit: commit, Repo: r}.FullPage())
	}
}
