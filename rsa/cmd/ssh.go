package cmd

import (
	"fmt"
	"os"
	"strings"
	"bufio"

	"github.com/chinx/rsa/handler"
	"github.com/urfave/cli"
	"log"
)

var SSH = cli.Command{
	Name:   "ssh",
	Action: openSSH,
	Usage:  "open new ssh",
}

func openSSH(context *cli.Context) {
	r, err := remoteFromContext(context)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := handler.ConnectSSH(r); err != nil {
		if strings.Contains(err.Error(), "unable to authenticate") {
			inputReader := bufio.NewReader(os.Stdin)
			fmt.Printf("Passwd wrong! please enter new: ")
			input, err1 := inputReader.ReadString('\n')
			if err1 != nil || input == "" {
				fmt.Errorf("Input passwd \"%s\"is error: %s ", input, err1)
				os.Exit(1)
			}
			r.Passwd = input
			log.Println(input)
			err = handler.ConnectSSH(r)
			log.Println(r, err)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
