package routes

import (
	"fmt"
	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/template"
	"net/http"
	"os"
)

type IndexRoute struct {
	RepoRoute string;
	ConfigGetter func() config.PageConfig;
}

func (i IndexRoute) Handler(w http.ResponseWriter, _ *http.Request) {
	var dirs []string = getGitDirs(i.ConfigGetter().RootDir, i.RepoRoute)
	config := i.ConfigGetter()
	fmt.Fprintf(w, template.IndexPage{Repos: dirs, Config: &config}.FullPage())
}

func getGitDirs(root string, displayRoot string) []string {
	d, err := os.ReadDir(root)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var dirs []string
	for _, dir := range d {
		if dir.IsDir() {
			dirs = append(dirs, displayRoot+dir.Name())
		}
	}
	return dirs
}
