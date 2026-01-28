package routes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"git.lewoof.xyz/gitbrowse/config"
	"github.com/go-git/go-git/v6"
)

type TagDownloadRoute struct {
	ConfigGetter func(repo string) config.PageConfig
}

func (route TagDownloadRoute) Handler(w http.ResponseWriter, r *http.Request) {
	repoPath := route.ConfigGetter(r.PathValue("repo")).RootDir
	tagName := r.PathValue("name")
	fileName := r.PathValue("fileName")

	fileType := filepath.Ext(fileName)
	allowedFileTypes := []string{".zip", ".tar.gz", ".tgz", ".tar"}
	isAllowedFileType := slices.Contains(allowedFileTypes, fileType)
	if !isAllowedFileType {
		http.Error(w, "file type not allowed", http.StatusBadRequest)
		return
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		http.Error(w, "failed to open repository", http.StatusInternalServerError)
		return
	}
	tag, err := repo.Tag(tagName)
	if err != nil {
		http.Error(w, "failed to get tag", http.StatusInternalServerError)
		return
	}
	if tag == nil {
		http.Error(w, "tag not found", http.StatusNotFound)
		return
	}

	cmd := exec.Command("git", "archive", "--format="+strings.TrimPrefix(fileType, "."), tag.Name().Short())
	fmt.Println(tag.Name().Short())
	fmt.Println(cmd.String())
	cmd.Dir = repoPath

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("git archive failed: %v\n%s", err, buf.String()),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", filepath.Base(repoPath), tagName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	io.Copy(w, &buf)
}
