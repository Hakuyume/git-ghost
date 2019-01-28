package git

import (
	"bytes"
	"errors"
	"fmt"
	"git-ghost/pkg/util"
	"io"
	"os"
	"os/exec"
)

func CreateDiffBundleFile(dir, filepath, fromRefspec, toRefspec string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command("git", "-C", dir, "format-patch", "--binary", "--stdout", fmt.Sprintf("%s..%s", fromRefspec, toRefspec))
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer reader.Close()
	err = cmd.Start()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return err
	}

	total := 0
	buf := make([]byte, 1024, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			d := buf[:n]
			_, err = f.Write(d)
			if err != nil {
				return err
			}
			total += n
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func ApplyDiffBundleFile(dir, filepath string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "am", filepath),
	)
}

func CreateDiffPatchFile(dir, filepath, refspec string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command("git", "-C", dir, "diff", "--binary", refspec)
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer reader.Close()
	err = cmd.Start()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return err
	}

	total := 0
	buf := make([]byte, 1024, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			d := buf[:n]
			_, err = f.Write(d)
			if err != nil {
				return err
			}
			total += n
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func ApplyDiffPatchFile(dir, filepath string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "apply", filepath),
	)
}
