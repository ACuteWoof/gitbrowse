// Gitbrowse: a simple web server for git.
// Copyright (C) 2026 Vithushan
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package config

const SizePrecision = 2


// "domain/" uses IndexPageConfig

var IndexPageConfig PageConfig = PageConfig{
	URLRoot: "/",
	RootDir: "/home/acutewoof/gitbrowse",
	Title: "Users on git.lewoof.xyz",
	Description: "Gitbrowse on git.lewoof.xyz",
	Thumbnail: "/static/thumbnail.png",
	Favicon: "/static/favicon.ico",
	Styles: []string{
		"/static/styles.css",
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
		"/static/styles.css",
	},
}
