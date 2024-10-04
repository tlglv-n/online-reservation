package app

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"reservation-system/internal/config"
	"reservation-system/internal/handler"
	"reservation-system/internal/repository"
	"reservation-system/internal/service/reservation"
	"reservation-system/pkg/log"
	"reservation-system/pkg/server"
	"syscall"
	"time"
)

func Run() {
	logger := log.LoggerFromContext(context.Background())

	configs, err := config.New()
	if err != nil {
		logger.Error("ERR_INIT_CONFIGS", zap.Error(err))
		return
	}

	repositories, err := repository.New(
		repository.WithMemoryStore())
	if err != nil {
		logger.Error("ERR_INIT_REPOSITORIES", zap.Error(err))
		return
	}
	defer repositories.Close()

	reservationService, err := reservation.New(
		reservation.WithCandidateRepository(repositories.Candidate),
		reservation.WithRecruiterRepository(repositories.Recruiter))
	if err != nil {
		logger.Error("ERR_INIT_LIBRARY_SERVICE", zap.Error(err))
		return
	}

	handlers, err := handler.New(
		handler.Dependencies{
			Configs:            configs,
			ReservationService: reservationService,
		},
		handler.WithHTTPHandler())
	if err != nil {
		logger.Error("ERR_INIT_HANDLERS", zap.Error(err))
		return
	}

	servers, err := server.New(
		server.WithHTTPServer(handlers.HTTP, configs.APP.Port))
	if err != nil {
		logger.Error("ERR_INIT_SERVERS", zap.Error(err))
		return
	}
	if err = servers.Run(logger); err != nil {
		logger.Error("ERR_RUN_SERVERS", zap.Error(err))
		return
	}
	logger.Info("http server started on http://localhost:" + configs.APP.Port + "/swagger/index.html")

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1) // Create channel to signify a signal being sent

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
	<-quit                                             // This blocks the main thread until an interrupt is received
	fmt.Println("gracefully shutting down...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	if err = servers.Stop(ctx); err != nil {
		panic(err) // failure/timeout shutting down the httpServer gracefully
	}

	fmt.Println("running cleanup tasks...")
	// Your cleanup tasks go here

	fmt.Println("server was successful shutdown.")
}
