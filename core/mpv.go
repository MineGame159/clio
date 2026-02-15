package core

import (
	"os/exec"
)

func OpenMpv(mediaName string, url string) {
	cmd := exec.Command(
		"mpv",
		"--terminal=no",
		"--alang=en",
		"--slang=en",
		"--window-maximized=yes",
		"--force-media-title="+mediaName,
		url,
	)

	cmd.SysProcAttr = getSysProcAttr()

	err := cmd.Start()
	if err != nil {
		panic(err.Error())
	}

	err = cmd.Process.Release()
	if err != nil {
		panic(err.Error())
	}
}
