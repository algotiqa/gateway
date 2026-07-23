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
	"github.com/algotiqa/core/boot"
	"github.com/algotiqa/gateway/pkg/app"
	"github.com/algotiqa/gateway/pkg/service"
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
	boot.RunHttpServer(engine, &cfg.Application)
}

//=============================================================================
