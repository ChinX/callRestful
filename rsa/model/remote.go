package model

type RemoteConfig struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Index int    `yaml:"index"`
	Name  string `yaml:"name,omitempty"`
}

type Group struct {
	Name string          `yaml:"name,omitempty"`
	List []*RemoteConfig `yaml:"list,omitempty"`
}

type RootGroup struct {
	Remotes     []*Group      `yaml:"remotes,omitempty"`
	Credentials []*Credential `yaml:"credentials,omitempty"`
}

type Credential struct {
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
}

type Remote struct {
	Host   string
	Port   int
	User   string
	Passwd string
}
