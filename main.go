package main

import (
	"github.com/docopt/docopt-go"
	"os"
	"strconv"
)

var version = "v0.1"
var defaultPort string = "8000"

func main() {
	usage := `Hafez Restaurant Websever.

Usage:
	hafez [--development][--port=<port> | --bind][--template-dir=<t>][--public-dir=<p>][--static-url=<s>]
	hafez -h | --help
	hafez -v | --version

Options:
	-h --help           Show this screen.
	-v --version        Show version.
	--port=<port>       Port to bind on [default: ` + defaultPort + `].
	--bind              Bind to $PORT.
	--development       Start server in development mode.
	--template-dir=<t>  Directory to serve templates from [default: templates/].
	--public-dir=<p>    Directory to serve assets from [default: public/].
	--static-url=<s>    URL to serve static assets from [default: /assets].`

	//Make sure that the first arg passed to Docopt is the one after `hafez`,
	//in case the program is being run through, say, github.com/codegangsta/gin,
	//where the command would look like `gin hafez --bind`.
	firstArg := 0
	for idx, arg := range os.Args {
		if arg == "hafez" {
			firstArg = idx + 1
		}
	}

	//Use docopts to parse our command-line arguments.
	args, err := docopt.Parse(usage, os.Args[firstArg:], true, "Hafez Restaurant "+version+".", false)
	if err != nil {
		panic(err)
	}

	//Parse whether we're in development mode.
	dev := false
	if args["--development"].(bool) {
		dev = true
	}

	var port int
	var portString string

	//If `--bind` is set, attempt to bind to $PORT. Otherwise use the `--port` argument.
	if args["--bind"].(bool) {
		portString = os.Getenv("PORT")
		port, err = strconv.Atoi(portString)
	} else {
		portString = args["--port"].(string)
		port, err = strconv.Atoi(portString)
	}

	//If $PORT or `--port` could not be resolved to a positive integer, bind to defaultPort.
	if err != nil || !(port > 0) {
		port, _ = strconv.Atoi(defaultPort)
	}

	//Parse the directories we'll use from the command-line args.
	templateDir := args["--template-dir"].(string)
	publicDir := args["--public-dir"].(string)
	staticUrl := args["--static-url"].(string)

	//Feed everything into a new AppOptions instance.
	opts := AppOptions{
		Development: dev,
		Port:        port,
		TemplateDir: templateDir,
		PublicDir:   publicDir,
		StaticUrl:   staticUrl,
	}

	//Run the app.
	NewApp(opts).Run()
}
