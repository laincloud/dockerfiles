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
do:{{range $context, $_ := .Diff}} {{$context | escapeSlash}}{{end}}

{{range $context, $image := .All}}{{$context | escapeSlash}}: {{if $image.From.In $.Diff}}{{$image.From.Context $.All | escapeSlash}}{{end}}
	$(buildImage) {{range $image.Tags}}{{printf "-t %s/%s:%s" $image.Organization $image.Repository .}} {{end}}{{$image.Context}}

{{end -}}

{{else if eq (showCommand .Args.Command) "pull"}}
do:{{range $context, $_ := .Diff}} {{$context | escapeSlash}}{{end}}

{{range $context, $image := .All}}{{$context | escapeSlash}}: {{if $image.From.In $.Diff}}{{$image.From.Context $.All | escapeSlash}}{{end}}
	{{range $image.Tags}}$(pullImage) {{if $.Args.RegistryHost}}{{printf "%s/" $.Args.RegistryHost}}{{end}}{{printf "%s/%s:%s" $.Args.Organization $image.Repository .}}
	{{end}}

{{end -}}

{{else if eq (showCommand .Args.Command) "push"}}
do: {{range $context, $_ := .Diff}} {{$context | escapeSlash}}{{end}}

{{range $context, $image := .All}}{{$context | escapeSlash}}: {{if $image.From.In $.Diff}}{{$image.From.Context $.All | escapeSlash}}{{end}}
	{{range $image.Tags}}$(pushImage) {{if $.Args.RegistryHost}}{{printf "%s/" $.Args.RegistryHost}}{{end}}{{printf "%s/%s:%s" $.Args.Organization $image.Repository .}}
	{{end}}

{{end -}}

{{else if eq (showCommand .Args.Command) "retag"}}
do: {{range $context, $_ := .Diff}} {{$context | escapeSlash}}{{end}}

{{range $context, $image := .All}}{{$context | escapeSlash}}: {{if $image.From.In $.Diff}}{{$image.From.Context $.All | escapeSlash}}{{end}}
	{{range $image.Tags}}$(retagImage) {{if $.Args.OldRegistryHost}}{{printf "%s/" $.Args.OldRegistryHost}}{{end}}{{printf "%s/%s:%s" $.Args.OldOrganization $image.Repository .}} {{if $.Args.NewRegistryHost}}{{printf "%s/" $.Args.NewRegistryHost}}{{end}}{{printf "%s/%s:%s" $.Args.NewOrganization $image.Repository .}}
	{{end}}

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

	diffImages, err := GetContext2Images(diffFiles)
	if err != nil {
		return nil, err
	}

	allFiles, err := Walk(".")
	if err != nil {
		return nil, err
	}

	allImages, err := GetContext2Images(allFiles)
	if err != nil {
		return nil, err
	}

	for _, image := range allImages {
		for _, diffImage := range diffImages {
			if image.IsFrom(diffImage, allImages) {
				diffImages[image.Context] = image
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
		"escapeSlash": EscapeSlash,
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
