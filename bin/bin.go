package bin

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func RunBin(app string, args []string) {
	var (
		p    string
		info fs.FileInfo
		err  error
		cmd  *exec.Cmd
	)

	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s exit: %v\n", app, err)
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "%s exit\n", app)
		}
	}()

	p = filepath.Join(filepath.Dir(os.Args[0]), app)

	if info, err = os.Stat(p); err != nil {
		// if errors.Is(err, os.ErrNotExist)
		return
	}

	if !info.Mode().IsRegular() {
		err = fmt.Errorf("not regular file: %s", p)
		return
	}
	// info.Mode().Perm()

	cmd = exec.Command(p, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
}
