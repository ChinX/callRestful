package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chinx/rsa/handler"
	"github.com/urfave/cli"
	"log"
)

var remoteDirFlag = cli.StringFlag{
	Name:  "remote, r",
	Usage: "remote open dir",
}

var currentDirFlag = cli.StringFlag{
	Name:  "current, c",
	Usage: "current open dir",
}

var SFTP = cli.Command{
	Name:   "sftp",
	Action: openSFTP,
	Usage:  "open new sftp",
	Flags: []cli.Flag{
		remoteDirFlag,
		currentDirFlag,
	},
}

func openSFTP(context *cli.Context) {
	r, err := remoteFromContext(context)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := handler.ConnectSSH(r); err != nil {
		if strings.Contains(err.Error(), "unable to authenticate") {
			inputReader := bufio.NewReader(os.Stdin)
			fmt.Printf("Passwd wrong! please enter new: ")
			input, err3 := inputReader.ReadString('\n')
			if err3 != nil || input == "" {
				err = fmt.Errorf("Input passwd \"%s\"is error: %s ", input, err3)
				return
			}
			r.Passwd = input
			log.Println(input)
			err = handler.ConnectSSH(r)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
