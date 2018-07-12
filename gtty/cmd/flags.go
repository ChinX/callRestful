package cmd

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/chinx/gtty/model"
	"github.com/chinx/gtty/setting"
	"github.com/urfave/cli"
)

var rgx = regexp.MustCompile(`^(\d+)\-(\d+)$`)

var remoteDirFlag = cli.StringFlag{
	Name:  "remote, r",
	Usage: "remote open dir",
}

var currentDirFlag = cli.StringFlag{
	Name:  "current, c",
	Usage: "current open dir",
}

func remoteFromContext(context *cli.Context) (r *model.Remote, err error) {
	args := context.Args()
	if len(args) == 0 {
		return
	}

	if strings.Index(args[0], ".") == -1 {
		r, err = remoteFromIndex(args[0])
		if err == nil {
			return
		}
	}

	r, err = remoteFromCmdLine(args[0])
	if err != nil {
		return
	}

	nr, err := remoteFromHost(r.Host)
	if err == nil {
		if r.Port <= 0 {
			r.Port = nr.Port
		}

		if r.User == "" {
			r.User = nr.User
		}
	}

	if r.Port <= 0 {
		r.Port = 22
	}

	if len(args) > 1 {
		r.Passwd = args[1]
	}
	return
}

func remoteFromCmdLine(cmdLine string) (r *model.Remote, err error) {
	r = &model.Remote{
		Host: cmdLine,
	}
	userHost := strings.Split(cmdLine, "@")
	if len(userHost) > 1 {
		r.User = userHost[0]
		r.Host = userHost[1]
	}

	hostPort := strings.Split(r.Host, ":")
	if len(hostPort) > 1 {
		r.Host = hostPort[0]
		r.Port, _ = strconv.Atoi(hostPort[1])
	}

	if r.Host == "" || net.ParseIP(r.Host) == nil {
		err = fmt.Errorf("remote command line \"%s\" must has effective IP", cmdLine)
		return
	}
	return
}

func remoteFromIndex(index string) (r *model.Remote, err error) {
	if strings.Index(index, "-") < 0 {
		err = fmt.Errorf("remote command line \"%s\" is wrong", index)
		return
	}
	matches := rgx.FindStringSubmatch(index)
	if len(matches) == 3 {
		remotes := setting.Root.Remotes
		i, _ := strconv.Atoi(matches[1])
		if len(remotes) < i+1 {
			err = fmt.Errorf("remote command line index \"%s\" is wrong", index)
			return
		}
		list := setting.Root.Remotes[i].List
		j, _ := strconv.Atoi(matches[2])
		if len(list) < j+1 {
			err = fmt.Errorf("remote command line index \"%s\" is wrong", index)
			return
		}
		credential := setting.Root.Credentials[list[j].Index]
		r = &model.Remote{
			Host:   list[j].Host,
			Port:   list[j].Port,
			User:   credential.User,
			Passwd: credential.Passwd,
		}
		return
	}

	names := strings.Split(index, "-")
	for _, val := range setting.Root.Remotes {
		if val.Name == names[0] {
			for _, v := range val.List {
				if v.Name == names[1] {
					credential := setting.Root.Credentials[v.Index]
					r = &model.Remote{
						Host:   v.Host,
						Port:   v.Port,
						User:   credential.User,
						Passwd: credential.Passwd,
					}
					return
				}

			}
			break
		}
	}
	err = fmt.Errorf("remote command line index \"%s\" is not fount", index)
	return
}

func remoteFromHost(host string) (r *model.Remote, err error) {
	for _, val := range setting.Root.Remotes {
		for _, v := range val.List {
			if v.Host == host {
				credential := setting.Root.Credentials[v.Index]
				r = &model.Remote{
					Host:   v.Host,
					Port:   v.Port,
					User:   credential.User,
					Passwd: credential.Passwd,
				}
				return
			}

		}
	}
	err = fmt.Errorf("remote host \"%s\" is not fount", host)
	return
}
