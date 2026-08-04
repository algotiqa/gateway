//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package service

import (
	"crypto/tls"
	"crypto/x509"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/algotiqa/core"
	"github.com/algotiqa/gateway/pkg/app"
	"github.com/gin-gonic/gin"
)

//=============================================================================

var gatewayCfg   *app.Config
var transportCfg *http.Transport

//=============================================================================

func Init(cfg *app.Config, router *gin.Engine, logger *slog.Logger) {
	gatewayCfg   = cfg
	transportCfg = createHttpTransport(logger)
	router.Use(handleUrl)
}

//=============================================================================

func createHttpTransport(logger *slog.Logger) *http.Transport {
	cert, err := os.ReadFile("config/ca.crt")
	if err != nil {
		core.ExitWithMessage("Could not open certificate file: " + err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	certificate, err := tls.LoadX509KeyPair("config/server.crt", "config/server.key")
	if err != nil {
		core.ExitWithMessage("Could not load certificate: " + err.Error())
	}

	return &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{certificate},
		},
	}
}

//=============================================================================

func handleUrl(c *gin.Context) {
	start := time.Now()
	slog.Info("New request", "client", c.ClientIP(), "context", c.Request.URL.String())
	path := c.Request.URL.Path

	targetURL := lookupTargetURL(path)
	if targetURL == "" {
		c.String(404, "Not Found")
		slog.Error("URL mapping not found", "client", c.ClientIP(), "path", path)
		return
	}

	proxyUrl(targetURL, c)
	duration := time.Since(start)
	slog.Info("Request served", "duration", duration.Seconds())
}

//=============================================================================

func lookupTargetURL(path string) string {
	prefix := ""
	target := ""

	for _, elem := range gatewayCfg.Routes {
		subPath, found := strings.CutPrefix(path, elem.Prefix)
		if found && len(prefix) < len(elem.Prefix) {
			if !strings.HasPrefix(subPath, "/") {
				subPath = "/" + subPath
			}

			prefix = elem.Prefix
			target = elem.Url + subPath
		}
	}

	return target
}

//=============================================================================

func proxyUrl(targetURL string, c *gin.Context) {
	target, err := url.Parse(targetURL)
	if err != nil {
		c.String(500, "Invalid target URL")
		slog.Error("Invalid target URL", "client", c.ClientIP(), "url", targetURL)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Transport = transportCfg
	proxy.Director  = func(request *http.Request) {
		request.URL.Scheme = target.Scheme
		request.URL.Host   = target.Host
		request.URL.Path   = target.Path
	}

	slog.Info("Forwarding request", "target", target)
	proxy.ServeHTTP(c.Writer, c.Request)
}

//=============================================================================
