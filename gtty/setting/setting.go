package setting

import (
	"io/ioutil"
	"os"

	"github.com/chinx/gtty/model"
	"github.com/go-yaml/yaml"
)

var Root = model.RootGroup{}

func init()  {
	OpenConfig("./config/config.yaml")
}

func OpenConfig(fileName string) error {
	if !IsFileExist(fileName) {
		return ioutil.WriteFile(fileName, []byte(""), 0750)
	}
	Root = model.RootGroup{}

	byteArr, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(byteArr, &Root); err != nil {
		return err
	}
	return nil
}

func IsFileExist(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir() || os.IsExist(err)
}
