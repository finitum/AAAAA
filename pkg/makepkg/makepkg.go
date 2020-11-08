package makepkg

import (
	"github.com/finitum/AAAAA/internal/cmdutil"
)

func Build() (string, string, error) {
	return cmdutil.RunCommand("/usr/bin/makepkg",
		"--syncdeps",
		"--install",
		"--nocheck",
		"--noconfirm",
		//"PKGDEST=" + pkgdest,
		"PKGEXT="+".pkg.tar.zst",
	)
}
