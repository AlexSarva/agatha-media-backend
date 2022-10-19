package dbutils

import "strings"

type ConfigDB struct {
	User         string
	Password     string
	DatabaseName string
}

func ParseConfigDB(text string) (*ConfigDB, error) {
	paths := strings.Split(text, " ")
	dict := make(map[string]string)
	for _, v := range paths {
		info := strings.Split(v, "=")
		key := info[0]
		value := info[1]
		dict[key] = value
	}
	return &ConfigDB{
		User:         dict["user"],
		Password:     dict["password"],
		DatabaseName: dict["dbname"],
	}, nil
}
