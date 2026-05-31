package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	moduleLib "margin-delver/lib"
	"margin-delver/middleware"
	v1 "margin-delver/router/v1"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewServer(
	lc fx.Lifecycle,
	cfg *moduleLib.AppConfig,
	log *moduleLib.BaseLog,
	db *gorm.DB,
) *http.Server {

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}

	router := gin.New()

	router.Use(
		middleware.CORS(),
		middleware.RequestLogger(log),
		middleware.Recovery(log),
	)

	v1.SetRouter(router, cfg, log, db)

	serverTimeout := time.Duration(cfg.ServerTimeOut)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AppPort),
		Handler:      router,
		ReadTimeout:  serverTimeout * time.Second,
		WriteTimeout: serverTimeout * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: startServer(srv, log),
		OnStop:  stopServer(srv, log),
	})

	return srv
}

func startServer(
	srv *http.Server,
	log *moduleLib.BaseLog,
) func(context.Context) error {

	go func() {
		_ = http.ListenAndServe(":6061", nil)
	}()

	return func(ctx context.Context) error {

		ln, err := net.Listen("tcp", srv.Addr)
		if err != nil {
			return err
		}

		log.SugarLog().Infof(
			"Starting HTTP server at %s",
			srv.Addr,
		)

		go func() {
			if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
				log.SugarLog().Errorf(
					"server stopped: %v",
					err,
				)
			}
		}()

		return nil
	}
}

func stopServer(
	srv *http.Server,
	log *moduleLib.BaseLog,
) func(context.Context) error {

	return func(ctx context.Context) error {

		log.SugarLog().Infof(
			"Stopping HTTP server at %s",
			srv.Addr,
		)

		return srv.Shutdown(ctx)
	}
}
