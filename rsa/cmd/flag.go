package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/chinx/rsa/model"
	"github.com/chinx/rsa/setting"
	"github.com/urfave/cli"
)

var rgx = regexp.MustCompile(`^(\d+)\-(\d+)$`)

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

	nr, err1 := remoteFromHost(r.Host)
	if err1 == nil {
		r = nr
		return
	}

	if r.Port <= 0 {
		r.Port = 22
	}

	if r.User == "" {
		r.User = "root"
	}

	if r.Passwd == "" && len(args) > 1 {
		r.Passwd = args[1]
	}

	if r.Passwd == "" {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Printf("Passwd is empty! Please enter new: ")
		input, err3 := inputReader.ReadString('\n')
		if err3 != nil || input == "" {
			err = fmt.Errorf("Input passwd \"%s\"is error: %s ", input, err3)
			return
		}
		r.Passwd = input
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
