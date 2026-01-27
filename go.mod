module git.lewoof.xyz/clone/gitbrowse

require git.lewoof.xyz/clone/gitbrowse/routes v1.0.0

require git.lewoof.xyz/clone/gitbrowse/config v1.0.0

require (
	git.lewoof.xyz/clone/gitbrowse/template v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.3.0 // indirect
	github.com/alecthomas/chroma/v2 v2.23.1 // indirect
	github.com/cloudflare/circl v1.6.1 // indirect
	github.com/cyphar/filepath-securejoin v0.6.1 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg/v2 v2.0.2 // indirect
	github.com/go-git/go-billy/v6 v6.0.0-20260114122816-19306b749ecc // indirect
	github.com/go-git/go-git/v6 v6.0.0-20260123133532-f99a98e81ce9 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/gomarkdown/markdown v0.0.0-20250810172220-2e2c11897d1a // indirect
	github.com/kevinburke/ssh_config v1.4.0 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/pjbgf/sha1cd v0.5.0 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
)

replace (
	git.lewoof.xyz/clone/gitbrowse/config => ./config
	git.lewoof.xyz/clone/gitbrowse/routes => ./routes
	git.lewoof.xyz/clone/gitbrowse/template => ./template
)

go 1.25.6
