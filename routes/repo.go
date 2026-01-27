package routes

import (
	"fmt"
	"net/http"
	"strings"

	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/template"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/gomarkdown/markdown/html"
	stdhtml "html"
)

type RepoReadmeRoute struct {
	ConfigGetter func(repo string) config.PageConfig
}

func (route RepoReadmeRoute) Handler(w http.ResponseWriter, req *http.Request) {
	repo := req.PathValue("repo")
	config := route.ConfigGetter(repo)

	r, err := git.PlainOpen(config.RootDir)
	errCheck(err)

	ref, err := r.Head()
	errCheck(err)

	commit, err := r.CommitObject(ref.Hash())
	errCheck(err)

	tree, err := commit.Tree()
	errCheck(err)

	// file, err := tree.File("README.md")
	possibleReadmes := []string{"README.md", "README.markdown", "README.txt", "README", "readme", "README", "LICENSE", "LICENSE.md", "LICENSE.txt", "license", "license.txt", "license.md"}
	var readme *object.File
	for _, possibleReadme := range possibleReadmes {
		file, err := tree.File(possibleReadme)
		if err == nil {
			readme = file
			break
		}
	}
	if readme == nil {
		http.Redirect(w, req, req.URL.Path+"/branch", http.StatusTemporaryRedirect)
		return
	}

	content, err := readme.Contents()
	errCheck(err)

	var readmeHtml string
	if strings.HasSuffix(readme.Name, ".md") {
		readmeHtml = markdownToHtml(content)
	} else {
		readmeHtml = stdhtml.EscapeString(string(content))
	}

	fmt.Fprintf(w, template.RepoReadmePage{Readme: readmeHtml, Config: &config}.FullPage())
}
