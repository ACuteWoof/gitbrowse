package routes

import (
	"fmt"
	"git.lewoof.xyz/gitbrowse/config"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"net/http"
	"git.lewoof.xyz/gitbrowse/template"
)

type RepoReadmeRoute struct {
	ConfigGetter func(repo string) config.PageConfig
}

func (route RepoReadmeRoute) Handler(w http.ResponseWriter, req *http.Request) {
	repo := req.PathValue("repo")
	config := route.ConfigGetter(repo)

	fmt.Println(config.RootDir)

	r, err := git.PlainOpen(config.RootDir)
	errCheck(err)

	ref, err := r.Head()
	errCheck(err)

	commit, err := r.CommitObject(ref.Hash())
	errCheck(err)

	tree, err := commit.Tree()
	errCheck(err)

	// file, err := tree.File("README.md")
	possibleReadmes := []string{"README.md", "README.markdown", "README.txt", "README", "readme", "README"}
	var readme *object.File
	for _, possibleReadme := range possibleReadmes {
		file, err := tree.File(possibleReadme)
		if err == nil {
			readme = file
			break
		}
	}
	if readme == nil {
		http.Redirect(w, req, req.URL.Path+"/tree/master", http.StatusTemporaryRedirect)
		return
	}

	content, err := readme.Contents()
	errCheck(err)
	readmeHtml := markdownToHtml(content)

	fmt.Fprintf(w, template.RepoReadmePage{Readme: readmeHtml, Config: &config}.FullPage())
}
