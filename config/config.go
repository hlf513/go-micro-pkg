package config

import (
	"errors"
	"fmt"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/util/log"
)

var (
	// c 初始化配置器
	c = &configurator{}
)

// Configurator 配置器接口
type Configurator interface {
	// Init 初始化 micro config，并监控改动
	Init(configs []func(config []byte) error, source ...source.Source) error
	// Get 实时读取配置信息并解析
	Get(name []string, config interface{}) error
}

// configurator 配置器
type configurator struct {
	conf config.Config
}

// Init 初始化 micro config，并监控改动
func (c *configurator) Init(configs []func(config []byte) error, source ...source.Source) error {
	c.conf = config.NewConfig()
	// 加载配置
	if err := c.conf.Load(source...); err != nil {
		log.Fatal("[initialise config] load error:", err.Error())
	}
	// 监控改动
	go func() {
		w, err := c.conf.Watch()
		if err != nil {
			log.Fatal("[initialise config] watch error:", err.Error())
		}
		for {
			v, err := w.Next()
			if err != nil {
				log.Fatalf("[initialise config] watch next error，%s", err)
				return
			}
			// log.Infof("config was changed， %s", string(v.Bytes()))
			// 更新配置
			modify := v.Bytes()
			for _, f := range configs {
				if err := f(modify); err != nil {
					log.Warn("[update config] ", err.Error())
				}
			}
		}
	}()

	// 读取配置
	for _, f := range configs {
		if err := f(nil); err != nil {
			log.Fatal("[initialise config] ", err.Error())
		}
	}

	return nil
}

// Get 实时读取配置信息并解析
func (c *configurator) Get(name []string, config interface{}) error {
	value := c.conf.Get(name...)
	if value == nil {
		return errors.New(fmt.Sprintf("[initialise config] config name (%s):was not found", name))
	}
	return value.Scan(config)
}

// GetConfigurator 获取配置器
func GetConfigurator() Configurator {
	return c
}
