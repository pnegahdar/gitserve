package main

import (
	"bytes"
	"flag"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

var prefix = flag.String("prefix", "", "The prefix path")
var listen = flag.String("listen", ":8020", "The listen arg to httpserver")
var root = flag.String("root", "", "The git root")
var pypiPrefix = flag.String("pypi-tag-prefix", "", "pypi tag prefix")
var pypiPrefixDelimiter = flag.String("pypi-tag-delimiter", "", "what to split to tag on to find the package version")

type GitArchive struct {
	gitRoot string // Directory .git resides
	prefix  string
}

func (ga *GitArchive) CommandBase() []string {
	dotGitDir := path.Join(ga.gitRoot, ".git")
	return []string{"git", "--git-dir", dotGitDir, "--work-tree", ga.gitRoot}
}

func (ga *GitArchive) CommandOutput(args ...string) (string, error) {
	baseCommand := ga.CommandBase()
	fullCmd := append(baseCommand, args...)
	cmd := exec.Command(fullCmd[0], fullCmd[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
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
		time.Sleep(time.Second * 10)
	}
}

func (ga *GitArchive) HttpHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rest := vars["rest"]
	dir := strings.TrimRight(path.Join(ga.prefix, rest), "/") // Trim off front slash
	strSplit := strings.Split(dir, ".")
	pathNoExt := strings.TrimRight(strSplit[0], "/")
	ext := "tar"
	if len(strSplit) > 1 {
		ext = strings.Join(strSplit[1:], ".")
	}
	treeish := r.URL.Query().Get("tree")
	if treeish == "" {
		treeish = "origin/master"
	}
	mimeType := mime.TypeByExtension(ext)
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", "attachment; filename=archive."+ext)
	err := ga.WriteArchive(pathNoExt, treeish, ext, w)
	if err != nil {
		handleError(err, w, r)
		return
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
	go ga.FetchLoop()

	router := mux.NewRouter()
	if *pypiPrefix != "" && *pypiPrefixDelimiter != "" {
		println("Running pypi server")
		router.HandleFunc("/_pypi-simple", ga.PypiIndexHandler)
		router.HandleFunc("/_pypi-simple/", ga.PypiIndexHandler)
		router.HandleFunc("/_pypi-simple/{package}", ga.PypiPackageHandler)
		router.HandleFunc("/_pypi-simple/{package}/", ga.PypiPackageHandler)
	}
	router.HandleFunc("/{rest:.*}", ga.HttpHandler)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	http.Handle("/", loggedRouter)
	err := http.ListenAndServe(*listen, nil)
	panic(err)

}
