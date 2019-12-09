package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/axetroy/denox/internal/deno"
	"github.com/axetroy/denox/internal/signals"
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

	defer d.Clean()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quit
		// if receive exit signal
		err = d.Clean()

		if err != nil {
			log.Printf("%+v\n", err)
		}

		os.Exit(1)
	}()

	executablePath, err := d.Download()

	// if download success. we don't need to listen quit signal
	signal.Stop(quit)

	if err != nil {
		err = errors.Wrap(err, "get executable path fail")
		return
	}

	signalProxy := make(chan os.Signal)
	signal.Notify(signalProxy, signals.AllSignals...)

	go func() {
		s := <-signalProxy
		if cmd != nil && !cmd.ProcessState.Exited() {
			if err = cmd.Process.Signal(s); err != nil {
				err = errors.Wrapf(err, "send signal `%s` fail", s)
			}
		}
	}()

	if err := os.Setenv("DENO_DIR", d.DenoDir); err != nil {
		err = errors.Wrapf(err, "set env $DENO_DIR=%s fail", d.DenoDir)
		return
	}

	cmd = exec.Command(executablePath, denoArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			denoExitCode = exitError.ExitCode()
		} else {
			err = errors.Wrap(err, "run command fail")
		}
		return
	}
}
