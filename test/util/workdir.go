package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type WorkDir struct {
	Dir string
	Env map[string]string
}

type CommandError struct {
	InternalError error
	Stdout        string
	Stderr        string
}

func (ce *CommandError) Error() string {
	return fmt.Sprintf("%s\n\nstdout:\n%s\n\nstderr:\n%s", ce.InternalError, ce.Stdout, ce.Stderr)
}

func CloneWorkDir(base *WorkDir) (*WorkDir, error) {
	wd, err := CreateWorkDir()
	if err != nil {
		return nil, err
	}
	_, _, err = wd.RunCommmand("git", "clone", base.Dir, wd.Dir)
	if err != nil {
		wd.Remove()
		return nil, err
	}
	return wd, nil
}

func CreateGitWorkDir() (*WorkDir, error) {
	wd, err := CreateWorkDir()
	if err != nil {
		return nil, err
	}
	_, _, err = wd.RunCommmand("git", "init")
	if err != nil {
		wd.Remove()
		return nil, err
	}
	_, _, err = wd.RunCommmand("git", "config", "user.email", "you@example.com")
	if err != nil {
		wd.Remove()
		return nil, err
	}
	_, _, err = wd.RunCommmand("git", "config", "user.name", "Your Name")
	if err != nil {
		wd.Remove()
		return nil, err
	}
	return wd, nil
}

func CreateWorkDir() (*WorkDir, error) {
	dir, err := ioutil.TempDir("", "git-ghost-e2e-test-")
	if err != nil {
		return nil, err
	}
	return &WorkDir{Dir: dir}, nil
}

func (wd *WorkDir) Remove() error {
	return os.RemoveAll(wd.Dir)
}

func (wd *WorkDir) RunCommmand(command string, args ...string) (string, string, error) {
	cmd := exec.Command(command, args...)
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")
	cmd.Dir = wd.Dir
	var env []string
	for _, e := range os.Environ() {
		env = append(env, e)
	}
	for key, val := range wd.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, val))
	}
	env = append(env, fmt.Sprintf("PWD=%s", wd.Dir))
	cmd.Env = env
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		err = &CommandError{
			InternalError: err,
			Stdout:        stdout.String(),
			Stderr:        stderr.String(),
		}
	}
	return stdout.String(), stderr.String(), err
}