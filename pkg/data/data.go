package data

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Query struct {
	Sql string `yaml:"sql"`
	Type string `yaml:"type"`
	Connection string `yaml:"connection"`
	MessageBefore string `yaml:"message_before"`
	MessageAfter string `yaml:"message_after"`
	Async bool `yaml:"async"`
	Sleep int `yaml:"sleep"`
	Timeout int `yaml:"timeout"`
}

type Connection struct {
	Name string `yaml:"name"`
	Dsn string `yaml:"dsn"`
}

type Flow struct {
	Connections []Connection `yaml:"connections"`
	Queries []Query `yaml:"queries"`
}

type Loader interface {

}

type loader struct {

}

func NewLoader() *loader {
	return &loader{}
}

func (l *loader) Load(fileName string) (Flow, error) {
	dataBytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		return Flow{}, err
	}

	var flow Flow
	err = yaml.Unmarshal(dataBytes, &flow)

	if err != nil {
		return Flow{}, err
	}

	return flow, nil
}
