package main

import (
	"git.lewoof.xyz/clone/gitbrowse/config"
	"git.lewoof.xyz/clone/gitbrowse/routes"
	"net/http"
	"strings"
	// "os"
)

const defaultUser = "acutewoof"

func main() {
	startHttpServer()
}

func startHttpServer() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	staticMux := http.NewServeMux()
	staticMux.Handle("/", fs)

	setupMultiUserHandlers(mux) // change to setupMultiUserHandlers(mux) to enable multi user support

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/static")
			staticMux.ServeHTTP(w, r)
			return
		}

		mux.ServeHTTP(w, r)
	})

	http.ListenAndServe(":8088", nil)
}

func HandleStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

// func singleUserHandler(mux *http.ServeMux) {
// 	mux.HandleFunc("/", routes.IndexRoute{RepoRoute: "/", ConfigGetter: config.GetIndexConfg}.Handler)
// 	mux.HandleFunc("/{repo}/", routes.RepoReadmeRoute{ConfigGetter: config.GetRepoConfg}.Handler)
// 	mux.HandleFunc("/{repo}/branch/", routes.RepoBranchesRoute{ConfigGetter: config.GetRepoConfg}.Handler)
// 	mux.HandleFunc("/{repo}/branch/{branch}/tree/{filepath...}", routes.RepoBranchTreeRoute{ConfigGetter: config.GetRepoConfg}.Handler)
// 	mux.HandleFunc("/{repo}/branch/{branch}/commit", routes.RepoBranchLogRoute{ConfigGetter: config.GetRepoConfg}.Handler)
// 	mux.HandleFunc("/{repo}/tag/", routes.RepoTagsRoute{ConfigGetter: config.GetRepoConfg}.Handler)
// 	mux.HandleFunc("/{repo}/show/{hash}", routes.RepoGitShowRoute{ConfigGetter: config.GetRepoConfg}.Handler)
// }

func getIndexConfigGetterUser(username string) func() config.PageConfig {
	var IndexPageConfig config.PageConfig = config.PageConfig{
		URLRoot:     "/" + username,
		RootDir:     "/home/" + username + "/gitbrowse", // directory in the unix filesystem, ls here is the list of repos displayed
		Title:       username,
		Description: username + " on git.lewoof.xyz",
		Thumbnail:   "/static/thumbnail.png",
		Favicon:     "/static/favicon.ico",
		Styles: []string{
			"/static/global.css",
		},
	}
	return func() config.PageConfig {
		return IndexPageConfig
	}
}

func getRepoConfigGetter(username string) func(repo string) config.PageConfig {
	return func(repo string) config.PageConfig {
		var RepoPageConfig config.PageConfig = config.PageConfig{
			URLRoot:     "/" + username + "/" + repo,         // url path, don't bother changing
			RootDir:     "/home/" + username + "/gitbrowse/" + repo, // directory in the unix filesystem where each repo is stored
			CloneURL:    "https://git.lewoof.xyz/clone/" + username + "/" + repo, // url used to clone the repo
			Title:       username + "/" + repo,
			Description: username + "/" + repo + " on git.lewoof.xyz",
			Thumbnail:   "/static/thumbnail.png",
			Favicon:     "/static/favicon.ico",
			Styles: []string{
				"/static/global.css",
			},
		}
		return RepoPageConfig
	}
}

func setupMultiUserHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/"+defaultUser, http.StatusMovedPermanently)
	})

	mux.HandleFunc("/{user}", func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user")
		routes.IndexRoute{RepoRoute: "/" + user + "/", ConfigGetter: getIndexConfigGetterUser(user)}.Handler(w, r)
	})

	mux.HandleFunc("/{user}/{repo}/", func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user")
		routes.RepoReadmeRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
	})
	mux.HandleFunc("/{user}/{repo}/branch/", func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user")
		routes.RepoBranchesRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
	})
	mux.HandleFunc("/{user}/{repo}/branch/{branch}/tree/{filepath...}", func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user")
		routes.RepoBranchTreeRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
	})
	mux.HandleFunc("/{user}/{repo}/branch/{branch}/commit", func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user")
		routes.RepoBranchLogRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
	})
	mux.HandleFunc("/{user}/{repo}/tag/", func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user")
		routes.RepoTagsRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
	})
	mux.HandleFunc("/{user}/{repo}/tag/{name}/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user")
		routes.TagDownloadRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
	})
	mux.HandleFunc("/{user}/{repo}/show/{hash}", func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user")
		routes.RepoGitShowRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
	})
}
