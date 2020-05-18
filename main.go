package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/logrusorgru/aurora"
	"github.com/rebeccajae/grterm/pkg/ttyrec"
	"github.com/riywo/loginshell"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

var au aurora.Aurora

func term(cmd, out string) error {
	rec, err := ttyrec.NewTTYRecorder(out)
	if err != nil {
		return err
	}
	defer rec.Close()

	c := exec.Command(cmd)
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	fmt.Println(au.Sprintf(au.Green("Recording %s to %s"), cmd, out))
	defer func() {
		_ = ptmx.Close()
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = terminal.Restore(int(os.Stdin.Fd()), oldState)
	}()

	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
	}()

	stdoutw := io.MultiWriter(os.Stdout, rec)
	_, _ = io.Copy(stdoutw, ptmx)

	return nil
}

func main() {
	shell, err := loginshell.Shell()
	if err != nil {
		log.Fatal(err)
	}

	disableColors := flag.Bool("nocolor", false, "Disable message colors")
	cmd := flag.String("command", shell, "Command to execute as shell")
	output := flag.String("output", "rec.ttyrec", "Save path of recording")
	flag.Parse()

	au = aurora.NewAurora(!*disableColors)
	if err := term(*cmd, *output); err != nil {
		fmt.Println(au.Sprintf(au.Red("Error recording: %s"), err))
	}
	fmt.Println(au.Sprintf(au.Green("Finished recording %s to %s"), *cmd, *output))
}
