package config

var IndexPageConfig PageConfig = PageConfig{
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
	RootDir: "/home/acutewoof/Projects/acutewoof/<REPO>",
	Title: "<REPO>",
	Description: "<REPO> on lewoof.xyz",
	Thumbnail: "/static/thumbnail.png",
	Favicon: "/static/favicon.ico",
	Styles: []string{
		"/static/styles.css",
	},
}
