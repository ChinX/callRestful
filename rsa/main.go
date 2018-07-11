package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chinx/rsa/cmd"
	"github.com/chinx/rsa/setting"
	"github.com/urfave/cli"
)

func main() {
	err := setting.OpenConfig("./config/config.yaml")
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "remote source access for golang"
	app.Commands = cli.Commands{
		//cmd.Name,
		cmd.List,
		cmd.SSH,
		cmd.SFTP,
	}

	if err = app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
