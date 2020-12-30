package config

import (
	"errors"
	"reflect"
)

type ConfigContainer struct {
	container map[string]*Congfig
}

func (this *ConfigContainer) NewConfig(name string, value interface{}) {
	reflect.TypeOf(value)
}

func (this *ConfigContainer) get(name string) (*Congfig, error) {
	if value, ok := this.container[name]; ok {
		return value, nil
	}
	return nil, errors.New("key not fund")
}

func (this *ConfigContainer) set(name string, value interface{}) error {
	this.container[name] = value
}

type Congfig struct {
	name  string
	path  string
	value interface{}
}

