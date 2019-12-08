package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/axetroy/denox/internal/deno"
	"github.com/pkg/errors"
)

func main() {
	args := os.Args

	var (
		err          error
		denoArgs     = args[1:]
		denoVersion  *string
		denoExitCode int
		cmd          *exec.Cmd
	)

	defer func() {
		if err != nil {
			log.Fatalf("%+v\n", err)
		}

		os.Exit(denoExitCode)
	}()

	if v := os.Getenv("DENO_VERSION"); v != "" {
		denoVersion = &v
	}

	d, err := deno.New(denoVersion)

	if err != nil {
		return
	}

	defer func() {
		err = d.Clean()
	}()

	executablePath, err := d.Download()

	if err != nil {
		err = errors.Wrap(err, "get executable path fail")
		return
	}

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		s := <-quit
		if cmd != nil {
			if err = cmd.Process.Signal(s); err != nil {
				err = errors.Wrapf(err, "send signal `%s` fail", s)
			}
		}
	}()

	cmd = exec.Command(executablePath, denoArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(cmd.Env, fmt.Sprintf("DENO_DIR=%s", d.DenoDir))

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			denoExitCode = exitError.ExitCode()
		} else {
			err = errors.Wrap(err, "run command fail")
		}
		return
	}
}
