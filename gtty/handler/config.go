package handler

import (
	"fmt"

	"github.com/chinx/gtty/model"
)

func ListRemotes(remotes []*model.Group, group string) {
	needMatch := group != ""
	for key, val := range remotes {
		equal := val.Name == group
		if !needMatch || equal {
			fmt.Printf("%s:\n", val.Name)
			for k, v := range val.List {
				if v.Name != "" {
					fmt.Printf("    %d-%d | %s-%s | %s\n", key, k, val.Name, v.Name, v.Host)
				} else {
					fmt.Printf("    %d-%d | %s\n", key, k, v.Host)
				}
			}
		}
		if equal {
			break
		}
	}
}

func findByIndex(remotes []*model.Group, keys []int) (config *model.RemoteConfig) {
	if len(keys) < 2 {
		return
	}

	if i := keys[0]; len(remotes) > i {
		list := remotes[i].List
		if j := keys[1]; len(list) > j {
			config = list[j]
		}
	}
	return
}

func findByName(remotes []*model.Group, keys []string) (config *model.RemoteConfig) {
	for _, val := range remotes {
		if val.Name == keys[0] {
			for _, v := range val.List {
				if v.Name == keys[1] {
					return v
				}

			}
		}
	}
	return
}

func findByHost(remotes []*model.Group, host string) (config *model.RemoteConfig) {
	for _, val := range remotes {
		for _, v := range val.List {
			if v.Host == host {
				return v
			}
		}
	}
	return
}
