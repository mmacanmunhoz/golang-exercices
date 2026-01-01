package config

type DeployConfig struct {
	Name     string            `json:"name" yaml:"name"`
	Image    string            `json:"image" yaml:"image"`
	Port     map[string]string `json:"ports" yaml:"ports"`
	Env      []string          `json:"env" yaml:"env"`
	Volume   []string          `json:"volumes" yaml:"volumes"`
	Restart  string            `json:"restart" yaml:"restart"`
	Networks []string          `json:"networks" yaml:"networks"`
	Command  []string          `json:"command,omitempty" yaml:"command,omitempty"`
}

func (s DeployConfig) String() string {
	return s.Name + " at " + s.Image + ":"
}

type Config struct {
	Deploy []DeployConfig `json:"servers" yaml:"servers"`
}
