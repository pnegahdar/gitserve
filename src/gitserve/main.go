package main

import (
	"path"
	"os/exec"
	"io"
	"net/http"
	"mime"
	"strings"
	"flag"
	"os"
	"bytes"
	"time"
)

var prefix = flag.String("prefix", "", "The prefix path")
var listen = flag.String("listen", ":8020", "The listen arg to httpserver")
var root = flag.String("root", "", "The git root")

type GitArchive struct {
	gitRoot string // Directory .git resides
	prefix  string
}

func (ga *GitArchive) CommandBase() []string {
	dotGitDir := path.Join(ga.gitRoot, ".git")
	return []string{"git", "--git-dir", dotGitDir, "--work-tree", ga.gitRoot}
}

func (ga *GitArchive) WriteArchive(dir, treeish, format string, f io.Writer) error {
	baseCommand := ga.CommandBase()
	println(time.Now().UTC().String(), "Performing archive on git:", ga.gitRoot, "dir:", dir, "format:", format, "treeish", treeish)
	addOnArgs := []string{"archive", "--format", format, treeish + ":" + dir}
	fullCmd := append(baseCommand, addOnArgs...)
	println(time.Now().UTC().String(), "Running command", strings.Join(fullCmd, " "))
	cmd := exec.Command(fullCmd[0], fullCmd[1:]...)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	err := cmd.Run()
	if err != nil {
		return err
	}
	io.Copy(f, &buffer)
	return nil
}

func (ga *GitArchive) FetchLoop() {
	for {
		baseCmd := ga.CommandBase()
		addonArgs := []string{"fetch"}
		println(time.Now().UTC().String(), "Running git fetch.")
		fullCmd := append(baseCmd, addonArgs...)
		cmd := exec.Command(fullCmd[0], fullCmd[1:]...)
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 30)
	}
}

func (ga *GitArchive) HttpHandler(w http.ResponseWriter, r *http.Request) {
	dir := strings.TrimRight(path.Join(ga.prefix, r.URL.Path[1:]), "/") // Trim off front slash
	pathNoExt := strings.TrimRight(strings.Split(dir, ".")[0], "/")
	ext := strings.TrimLeft(path.Ext(dir), ".")
	if ext == "" {
		ext = "tar"
	}
	treeish := r.URL.Query().Get("tree")
	if treeish == ""{
		treeish = "HEAD"
	}
	mimeType := mime.TypeByExtension(ext)
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", "attachment; filename=archive." + ext)
	err := ga.WriteArchive(pathNoExt, treeish, ext, w)
	if err != nil {
		w.WriteHeader(500)
		panic(err)
	}

}

func main() {
	flag.Parse()
	dir := *root
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}
	ga := &GitArchive{gitRoot: dir, prefix: *prefix}
	http.HandleFunc("/", ga.HttpHandler)
	go ga.FetchLoop()
	err := http.ListenAndServe(*listen, nil)
	panic(err)

}
