package main

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/opentelemetry"
	"flamingo.me/pugtemplate"
	"saitho.me/styx-app/src/lib"

	"saitho.me/styx-app/src/api"
	"saitho.me/styx-app/src/module"
	"saitho.me/styx-app/src/rpc"
)

// main is our entry point
func main() {
	flamingo.App([]dingo.Module{
		new(zap.Module),
		//new(healthcheck.Module),
		new(opentelemetry.Module),
		// log formatter
		new(requestlogger.Module), // request logger show request logs
		new(pugtemplate.Module),
		new(lib.ConfigModule),
		new(module.Module),
		new(api.Module),
		new(rpc.Module),
	})
}
