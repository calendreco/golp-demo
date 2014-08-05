package main

import (
	"github.com/calendreco/golp"
	"github.com/calendreco/golp-contrib/hash"
	"github.com/calendreco/golp-contrib/less"
	"github.com/calendreco/golp-contrib/min"
	// "github.com/calendreco/golp-contrib/reload"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
)

func development(c *cli.Context) {

	// POC - 
	builder := gin.NewBuilder(c.GlobalString("path"), c.GlobalString("bin"))
	runner := gin.NewRunner(filepath.Join(wd, builder.Binary()), c.Args()...)

	// Emulate livereload
	golp.Watch("**.go").Pipe(func() {
		err := builder.Build()
		_, err := runner.Run()
	})

	// Compile our less on the fly
	golp.Watch("assets/less/*.less").Pipe(Less).Dest("assets/less/*.css")

}

func deploy(c *cli.Context) {
	hashAssets = hash.Assets{}
	less := golp.Src("asset/less/*.less").
		Pipe(less.Less).
		Pipe(min.Min).
		Pipe(hash.Hash(&hashAssets)).
		Write("dist/")
	html := golp.Src("template/index.html").
		Pipe(hash.ReplacePaths(&hashAssets)).
		Write("dist/index.html")

	less.Pipe(html).Pipe(func() {
		exec.Command("git subtree push --prefix dist origin gh-pages")
	})
}

func main() {
	app := cli.NewApp()
	app.Name = "symphony api"
	app.Usage = "The api behind symphony"
	flags := []cli.Flag{
		cli.StringFlag{"env,e", "", "the enviroment to run with"},
	}
	app.Commands = []cli.Command{
		{
			Name:   "development",
			Usage:  "start the web process",
			Action: development,
			Flags:  flags,
		},
		{
			Name:   "deploy",
			Usage:  "compile and push to s3",
			Action: deploy,
			Flags:  flags,
		},
	}
	app.Run(os.Args)
}
