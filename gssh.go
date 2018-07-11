package cmd

import (
	"fmt"
	"net"

	"github.com/urfave/cli"
	"github.com/chinx/ctty/setting"
)

var groupFlag = cli.StringFlag{
	Name:  "group, g",
	Usage: "list in group",
}

var List = cli.Command{
	Name:      "list",
	ShortName: "ls",
	Action:    list,
	Flags: []cli.Flag{
		groupFlag,
	},
}

func list(context *cli.Context) {
	group := context.String(groupFlag.Name)
	matchGroup := group != ""
	for key, val := range setting.RemoteList {
		equal := val.Name == group
		if !matchGroup || equal {
			fmt.Printf("%s:\n", val.Name)
			for k, v := range val.List {
				fmt.Printf("    %d-%d : %s\n", key, k, v.Host)
			}
		}
		if equal {
			break
		}
	}
}

var SSH = cli.Command{
	Name:      "ssh",
	ShortName: "sh",
	Action:    openSSH,
	Flags: []cli.Flag{
		hostFlag,
		portFlag,
		userFlag,
		passwdFlag,
	},
}

type SSHRemote struct {
	*remote
}

func openSSH(context *cli.Context) {
	args := context.Args()
	r, err := remoteFromFlag(context)
	if len(args) != 0 {

	}
}

var remoteDirFlag = cli.StringFlag{
	Name:  "remote, r",
	Usage: "remote open dir",
}

var currentDirFlag = cli.StringFlag{
	Name:  "current, c",
	Usage: "current open dir",
}

var SFTP = cli.Command{
	Name:      "sftp",
	ShortName: "cp",
	Action:    openSFTP,
	Flags: []cli.Flag{
		hostFlag,
		portFlag,
		userFlag,
		passwdFlag,
		remoteDirFlag,
		currentDirFlag,
	},
}

func openSFTP(context *cli.Context) {
	cmdLine := context.Args()[0]
	if cmdLine == "" {
		cmdLine = "default"
	}
	log.Println(context.Command.Name)
	log.Println(context.Args())
}

var hostFlag = cli.StringFlag{
	Name:  "host, H",
	Usage: "remote host",
}

var portFlag = cli.IntFlag{
	Name:  "port, P",
	Usage: "remote port",
	Value: 22,
}

var userFlag = cli.StringFlag{
	Name:  "user, u",
	Usage: "login user",
}

var passwdFlag = cli.StringFlag{
	Name:  "passwd, p",
	Usage: "login passwd",
}

type remote struct {
	Host   string
	Port   int
	User   string
	Passwd string
}

func remoteFromContext(context *cli.Context) (r *remote, err error) {
	r, err = remoteFromFlag(context)
	args := context.Args()
	if len(args) == 0 {
		return
	}
}

func remoteFromFlag(context *cli.Context) (r *remote, err error) {
	host := context.String(hostFlag.Name)
	if host == "" || net.ParseIP(host) == nil {
		err = fmt.Errorf("remote host \"%s\" must be an effective IP", host)
		return
	}
	r = &remote{
		Host:   host,
		Port:   context.Int(portFlag.Name),
		User:   context.String(userFlag.Name),
		Passwd: context.String(passwdFlag.Name),
	}
	return
}

func remoteFromCmdLine(cmdLine string) (r *remote, err error) {
	/*
	100.101.190.114
	100.101.190.114:9999
	root@100.101.190.114
	root@100.101.190.114:9999
	root:Changeme_123@100.101.190.114
	root:Changeme_123@100.101.190.114:9999

	*/
	return
}
