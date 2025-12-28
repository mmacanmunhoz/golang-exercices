package config

type ServerConfig struct {
	Name     string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Replicas int    `json:"replicas" yaml:"replicas"`
}

func (s ServerConfig) String() string {
	return s.Name + " at " + s.Host + ":" + string(rune(s.Port))
}

type DatabaseConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
}


type Config struct {
	Servers  []ServerConfig   `json:"servers" yaml:"servers"`
	Database DatabaseConfig   `json:"database" yaml:"database"`
}