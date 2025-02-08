//go:build tools
// +build tools

//

//go:generate oapi-codegen --config=../../api/config.yaml -o ../../api/v1/go/api.gen.go ../../api/v1/openapi.yaml
package main

import (
	"fmt"
	"log"

	"github.com/karthik2304/xm-go-service/configs"
	"github.com/karthik2304/xm-go-service/internal/repository"
	"github.com/karthik2304/xm-go-service/internal/server"
	"github.com/karthik2304/xm-go-service/pkg"
)

var (
	Version = "dev"
	service *server.Server
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if err := bootstrap(); err != nil {
		log.Println(err)
		return
	}

	if Version == "dev" {
		if err := service.Server.Start(fmt.Sprintf(":%d", configs.Settings.APP_PORT)); err != nil {
			log.Println(err)
		}
		return
	}
}

func bootstrap() error {

	configs.LoadConfig() // load environment configs

	kr, kw := pkg.ConnectKafka()
	db, err := pkg.ConnectMongo()
	if err != nil {
		panic(err)
	}

	r := pkg.New() // route configs

	service = server.New(r, repository.NewDB(db), kw, kr)

	return nil

}
