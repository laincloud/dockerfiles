package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

const (
	makefileTemplate = `
buildImage := docker build
pullImage := docker pull
pushImage := docker push
retagImage := docker tag

.PHONY: do

{{if eq (showCommand .Args.Command) "build"}}
do:{{range $_, $v := .Diff}} {{$v.Name | escapeColon}}{{end}}

{{range $_, $v := .All}}{{$v.Name | escapeColon}}: {{if $v.From.In $.Diff}}{{$v.From.Name | escapeColon}}{{end}}
	$(buildImage) {{printf "-t %s %s" $v.Name $v.Context}}

{{end -}}

{{else if eq (showCommand .Args.Command) "pull"}}
do:{{range $_, $v := .Diff}} {{$v.Name | escapeColon}}{{end}}

{{range $_, $v := .All}}{{$v.Name | escapeColon}}: {{if $v.From.In $.Diff}}{{$v.From.Name | escapeColon}}{{end}}
	$(pullImage) {{if $.Args.RegistryHost}}{{printf "%s/" $.Args.RegistryHost}}{{end}}{{printf "%s/%s:%s" $.Args.Organization $v.Repository $v.Tag}}

{{end -}}

{{else if eq (showCommand .Args.Command) "push"}}
do: {{range $_, $v := .Diff}} {{$v.Name | escapeColon}}{{end}}

{{range $_, $v := .All}}{{$v.Name | escapeColon}}: {{if $v.From.In $.Diff}}{{$v.From.Name | escapeColon}}{{end}}
	$(pushImage) {{if $.Args.RegistryHost}}{{printf "%s/" $.Args.RegistryHost}}{{end}}{{printf "%s/%s:%s" $.Args.Organization $v.Repository $v.Tag}}

{{end -}}

{{else if eq (showCommand .Args.Command) "retag"}}
do: {{range $_, $v := .Diff}} {{$v.Name | escapeColon}}{{end}}

{{range $_, $v := .All}}{{$v.Name | escapeColon}}: {{if $v.From.In $.Diff}}{{$v.From.Name | escapeColon}}{{end}}
	$(retagImage) {{if $.Args.OldRegistryHost}}{{printf "%s/" $.Args.OldRegistryHost}}{{end}}{{printf "%s/%s:%s" $.Args.OldOrganization $v.Repository $v.Tag}} {{if $.Args.NewRegistryHost}}{{printf "%s/" $.Args.NewRegistryHost}}{{end}}{{printf "%s/%s:%s" $.Args.NewOrganization $v.Repository $v.Tag}}

{{end -}}
{{end -}}
`
)

// Args denotes command arguments
type Args struct {
	Command         Command
	Commit1         string
	Commit2         string
	NewOrganization string
	NewRegistryHost string
	OldOrganization string
	OldRegistryHost string
	Organization    string
	Password        string
	RegistryHost    string
	User            string
}

// MakefileData is for makefileTemplate
type MakefileData struct {
	All  map[string]Image
	Diff map[string]Image
	Args Args
}

// newMakefileData create *MakefileData
func newMakefileData(args Args) (*MakefileData, error) {
	diffFiles, err := Diff(args.Commit1, args.Commit2)
	if err != nil {
		return nil, err
	}

	allFiles, err := Walk(".")
	if err != nil {
		return nil, err
	}

	diffImages, err := parseImages(diffFiles)
	if err != nil {
		return nil, err
	}

	allImages, err := parseImages(allFiles)
	if err != nil {
		return nil, err
	}
	for _, image := range allImages {
		for _, diffImage := range diffImages {
			if image.IsFrom(diffImage, allImages) {
				diffImages[image.Name()] = image
			}
		}
	}

	return &MakefileData{
		All:  allImages,
		Diff: diffImages,
		Args: Args{
			Command:         args.Command,
			Commit1:         args.Commit1,
			Commit2:         args.Commit2,
			NewOrganization: args.NewOrganization,
			NewRegistryHost: args.NewRegistryHost,
			OldOrganization: args.OldOrganization,
			OldRegistryHost: args.OldRegistryHost,
			Organization:    args.Organization,
			Password:        args.Password,
			RegistryHost:    args.RegistryHost,
			User:            args.User,
		},
	}, nil
}

// generateMakefile generate makefile for build/pull/push/retag images
func generateMakefile(args Args) (string, error) {
	data, err := newMakefileData(args)
	if err != nil {
		return "", err
	}

	funcMap := template.FuncMap{
		"escapeColon": EscapeColon,
		"showCommand": ShowCommand,
	}
	t := template.Must(template.New("makefile").Funcs(funcMap).Parse(makefileTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, *data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Make do make
func Make(args Args) error {
	makefile, err := generateMakefile(args)
	if err != nil {
		return err
	}

	fmt.Printf("makefile:\n%s\n", makefile)
	doMake := exec.Command("make", "-f", "-", "do")
	buf := bytes.NewBufferString(makefile)
	doMake.Stdin = buf
	doMake.Stdout = os.Stdout
	doMake.Stderr = os.Stdout
	return doMake.Run()
}
