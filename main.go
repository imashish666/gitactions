package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"www-api/config"
	"www-api/internal/logger"
	"www-api/internal/server"
)

// @title           WWW-API
// @version         1.0
// @description     This is www-api server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Securly
// @contact.url    https://https://www.securly.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  Bearer

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	configFile := flag.String("config", "config/config.yaml", "config file location")
	region := flag.String("region", "", "aws region where service is deployed")
	secret := flag.String("secret", "", "secret name where configs can be found")
	deployment := flag.String("deployment", "", "prod or dev")
	flag.Parse()
	logger, err := logger.NewZapLogger()
	if err != nil {
		fmt.Println("failed to initialize logger")
		log.Fatalf("unable to initiate logger %v", err)
	}

	conf, err := config.LoadConfig(*configFile, *region, *deployment, *secret, logger)
	if err != nil {
		fmt.Println("failed to fetch configs")
		log.Fatal("unable to fetch configs", map[string]interface{}{
			"error": err,
		})
	}
	r := server.GetRouter(conf, logger)
	server := &http.Server{Addr: ":" + conf.Server.Port, Handler: r}

	go func() {
		serverErr := server.ListenAndServe()
		log.Println(serverErr)
		logger.Fatal("error starting server", map[string]interface{}{"error": serverErr})
	}()

	log.Println("server started at port ", conf.Server.Port)
	logger.Info("server started", nil)
	//creating a channel to receive an os interruption signal
	stopC := make(chan os.Signal, 1)

	//notifying the channel when an interrupt signal is received
	signal.Notify(stopC, os.Interrupt)

	// blocking untill a interruption signal is received
	<-stopC

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("server stopping...")
	logger.Info("server stopping...", nil)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Println(err)
		logger.Error("error shutting down server", map[string]interface{}{"error": err})
	}
}
