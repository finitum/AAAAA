package makepkg

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os/exec"
)

func Build() error {
	cmd := exec.Command("/usr/bin/makepkg",
		"--syncdeps",
		"--install",
		"--nocheck",
		"--noconfirm",
		//"PKGDEST=" + pkgdest,
		"PKGEXT="+".pkg.tar.zst",
	)

	errp, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer errp.Close()

	stdp, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdp.Close()

	if err := cmd.Run(); err != nil {
		errb, _ := ioutil.ReadAll(errp)
		stdb, _ := ioutil.ReadAll(stdp)

		errStr := fmt.Sprintf("\nRunning makepkg failed\nstdout: %s \n\n stderr: %s \n", stdb, errb)

		return errors.Wrap(err, errStr)
	}

	return nil
}
