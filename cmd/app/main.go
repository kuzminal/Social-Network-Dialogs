package main

import (
	"Social-Net-Dialogs/internal/handler"
	"Social-Net-Dialogs/internal/helper"
	"Social-Net-Dialogs/internal/router"
	"Social-Net-Dialogs/internal/service"
	"Social-Net-Dialogs/internal/store/pg"
	"Social-Net-Dialogs/internal/store/tarantool"
	"Social-Net-Dialogs/models"
	"context"
	"fmt"
	"log"
	"net/http"
)

var (
	master          *pg.Postgres
	tarantoolMaster *tarantool.Tarantool
)

func main() {
	initDb()
	initTarantoolDb()
	connectToWsChan := make(chan models.ActiveWsUsers, 10)
	disconnectToWsChan := make(chan models.ActiveWsUsers, 10)
	tokenService := service.NewTokenServiceClient(tarantoolMaster)
	app := handler.NewInstance(
		tarantoolMaster,
		master,
		tokenService,
		connectToWsChan,
		disconnectToWsChan,
	)

	r := router.NewRouter(app)
	log.Fatalln(http.ListenAndServe(":8081", r))
}

func initDb() {
	pghost := helper.GetEnvValue("PGHOST", "localhost")
	pgport := helper.GetEnvValue("PGPORT", "5432")
	master, _ = pg.NewMaster(context.Background(), fmt.Sprintf("postgresql://postgres:postgres@%s:%s/postgres?sslmode=disable", pghost, pgport))
}

func initTarantoolDb() {
	thost := helper.GetEnvValue("TARANTOOL_HOST", "localhost")
	tport := "3301" //пока так
	tuser := helper.GetEnvValue("TARANTOOL_USER_NAME", "user")
	tpassword := helper.GetEnvValue("TARANTOOL_USER_PASSWORD", "password")
	tarantoolMaster, _ = tarantool.NewTarantoolMaster(thost, tport, tuser, tpassword)
}
