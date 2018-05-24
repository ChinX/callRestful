package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"strconv"
	"errors"
	"os/exec"
	"io/ioutil"
)

var conf = &Config{}

type Config struct {
	User string
	Pass string
	Host string
	Port uint
}

func (c *Config) ParseParams(param string) error {
	if len(param) == 0{
		if c.Host != ""{
			return nil
		}
		return errors.New("host muse not empty")
	}
	userIP := strings.Split(param,"@")
	count := len(userIP)
	if count > 1{
		userPass := strings.Split(userIP[0],":")
		if len(userPass) > 1 && c.Pass == "" {
			c.Pass = userPass[1]
		}
		if c.User == "" || c.User == "root"{
			c.User = userPass[0]
		}
	}
	hostPort := strings.Split(userIP[count-1],":")
	if len(hostPort) > 1{
		i, err := strconv.Atoi(hostPort[1])
		if err == nil && (i != 0 || i != 22) && c.Port == 22 {
			c.Port = uint(i)
		}
	}
	if c.Host == ""{
		c.Host = hostPort[0]
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name=filepath.Base(os.Args[0])
	app.Usage="golang ssh tools"
	app.Commands=[]cli.Command{
		{
			Name:"open",
			ShortName:"o",
			Usage:"open new ssh",
			Description:"params like:[user[:passwd]]host[:port]",
			Flags: getFlags("ssh", conf),
			Action: openSSH,
		},{
			Name:"sftp",
			ShortName:"c",
			Usage:"open new sftp",
			Description:"params like:[user[:passwd]]host[:port]",
			Flags:  getFlags("sftp", conf),
			Action: openWSCP,
		},{
			Name:"list",
			ShortName:"l",
			Usage:"show ssh list",
			Action: listSSH,
		},
	}
	app.Run(os.Args)
}

func getFlags(cmdType string, conf *Config) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name: "user,u",
			Value: "root",
			Usage: fmt.Sprintf("%s user", cmdType),
			Destination: &conf.Host,
		},
		cli.StringFlag{
			Name: "pass,p",
			Usage: fmt.Sprintf("%s pass word", cmdType),
			Destination: &conf.Pass,
		},
		cli.StringFlag{
			Name: "host,H",
			Usage: fmt.Sprintf("%s host", cmdType),
			Destination: &conf.Host,
		},
		cli.UintFlag{
			Name:"port,P",
			Value: 22,
			Usage: fmt.Sprintf("%s port", cmdType),
			Destination: &conf.Port,
		},

	}
}

func openSSH(ctx *cli.Context){
	args := ctx.Args()
	if len(args) > 0{
		conf.ParseParams(args[0])
	}
	if conf.Host == ""{
		cli.ShowCommandHelp(ctx, ctx.Command.Name)
		return
	}
	log.Println(*conf)
	cmd := exec.Command("ssh", conf.Host+":"+strconv.Itoa(int(conf.Port)), "-u", conf.User, "-p", conf.Pass)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(opBytes))
}

func openWSCP(ctx *cli.Context){
	args := ctx.Args()
	if len(args) > 0{
		conf.ParseParams(args[0])
	}
	if conf.Host == ""{
		cli.ShowCommandHelp(ctx, ctx.Command.Name)
		return
	}
	log.Println(*conf)
	cmd := exec.Command("sftp", conf.Host+":"+strconv.Itoa(int(conf.Port)), "-u", conf.User, "-p", conf.Pass)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(opBytes))
}

func listSSH(ctx *cli.Context){
	log.Println(ctx.Args())
}

