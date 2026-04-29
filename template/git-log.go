package template

import (
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func GetLog1File(repo *git.Repository, repoPath string, path string) (*object.Commit) {
	cmd := exec.Command("git", "-c", "safe.directory="+repoPath, "log", `--format=%H`, "-1", "--", "."+path)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	outstr := strings.TrimSpace(string(out))
	hash, exists := plumbing.FromHex(outstr)
	if !exists {
		return nil
	}

	commit, err := repo.CommitObject(hash)
	checkErr(err)

	return commit
}
