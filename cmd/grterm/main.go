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

	"github.com/TwinProduction/go-color"
	"github.com/rebeccajae/grterm/pkg/ttyrec"
	"github.com/riywo/loginshell"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

var disableResize bool

func term(cmd string, f io.Writer) error {
	rec := ttyrec.NewTTYRecorder(f)

	c := exec.Command(cmd)
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	fmt.Println(color.Ize(color.Green, fmt.Sprintf("Recording %s", cmd)))
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
			if !disableResize {
				y, x, err := pty.Getsize(ptmx)
				if err != nil {
					log.Printf("error resizing recorded pty: %s", err)
				}
				code := fmt.Sprintf("\x1b[8;%d;%dt", y, x)
				rec.Write([]byte(code))
			}
		}
	}()
	ch <- syscall.SIGWINCH

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	defer func() {
		_ = terminal.Restore(int(os.Stdin.Fd()), oldState)
	}()

	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
	}()

	stdoutw := io.MultiWriter(os.Stdout, rec)
	_, err = io.Copy(stdoutw, ptmx)

	return err
}

func main() {
	shell, err := loginshell.Shell()
	if err != nil {
		log.Fatal(err)
	}

	cmd := flag.String("command", shell, "Command to execute as shell")
	output := flag.String("output", "rec.ttyrec", "Save path of recording")
	flag.BoolVar(&disableResize, "noresize", false, "Disables insertion of resize escape codes")
	flag.Parse()

	f, err := os.Create(*output)
	if err != nil {
		fmt.Println(color.Ize(color.Red, fmt.Sprintf("Error recording: %s", err)))
	}
	defer f.Close()
	fmt.Println(color.Ize(color.Green, fmt.Sprintf("Recording to %s", *output)))

	if err := term(*cmd, f); err != nil {
		fmt.Println(color.Ize(color.Red, fmt.Sprintf("Error recording: %s", err)))
	}
	fmt.Println(color.Ize(color.Green, fmt.Sprintf("Finished recording %s to %s", *cmd, *output)))
}
