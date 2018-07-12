package handler

import (
	"fmt"
	"os"

	"github.com/chinx/gtty/model"
	"github.com/chinx/gtty/service"
)

func OpenSSH(config *model.Remote) {
	checkRemoteConfig(config)
	svc := &service.Server{
		User: config.User,
		Password: config.Passwd,
		Host: config.Host,
		Port: config.Port,
	}
	svc.Connection()
	fmt.Println("over")
}

func OpenSFTP(config *model.Remote) {
	checkRemoteConfig(config)
}

func checkRemoteConfig(config *model.Remote) {
	if config.User == "" {
		config.User = flagFromInput("user")
	}
	if config.Passwd == "" {
		config.Passwd = flagFromInput("passwd")
	}
}

func flagFromInput(name string) (value string) {
	fmt.Printf("%s is empty! Please enter new: ", name)
	fmt.Scanln(&value)
	if value == "" {
		fmt.Printf("input %s value is empty", name)
		os.Exit(1)
		return
	}
	return
}
