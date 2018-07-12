package cmd

import (
	"fmt"
	"os"

	"github.com/chinx/gtty/handler"
	"github.com/chinx/gtty/model"
	"github.com/chinx/gtty/setting"
	"github.com/urfave/cli"
)

var List = cli.Command{
	Name:      "list",
	ShortName: "ls",
	Action: func(context *cli.Context) {
		group := model.DefaultGroup
		if args := context.Args(); len(args) > 0 {
			group = args[0]
		}
		handler.ListRemotes(setting.Root.Remotes, group)
	},
	Usage: "show remote list",
}

var SSH = cli.Command{
	Name: "ssh",
	Action: func(context *cli.Context) {
		r, err := remoteFromContext(context)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		handler.OpenSSH(r)
	},
	Usage: "open new ssh",
}

var SFTP = cli.Command{
	Name: "sftp",
	Action: func(context *cli.Context) {
		r, err := remoteFromContext(context)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		handler.OpenSFTP(r)
	},
	Usage: "open new sftp",
	Flags: []cli.Flag{
		remoteDirFlag,
		currentDirFlag,
	},
}
