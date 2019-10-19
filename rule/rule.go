package rule

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/vipsimage/vipsimage/operation"
)

// Version vipsimage
const Version = "1.0.0"

var (
	short       = make(map[string]operation.Operation)
	disableSave = true
	mu          sync.RWMutex
)

// Init vipsimage operation rule.
// load vipsimage.toml, check config version, if expire upgrade config.
func Init() {
	version := viper.GetString("version")
	// upgrade config
	if version != Version {
		viper.Set("version", Version)
		err := viper.WriteConfig()
		if err != nil {
			logrus.Errorln(err.Error())
		}
	}

	// parse operation rule
	ruleConf := viper.GetStringMapString("operation-rule")
	for alias, operationRule := range ruleConf {
		err := Add(alias, operationRule)
		if err != nil {
			panic(err)
		}
	}

	disableSave = false
}

// GetAll return all parsed operation
func GetAll() map[string]operation.Operation {
	return short
}

// Get operation by alias
func Get(alias string) (op operation.Operation, ok bool) {
	mu.RLock()
	op, ok = short[alias]
	mu.RUnlock()

	return
}

// Add operation rule
func Add(alias, operationRule string) (err error) {
	if url.PathEscape(alias) != alias {
		return errors.New("alias cant path escape")
	}

	_, ok := Get(alias)
	if ok {
		return fmt.Errorf("alias: %s, already existed", alias)
	}

	return Set(alias, operationRule)
}

// Set alias operation rule to operationRule
func Set(alias, operationRule string) (err error) {
	opRule, err := operation.Parse(operationRule)
	if err != nil {
		return
	}

	add2conf(fmt.Sprintf("operation-rule.%s", alias), operationRule)

	mu.Lock()
	short[alias] = opRule
	mu.Unlock()

	return
}

// add2conf add config to vipsimage.toml
func add2conf(k, v string) {
	if disableSave {
		return
	}

	viper.Set(k, v)
	err := viper.WriteConfig()
	if err != nil {
		logrus.Errorln(err.Error())
	}
}

// Del delete alias
func Del(alias string) {
	mu.Lock()
	delete(short, alias)
	mu.Unlock()

	err := unset(fmt.Sprintf(`operation-rule.%s`, alias))
	if err != nil {
		logrus.Errorln(err.Error())
	}
}

// unset vipsimage.toml key, and write
func unset(key string) error {
	configMap := viper.AllSettings()
	delete(configMap, key)

	encodedConfig, err := json.MarshalIndent(configMap, "", " ")
	if err != nil {
		return err
	}

	err = viper.ReadConfig(bytes.NewReader(encodedConfig))
	if err != nil {
		return err
	}

	return viper.WriteConfig()
}
