# httpwatch

![Go Build and Test](https://github.com/bentekkie/httpwatch/actions/workflows/go.yml/badge.svg) [![PkgGoDev](https://pkg.go.dev/badge/github.com/bentekkie/httpwatch)](https://pkg.go.dev/github.com/bentekkie/httpwatch)

A basic equivilant to the `watch` command that serves the output via an http server.

Functionality is mostly the same for flags -n, -t and -c. --address is used to specify a tcp ip and port.

## Installation

```
go install github.com/bentekkie/httpwatch/cmd/httpwatch
```

## Usage

```
$ httpwatch -h
Usage: httpwatch [-flags] -- command ...
  -a, --address string      Address to serve the output of the command to. (default "127.0.0.1:8000")
  -c, --color               Interpret ANSI color and style sequences.
  -n, --interval duration   Specify update interval.  The command will not allow quicker than 0.1 second interval, in which the smaller values are converted. The WATCH_INTERVAL environment can be used to persistently set a non-default interval (following the same rules and formatting). (default 2s)
  -t, --no-title            Turn off the header showing the interval, command, and current time at the top of the display, as well as the following blank line.
```

```
$ httpwatch -n 0.5s -- date
Listening at http://127.0.0.1:8000
```
![Screenshot](6T5sBfyXnDKKddu.png?raw=true "Title")
