package main

import (
	"bush/handler"
	"bush/repository"
	"bush/service"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func main() {
	//db := config.FirebaseDB()
	logrus.SetFormatter(new(logrus.JSONFormatter))
	repos, err := repository.NewRepository()
	if err != nil {
		print("error in repository")
	}
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(Server)

	started := make(chan struct{})

	go func() {
		if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
		close(started)
	}()

	logrus.Print("App Started")

	<-started // wait for the server to start

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
