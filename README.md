# `grterm` - A Terminal Recorder
A really simple terminal recorder compatible with the ttyrec format written in Go.

![Run Tests](https://github.com/rebeccajae/grterm/workflows/Run%20Tests/badge.svg?branch=default)
![Build](https://github.com/rebeccajae/grterm/workflows/Build/badge.svg?branch=default)

## Usage
You can bring this up by doing `grterm --help`.
```
Usage:
  -command string
    	Command to execute as shell (default "<YOUR_LOGIN_SHELL>")
  -noresize
    	Disables insertion of resize escape codes
  -output string
    	Save path of recording (default "rec.ttyrec")
```

## Installation 
Grab a build from the default branch over in the Build action. I try to 
make sure that the default branch is always prod-ready.

## Resizability
Some terminal emulators allow you to configure escape-code driven resizes. 
The specific configuration for your terminal emulator may vary.

For iTerm (which is what I use), it is by default disabled by the setting
```
Profile > Terminal > Disable session-initiated window resizing
```

Note that not all terminal emulators know what to do with this, so
be careful. You can disable inserting them by using the `--noresize` flag.

If a terminal emulator is not compatible with these, it may render the escape
sequence. I know iTerm works, and there's no reason xterm wouldn't either, but
YMMV.

## Reusable Bits
`pkg/ttyrec` implements a writer that is compatible with the ttyrec format.
I use a pty library and just multiwrite to a ttyrec writer.
