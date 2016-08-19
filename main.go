package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	serverURL = "http://%v:%v/vessel/v1/pipeline"
)

type requestType string
type CallArgs struct {
	Path   string
	Buffer *bytes.Buffer
	Kind   string
}

func (m *CallArgs) Exec() {
	newBuffer := *m.Buffer
	req, _ := http.NewRequest(m.Kind, m.Path, &newBuffer)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	} else {
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(data))
		log.Println(resp.StatusCode)
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	args := formatOsArgs(os.Args)
	if args == nil {
		log.Println("请设置运行所需参数 requestType jsonfile 【host】【port】")
	} else {
		args.Exec()
	}
}

func formatOsArgs(args []string) *CallArgs {
	argsLen := len(args)
	if argsLen <= 2 {
		log.Println("The jsonfile is not specified!")
		return nil
	}
	kind := strings.ToUpper(args[1])
	switch kind {
	case "POST":
	case "DELETE":
	case "DEL":
		kind = "DELETE"
	case "PUT":
	case "PATCH":
	case "GET":
	default:
		kind = ""
	}
	if kind == "" {
		log.Println("The requestType is wrong!")
		return nil
	}
	buffer, err := getBodyBuffer("./" + args[2])
	if err != nil {
		log.Println("The jsonfile has an error :", err)
		return nil
	}
	host := "127.0.0.1"
	port := "4488"
	if argsLen > 3 {
		host = args[3]
	}
	if argsLen > 4 {
		port = args[4]
	}
	return &CallArgs{Buffer: buffer, Kind: kind, Path: fmt.Sprintf(serverURL, host, port)}
}

func getBodyBuffer(filePath string) (*bytes.Buffer, error) {
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBytes), nil
}
