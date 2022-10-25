package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/floppyisadog/foauthserver/client/cmd"
)

var (
	cliApp     *cli.App
	configFile string
)

func init() {
	// Initialize a CLI app
	cliApp = cli.NewApp()
	cliApp.Name = "foauth-tool"
	cliApp.Usage = "fauth-tool for load data"
	cliApp.Author = "Ameng"
	cliApp.Email = "ameng2001.liu@gmail.com"
	cliApp.Version = "0.5.0"
	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "configFile",
			Value:       "foauthtool.conf",
			Destination: &configFile,
		},
	}
}

func main() {
	// comm := tars.NewCommunicator()
	// obj := fmt.Sprintf("floppyisadog.foauthserver.foauthObj@tcp -h 127.0.0.1 -p 10015 -t 60000")
	// app := new(floppyisadog.Foauth)
	// comm.StringToProxy(obj, app)
	// var out, i int32
	// i = 123
	// ret, err := app.Add(i, i*2, &out)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(ret, out)
	// Set the CLI app commands
	cliApp.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "run migrations",
			Action: func(c *cli.Context) error {
				return cmd.Migrate(configFile)
			},
		},
		{
			Name:  "loaddata",
			Usage: "load data from fixture",
			Action: func(c *cli.Context) error {
				return cmd.LoadData(c.Args(), configFile)
			},
		},
	}

	// Run the CLI app
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
