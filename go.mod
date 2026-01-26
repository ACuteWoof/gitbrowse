module git.lewoof.xyz/gitbrowse

require git.lewoof.xyz/gitbrowse/routes v1.0.0

require git.lewoof.xyz/gitbrowse/config v1.0.0 // indirect
require git.lewoof.xyz/gitbrowse/template v1.0.0 // indirect

replace (
	git.lewoof.xyz/gitbrowse/config => ./config
	git.lewoof.xyz/gitbrowse/routes => ./routes
	git.lewoof.xyz/gitbrowse/template => ./template
)

go 1.25.6
