package config

const SizePrecision = 2

var IndexPageConfig PageConfig = PageConfig{
	URLRoot: "/",
	RootDir: "/home/acutewoof/Projects/acutewoof/",
	Title: "Gitbrowse on lewoof.xyz",
	Description: "lewoof on lewoof.xyz",
	Thumbnail: "/static/thumbnail.png",
	Favicon: "/static/favicon.ico",
	Styles: []string{
		"/static/styles.css",
	},
}

// <REPO> will be replaced by the name of the repository
// this is the config used for all the tabs on the repo page
var RepoPageConfig PageConfig = PageConfig{
	URLRoot: "/browse/<REPO>",
	RootDir: "/home/acutewoof/Projects/acutewoof/<REPO>",
	Title: "lewoof/<REPO>",
	Description: "<REPO> on lewoof.xyz",
	Thumbnail: "/static/thumbnail.png",
	Favicon: "/static/favicon.ico",
	Styles: []string{
		"/static/styles.css",
	},
}
