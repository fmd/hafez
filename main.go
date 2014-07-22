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
	hafez [--development][--port=<port> | --bind]
	hafez -h | --help
	hafez -v | --version

Options:
	-h --help       Show this screen.
	-v --version    Show version.
	--port=<port>   Port to bind on [default: `+ defaultPort +`].
	--bind          Bind to $PORT.
	--development   Start server in development mode.`

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

	opts := AppOptions{
		Development: dev,
		Port:		 port,
		TemplateDir: "templates/",
		PublicDir:	 "public/",
		StaticUrl:	 "/assets",
	}

	NewApp(opts).Run()
}
