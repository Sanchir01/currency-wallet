package grpcapp

import (
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	Port       string
}

func New(lg *slog.Logger, port string, gRPCServer *grpc.Server) *App {
	return &App{
		log:        lg,
		gRPCServer: gRPCServer,
		Port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Start(); err != nil {
		panic(err)
	}
}

func GetGrpcServer() *grpc.Server {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(RecoveryInterceptor),
	)
	return srv
}
func (a *App) Start() error {
	const op = "grpcapp.App.Start"
	log := a.log.With(slog.String("op", op), "Starting gRPC server", "port", a.Port)
	l, err := net.Listen("tcp", a.Port)
	if err != nil {
		log.Error("Failed to start gRPC server", "error", err)
		return err
	}
	log.Info("gRPC server started", "address", l.Addr().String())
	if err := a.gRPCServer.Serve(l); err != nil {
		a.log.Error("Failed to start gRPC server", "error", err)
		return err
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.App.Stop"
	log := a.log.With(slog.String("op", op), "Stopping gRPC server")
	a.gRPCServer.GracefulStop()
	log.Info("gRPC server stopped")
}
