package main

import (
	"context"
	"fmt"
	"github.com/ShobenHou/monitor/internal/agent"
	"github.com/ShobenHou/monitor/internal/cp/api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//TODO:do as template do:
//https://github.com/gin-gonic/examples/blob/master/template/main.go

func main() {
	// init agent
	//TODO: STEP1 connect to default agent
	//TODO: STEP2 1.deploy by ansible 2.connect automatically to all agents
	//TODO: STEP3 stop deleted agent

	//STEP1. register the  default agent . Temp okay
	agentConf := &agent.AgentConf{
		Interval: "1s",
		Metrics: []string{
			"cpu",
			"mem",
			"host",
		},
		Addr: "localhost:55555",
	}
	agentInstance := agent.NewAgent(agentConf)

	//STEP2: init API server.
	//Temp okay. To update: api.NewServer() should receive a list of agents
	apiServer := api.NewServer(agentInstance)

	//start HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: apiServer.Router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	//-------------------------choice 1:  handle shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("shutting down server...")

	// shutdown gracefully
	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("server forced to shutdown: %s", err)
	}

	//agentInstance.Stop()

	/*-------------choice 2: gin official graceful shut-down: https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go

		// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	*/

	fmt.Println("server exiting")

	/* --------------old code-------------
	router := gin.Default()
	api.SetupRoutes(router)
	router.Run()
	*/
}
