package cmd

import (
	"fmt"

	"github.com/chinx/rsa/setting"
	"github.com/urfave/cli"
)

var List = cli.Command{
	Name:      "list",
	ShortName: "ls",
	Action:    list,
	Usage:     "show remote list",
}

func list(context *cli.Context) {
	args := context.Args()
	group, needMatch := "", false
	if len(args) > 0 {
		group, needMatch = args[0], true
	}
	for key, val := range setting.Root.Remotes {
		equal := val.Name == group
		if !needMatch || equal {
			fmt.Printf("%s:\n", val.Name)
			for k, v := range val.List {
				if v.Name != ""{
					fmt.Printf("    %d-%d | %s-%s | %s\n", key, k, val.Name, v.Name, v.Host)
				}else{
					fmt.Printf("    %d-%d | %s\n", key, k, v.Host)
				}

			}
		}
		if equal {
			break
		}
	}
}
