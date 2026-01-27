package main

import (
	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/routes"
	"net/http"
	// "os"
)

func main() {
	startHttpServer()
}

func startHttpServer() {
	setupSingleUserHandlers() // change to setupMultiUserHandlers() to enable multi user support

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8088", nil)
}

func HandleStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func setupSingleUserHandlers() {
	http.HandleFunc("/", routes.IndexRoute{RepoRoute: "/browse/", ConfigGetter: config.GetIndexConfg}.Handler)
	http.HandleFunc("/browse/{repo}/", routes.RepoReadmeRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	http.HandleFunc("/browse/{repo}/branch/", routes.RepoBranchesRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	http.HandleFunc("/browse/{repo}/branch/{branch}/tree/{filepath...}", routes.RepoBranchTreeRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	http.HandleFunc("/browse/{repo}/branch/{branch}/commit", routes.RepoBranchLogRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	http.HandleFunc("/browse/{repo}/tag/", routes.RepoTagsRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	http.HandleFunc("/browse/{repo}/commit/{hash}", routes.RepoGitShowRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	http.HandleFunc("/browse/{repo}/tag/{hash}", routes.RepoGitShowRoute{ConfigGetter: config.GetRepoConfg}.Handler)
}

// uncomment the code below, and comment out the function setupSingleUserHandlers() to enable multi user support

//
// func getIndexConfigGetterUser(username string) func() config.PageConfig {
// 	var IndexPageConfig config.PageConfig = config.PageConfig{
// 		URLRoot:     "/" + username,
// 		RootDir:     "/home/" + username + "/gitbrowse", // directory in the unix filesystem, ls here is the list of repos displayed
// 		Title:       username,
// 		Description: username + " on git.lewoof.xyz",
// 		Thumbnail:   "/static/thumbnail.png",
// 		Favicon:     "/static/favicon.ico",
// 		Styles: []string{
// 			"/static/global.css",
// 		},
// 	}
// 	return func() config.PageConfig {
// 		return IndexPageConfig
// 	}
// }
//
// func getRepoConfigGetter(username string) func(repo string) config.PageConfig {
// 	return func(repo string) config.PageConfig {
// 		var RepoPageConfig config.PageConfig = config.PageConfig{
// 			URLRoot:     "/browse/" + username + "/" + repo, // url path, don't bother changing
// 			RootDir:     "/home/" + username + "/gitbrowse/" + repo, // directory in the unix filesystem where each repo is stored
// 			Title:       username + "/" + repo,
// 			Description: username + "/" + repo + " on git.lewoof.xyz",
// 			Thumbnail:   "/static/thumbnail.png",
// 			Favicon:     "/static/favicon.ico",
// 			Styles: []string{
// 				"/static/global.css",
// 			},
// 		}
// 		return RepoPageConfig
// 	}
// }
//
// func setupMultiUserHandlers() {
// 	// by default "/" redirects to user woof, change it to your default user
//      const defaultUser = "woof"
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		http.Redirect(w, r, "/"+defaultUser, http.StatusTemporaryRedirect)
// 	})
//
// 	http.HandleFunc("/{user}", func(w http.ResponseWriter, r *http.Request) {
// 		user := r.PathValue("user")
// 		routes.IndexRoute{RepoRoute: "/browse/" + user + "/", ConfigGetter: getIndexConfigGetterUser(user)}.Handler(w, r)
// 	})
//
// 	// you probably don't want to change any of this below
//
// 	http.HandleFunc("/browse/{user}/{repo}/", func(w http.ResponseWriter, r *http.Request) {
// 		user := r.PathValue("user")
// 		routes.RepoReadmeRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
// 	})
// 	http.HandleFunc("/browse/{user}/{repo}/branch/", func(w http.ResponseWriter, r *http.Request) {
// 		user := r.PathValue("user")
// 		routes.RepoBranchesRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
// 	})
// 	http.HandleFunc("/browse/{user}/{repo}/branch/{branch}/tree/{filepath...}", func(w http.ResponseWriter, r *http.Request) {
// 		user := r.PathValue("user")
// 		routes.RepoBranchTreeRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
// 	})
// 	http.HandleFunc("/browse/{user}/{repo}/branch/{branch}/commit", func(w http.ResponseWriter, r *http.Request) {
// 		user := r.PathValue("user")
// 		routes.RepoBranchLogRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
// 	})
// 	http.HandleFunc("/browse/{user}/{repo}/tag/", func(w http.ResponseWriter, r *http.Request) {
// 		user := r.PathValue("user")
// 		routes.RepoTagsRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
// 	})
// 	http.HandleFunc("/browse/{user}/{repo}/commit/{hash}", func(w http.ResponseWriter, r *http.Request) {
// 		user := r.PathValue("user")
// 		routes.RepoGitShowRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
// 	})
// 	http.HandleFunc("/browse/{user}/{repo}/tag/{hash}", func(w http.ResponseWriter, r *http.Request) {
// 		user := r.PathValue("user")
// 		routes.RepoGitShowRoute{ConfigGetter: getRepoConfigGetter(user)}.Handler(w, r)
// 	})
// }
