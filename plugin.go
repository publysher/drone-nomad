package main

import (
	"fmt"
	"os"
	"errors"
	"io/ioutil"
	"text/template"
	"bytes"
	"os/exec"
)

type (
	Repo struct {
		Owner   string
		Name    string
		Link    string
		Avatar  string
		Branch  string
		Private bool
		Trusted bool
	}

	Build struct {
		Number   int
		Event    string
		Status   string
		Deploy   string
		Created  int64
		Started  int64
		Finished int64
		Link     string
	}

	Commit struct {
		Remote  string
		Sha     string
		Ref     string
		Link    string
		Branch  string
		Message string
		Author  Author
	}

	Author struct {
		Name   string
		Email  string
		Avatar string
	}

	Config struct {
		NomadAddr   string
		Job         string
		UseTemplate bool
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Commit Commit
		Config Config
	}
)

func (p Plugin) Exec() error {
	if p.Config.Job == "" {
		return errors.New("no job file specified")
	}
	if p.Config.NomadAddr == "" {
		p.Config.NomadAddr = "http://nomad.service.consul:4646"
	}

	jobFile, err := os.Open(p.Config.Job)
	if err != nil {
		return err
	}
	defer jobFile.Close()

	jobSpecBytes, err := ioutil.ReadAll(jobFile)
	if err != nil {
		return fmt.Errorf("could not read job file %s: %v", p.Config.Job, err)
	}

	var jobSpec *bytes.Buffer
	if p.Config.UseTemplate {
		tmpl, err := template.New("job").Parse(string(jobSpecBytes))
		if err != nil {
			return err
		}
		jobSpec = new(bytes.Buffer)
		if err := tmpl.Execute(jobSpec, p); err != nil {
			return err
		}
	} else {
		jobSpec = bytes.NewBuffer(jobSpecBytes)
	}

	cmd := exec.Command("/bin/nomad", "run", "-address", p.Config.NomadAddr, "-verbose", "-")
	cmd.Stdin = jobSpec
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	return err
}
