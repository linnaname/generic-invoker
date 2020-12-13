package config

/**
implement this method so you can have your owner config loader
*/
type Loader interface {
	Load() (Configs, error)
}

const (
	FILE  = "file"
	NACOS = "nacos"
)

/**
default file loader
TODO extend other kind config loader
*/
func NewLoader(kind string) Loader {
	return &FileLoader{
		name: "config.json",
		path: ".",
	}

}
