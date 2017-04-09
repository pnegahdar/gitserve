package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"runtime/debug"
	"strings"
	"text/template"
)

type Package struct {
	Versions map[string]string
	Name     string
	Dir      string
}

type Packages map[string]Package

const pypiIndexTemplate = `<html>
 <body>
	{{- range $packageName, $package := . -}}
	<a href="/_pypi-simple/{{ $packageName }}">{{ $packageName }}</a></br>
	{{- end -}}
 </body>
</html>
`

const pypiPackageTemplate = `<html>
<body>
	{{- $meta := . -}}

	{{- range $version, $tag := .Versions -}}
	<a href="/{{ $meta.Dir }}.tar.gz?tree={{ $tag }}#egg={{ $meta.Name }}-{{ $version }}<">{{ $meta.Name }}-{{ $version }}</a></br>
	{{- end -}}
</body>
</html>
`

func splitOnceRight(s, delimiter string) []string {
	parts := strings.Split(s, delimiter)
	return []string{strings.Join(parts[0:len(parts)-1], ""), parts[len(parts)-1]}
}

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	debug.PrintStack()
	fmt.Printf("Ran into error: %v", err.Error())
}

func (ga *GitArchive) PypiIndexHandler(w http.ResponseWriter, r *http.Request) {
	packages, err := ga.AvailablePackages()
	if err != nil {
		handleError(err, w, r)
		return
	}
	tmpl, err := template.New("test").Parse(pypiIndexTemplate)
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, packages)
}

func (ga *GitArchive) PypiPackageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	packageName := vars["package"]
	packages, err := ga.AvailablePackages()
	if err != nil {
		handleError(err, w, r)
		return
	}
	packageMeta, ok := packages[packageName]
	if !ok {
		w.WriteHeader(404)
		return
	}
	tmpl, err := template.New("test").Parse(pypiPackageTemplate)
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, packageMeta)
}

func (ga *GitArchive) AvailablePackages() (Packages, error) {
	resp, err := ga.CommandOutput("tag")
	if err != nil {
		return nil, err
	}
	packages := Packages{}
	tags := strings.Split(resp, "\n")
	for _, tag := range tags {
		if strings.HasPrefix(tag, *pypiPrefix) {
			parts := splitOnceRight(tag, *pypiPrefixDelimiter)
			// Per pypi _ is replaced with -
			packageName := strings.Replace(parts[0], "_", "-", -1)
			pack, ok := packages[packageName]
			if !ok {
				pack = Package{Name: packageName, Versions: map[string]string{}, Dir: parts[0]}
			}
			pack.Versions[parts[1]] = tag
			packages[packageName] = pack
		}
	}
	return packages, err

}
