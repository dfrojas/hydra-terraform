package config

type LocalstackConfig struct {
	Source           string
	Output           string
	KeepResources    []string
	RemoveAttributes map[string][]string
}
