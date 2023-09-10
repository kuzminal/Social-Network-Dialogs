package main

import (
	"Social-Net-Dialogs/internal/counters"
	"Social-Net-Dialogs/internal/handler"
	"Social-Net-Dialogs/internal/helper"
	"Social-Net-Dialogs/internal/metrics"
	"Social-Net-Dialogs/internal/router"
	"Social-Net-Dialogs/internal/service"
	"Social-Net-Dialogs/internal/session"
	"Social-Net-Dialogs/internal/store/pg"
	"Social-Net-Dialogs/internal/store/tarantool"
	"Social-Net-Dialogs/internal/tracing"
	"Social-Net-Dialogs/models"
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
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

	//metrics
	prometheus.MustRegister(metrics.SendMessage, metrics.GetMessage, metrics.MarkAsRead)

	traceServer := helper.GetEnvValue("TRACE_SERVER", "trace")
	tracePort := helper.GetEnvValue("TRACE_PORT", "14268")
	tracer, err := tracing.TracerProvider(fmt.Sprintf("http://%s:%s/api/traces", traceServer, tracePort))
	if err != nil {
		log.Fatal(err)
	}
	defer tracer.Shutdown(context.Background())

	connectToWsChan := make(chan models.ActiveWsUsers, 10)
	disconnectToWsChan := make(chan models.ActiveWsUsers, 10)
	tokenService := service.NewTokenServiceClient(tarantoolMaster, tracer)
	counterPublisher := counters.NewCountersPublisher()
	app := handler.NewInstance(
		tarantoolMaster,
		master,
		tokenService,
		connectToWsChan,
		disconnectToWsChan,
		tracer,
		counterPublisher,
	)
	sessionConsumer := session.NewSessionConsumer(tarantoolMaster)
	go sessionConsumer.ReadSessionInfo(context.Background())

	r := router.NewRouter(app)
	appPort := helper.GetEnvValue("PORT", "8081")
	log.Fatalln(http.ListenAndServe(":"+appPort, r))
}

func initDb() {
	pghost := helper.GetEnvValue("PGHOST", "localhost")
	pgport := helper.GetEnvValue("PGPORT", "5432")
	master, _ = pg.NewMaster(context.Background(), fmt.Sprintf("postgresql://postgres:postgres@%s:%s/postgres?sslmode=disable", pghost, pgport))
}

func initTarantoolDb() {
	thost := helper.GetEnvValue("TARANTOOL_HOST", "localhost")
	tport := helper.GetEnvValue("TARANTOOL_PORT", "3301")
	tuser := helper.GetEnvValue("TARANTOOL_USER_NAME", "user")
	tpassword := helper.GetEnvValue("TARANTOOL_USER_PASSWORD", "password")
	tarantoolMaster, _ = tarantool.NewTarantoolMaster(thost, tport, tuser, tpassword)
}
