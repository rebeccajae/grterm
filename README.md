# `grterm` - A Terminal Recorder
A really simple terminal recorder compatible with the ttyrec format written in Go.

## Usage
You can bring this up by doing `grterm --help`.
```
Usage:
  -command string
    	Command to execute as shell (default "<YOUR_LOGIN_SHELL>")
  -nocolor
    	Disable message colors
  -output string
    	Save path of recording (default "rec.ttyrec")
```

## Installation 
```
go get github.com/rebeccajae/grterm
```

## Reusable Bits
`pkg/ttyrec` implements a writer that is compatible with the ttyrec format.
I use a pty library and just multiwrite to a ttyrec writer.
