package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/bentekkie/httpwatch/web"
	"github.com/robert-nix/ansihtml"
)

var (
	interval = flag.Duration("n", 2*time.Second, "Specify update interval.  The command will not allow quicker than 0.1 second interval, in which the smaller values are converted. The WATCH_INTERVAL environment can be used to persistently set a non-default interval (following the same rules and formatting).")
	noTitle  = flag.Bool("t", false, "Turn off the header showing the interval, command, and current time at the top of the display, as well as the following blank line.")
	color    = flag.Bool("c", false, "Interpret ANSI color and style sequences.")
	address  = flag.String("address", "127.0.0.1:8000", "Address to serve the output of the command to.")
)

const minInterval = time.Second / 10

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %v [-flags] -- command ...\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	if *interval < minInterval {
		*interval = minInterval
	}
	c := &Command{cmd: flag.Args()}
	if len(c.cmd) == 0 {
		log.Fatal("No command specified")
	}
	err := mime.AddExtensionType(".js", "text/javascript")
	if err != nil {
		log.Fatalf("Error in mime js: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go c.Run(ctx)

	r := http.NewServeMux()
	r.HandleFunc("/update", c.WriteUpdate)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		c.WriteIndex(w, r)
	})

	fmt.Printf("Listening at http://%v\n", *address)
	server := &http.Server{Addr: *address, Handler: r}
	go func() {
		if err := http.ListenAndServe(*address, r); err != nil {
			log.Println(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	sctx, scancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer scancel()
	if err := server.Shutdown(sctx); err != nil {
		log.Fatalf("Error shutting down http server: %v", err)
	}
}

// Command describes the command to be run
type Command struct {
	cmd       []string
	mu        sync.RWMutex
	updatedAt time.Time
	output    template.HTML
	err       error
}

// Run executues the command every interval until the context is canceled
func (c *Command) Run(ctx context.Context) {
	next := time.Now().Add(*interval)
	c.update(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Until(next)):
			next = next.Add(*interval)
			c.update(ctx)
		}
	}
}

func (c *Command) update(ctx context.Context) {
	bs, err := exec.CommandContext(ctx, c.cmd[0], c.cmd[1:]...).CombinedOutput()
	if ctx.Err() != nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if *color {
		bs = ansihtml.ConvertToHTML(bs)
	}
	raw := strings.ReplaceAll(string(bs), "\n", "<br/>")
	raw = strings.ReplaceAll(raw, "\r\n", "<br/>")
	c.output = template.HTML(raw)
	c.err = err
	c.updatedAt = time.Now()
}

type tplData struct {
	Title     string
	Cmd       string
	UpdatedAt time.Time
	Interval  time.Duration
	Output    template.HTML
	Err       error
	NoTitle   bool
}

func (c *Command) tplData() tplData {
	return tplData{
		Title:     c.cmd[0],
		Cmd:       strings.Join(c.cmd, " "),
		Output:    c.output,
		Err:       c.err,
		UpdatedAt: c.updatedAt,
		Interval:  *interval,
		NoTitle:   *noTitle,
	}
}

// WriteIndex writes the index html page with the command results
func (c *Command) WriteIndex(w http.ResponseWriter, r *http.Request) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	err := web.Template.Lookup("index.html").Execute(w, c.tplData())
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

// WriteUpdate sends the updated command result
func (c *Command) WriteUpdate(w http.ResponseWriter, r *http.Request) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var content strings.Builder
	err := web.Template.Lookup("content.html").Execute(&content, c.tplData())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"interval": (*interval).Milliseconds(),
		"content":  content.String(),
	})

	if err != nil {
		log.Printf("Error: %v", err)
	}
}
