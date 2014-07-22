package main

import (
	"os"
	"strconv"
	"github.com/docopt/docopt-go"
)

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
	--port=<port>       Port to bind on [default: `+ defaultPort +`].
	--bind              Bind to $PORT.
	--development       Start server in development mode.
	--template-dir=<t>  Directory to serve templates from [default: templates/].
	--public-dir=<p>    Directory to serve assets from [default: public/].
	--static-url=<s>    URL to serve static assets from [default: /assets].`

	args, err := docopt.Parse(usage, nil, true, "Hafez Restaurant v0.1", false)
	if err != nil {
		panic(err)
	}

	dev := false
	if args["--development"].(bool) {
		dev = true
	}

	var port int
	var portString string

	if args["--bind"].(bool) {
		portString = os.Getenv("PORT")
		port, err = strconv.Atoi(portString)
	} else {
		portString = args["--port"].(string)
		port, err = strconv.Atoi(portString)
	}

	if err != nil {
		port, _ = strconv.Atoi(defaultPort)
	}

	templateDir := args["--template-dir"].(string)
	publicDir := args["--public-dir"].(string)
	staticUrl := args["--static-url"].(string)

	opts := AppOptions{
		Development: dev,
		Port:		 port,
		TemplateDir: templateDir,
		PublicDir:	 publicDir,
		StaticUrl:	 staticUrl,
	}

	NewApp(opts).Run()
}
