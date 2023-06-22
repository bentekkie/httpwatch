# httpwatch

A basic equivilant to the `watch` command that serves the output via an http server.

Functionality is mostly the same for flags -n, -t and -c. --address is used to specify a tcp ip and port.

## Installation

```
go install github.com/bentekkie/httpwatch/cmd/httpwatch
```


## Usage

```
$ httpwatch -n 0.5s -- date
Listening at http://127.0.0.1:8000
```
![Screenshot](6T5sBfyXnDKKddu.png?raw=true "Title")
