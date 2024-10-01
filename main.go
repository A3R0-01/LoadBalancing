package main

import (
	"context"
	"time"

	serverpool "github.com/A3R0-01/loadbalance/serverPool"
)

func HealthCheck(ctx context.Context, pool serverpool.ServerPool) {
	aliveChanel := make(chan bool, 1)
	for _, backend := range pool.GetBackends() {
		backend := backend

		requestCtx, stop := context.WithTimeout(ctx, 10*time.Second)
		defer stop()
		go IsBackendAlive(requestCtx, aliveChanel, backend.GetUrl())

		select {
		case <-ctx.Done():
			utils.Logger.Info("Gracefully shutting down the health check")
		case alive := <-aliveChanel:
			backend.SetAlive(alive)
		}
	}
}
func main() {

}
