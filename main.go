//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package main

import (
	"crypto/tls"
	"log/slog"
	"net/http"

	"github.com/algotiqa/core"
	"github.com/algotiqa/core/boot"
	"github.com/algotiqa/gateway/pkg/app"
	"github.com/algotiqa/gateway/pkg/service"
	"github.com/gin-gonic/gin"
)

//=============================================================================

const component = "gateway"
var   version   = "dev"

//=============================================================================

func main() {
	cfg := &app.Config{}
	boot.ReadConfig(component, cfg)
	logger := boot.InitLogger(component, version, &cfg.Application)
	engine := boot.InitEngine(logger, &cfg.Application)
	service.Init(cfg, engine, logger)
	runHttpServer(engine, &cfg.Application)
}

//=============================================================================

func runHttpServer(router *gin.Engine, app *core.Application) {
	slog.Info("Starting HTTPS server...")

	tlsConfig := &tls.Config{}

	server := &http.Server{
		Addr     : app.BindAddress,
		TLSConfig: tlsConfig,
		Handler  : router,
	}

	slog.Info("Running")
	err := server.ListenAndServeTLS("config/server.crt", "config/server.key")
	core.ExitIfError(err)
}

//=============================================================================
