#Hafez restaurant website repository

This is the GitHub repository for the Hafez Restaurant website. The site is served with Martini. It uses standard Go HTML templates.

##Config
To configure the dokku domains on live, use the following command (Note that this requires the dokku domains extension):

`dokku domains:set hafez hafezrestaurant.co.uk hafez.fmdud.com`
You can swap out fmdud.com for your own domain if this is being hosted elsewhere.

##Development
* To compile the program, make sure your `$GOPATH` is set and run `go install github.com/fmd/hafez`.
`hafez` will run the server. Make sure you are in the root directory of this repo when the command is run.

##Deployment
* Ensure that the git remote is set up properly: `git remote add live dokku@fmdud.com:hafez`.
* Push the `master` branch to `live`: `git push live master`.
