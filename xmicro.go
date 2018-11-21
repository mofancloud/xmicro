package xmicro

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	rcache "github.com/micro/go-rcache"
)

var (
	// AppPath is the absolute path to the app
	AppPath string
)

func init() {
	var err error
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
}

func Run(serviceName string, fn func(c *cli.Context, service micro.Service)) {
	service := micro.NewService(
		micro.Name(serviceName),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
		micro.Registry(rcache.New(registry.DefaultRegistry)),
		micro.Flags(
			cli.StringFlag{
				Name:  "config",
				Usage: "config file",
			},
		),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			configFilePath := c.String("config")
			if len(configFilePath) == 0 {
				fmt.Printf("Usage: %s --config=path\n", AppPath)
				os.Exit(-1)
			}

			// 1. 判断文件是否存在
			_, err := os.Stat(configFilePath)
			if err != nil {
				if os.IsNotExist(err) {
					panic(err)
				}
			}

			// 2. 读取配置文件, 理论可支持不同策略
			LoadAppConfig("json", configFilePath)

			fn(c, service)
		}),
	)

	if err := service.Run(); err != nil {
		panic(err)
	}
}
