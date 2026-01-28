package config

const SizePrecision = 2


// "domain/" uses IndexPageConfig

var IndexPageConfig PageConfig = PageConfig{
	URLRoot: "/",
	RootDir: "/home/",
	Title: "Users on git.lewoof.xyz",
	Description: "Gitbrowse on git.lewoof.xyz",
	Thumbnail: "/static/thumbnail.png",
	Favicon: "/static/favicon.ico",
	Styles: []string{
		"/static/global.css",
	},
}

// IRRELEVANT FOR MULTI USER SETUP

// <REPO> will be replaced by the name of the repository
// this is the config used for all the tabs on the repo page
var RepoPageConfig PageConfig = PageConfig{
	URLRoot: "/<REPO>",
	CloneURL: "https://git.lewoof.xyz/<REPO>",
	RootDir: "/home/acutewoof/gitbrowse/<REPO>",
	Title: "lewoof/<REPO>",
	Description: "<REPO> on git.lewoof.xyz",
	Thumbnail: "/static/thumbnail.png",
	Favicon: "/static/favicon.ico",
	Styles: []string{
		"/static/global.css",
	},
}
