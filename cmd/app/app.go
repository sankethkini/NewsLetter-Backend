package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/pkg/errors"
	"github.com/sankethkini/NewsLetter-Backend/internal/config"
	"github.com/sankethkini/NewsLetter-Backend/internal/service"
	transport "github.com/sankethkini/NewsLetter-Backend/internal/transport"
	"github.com/sankethkini/NewsLetter-Backend/pkg/auth"
	"github.com/sankethkini/NewsLetter-Backend/pkg/log"
	adminpb "github.com/sankethkini/NewsLetter-Backend/proto/adminpb/v1"
	newsletterpb "github.com/sankethkini/NewsLetter-Backend/proto/newsletterpb/v1"
	subscriptionpb "github.com/sankethkini/NewsLetter-Backend/proto/subscriptionpb/v1"
	userpb "github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Start(ctx context.Context) {
	logger := log.Build()
	cfg, err := IntializeServerConfig()
	if err != nil {
		logger.Sugar().Fatal("error:", err)
	}

	registry, clean, err := IntializeServiceRegistry()
	if err != nil {
		logger.Sugar().Fatal("error:", err)
	}

	authInterceptor, err := IntializeJWT()
	if err != nil {
		logger.Sugar().Fatal("error:", err)
	}

	defer func() {
		e := logger.Sync()
		if err != nil {
			logger.Sugar().Fatal("error:", e)
		}
	}()

	newCtx := ctxzap.ToContext(ctx, logger)

	if err := run(newCtx, authInterceptor, clean, registry, cfg); err != nil {
		logger.Sugar().Fatal("error: ", err)
	}
}

func run(ctx context.Context, auth *auth.AuthInterceptor, clean func(), registry *service.Registry, cfg config.ServerConfig) error {
	var err error
	logger := ctxzap.Extract(ctx)

	defer func() {
		logger.Sugar().Info("invoking database cleanup function")
		clean()
		logger.Sugar().Info("database closed")
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	var listener net.Listener
	var server *grpc.Server
	serverErrors := make(chan error, 1)

	// kafka consumer
	c, err := IntializeConsumer()
	if err != nil {
		return err
	}

	go c.Consume(ctx)

	go func() {
		addr := cfg.Host + ":" + cfg.Port
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			logger.Sugar().Fatal("server not started: ", err)
		}

		grpcUsr := transport.NewUserGrpcServer(ctx, registry.UserService)
		grpcSub := transport.NewSubscriptionService(ctx, registry.SubscriptionService)
		grpcNews := transport.NewNewsGrpcServer(ctx, registry.NewsService)
		grpcAdm := transport.NewAdminGrpcServer(ctx, registry.AdminService)

		server = grpc.NewServer(
			grpc.UnaryInterceptor(auth.Unary()),
		)
		userpb.RegisterUserServiceServer(server, grpcUsr)
		subscriptionpb.RegisterSubscriptionServiceServer(server, grpcSub)
		adminpb.RegisterAdminServiceServer(server, grpcAdm)
		newsletterpb.RegisterNewsLetterServiceServer(server, grpcNews)
		reflection.Register(server)
		serverErrors <- server.Serve(listener)
	}()

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")
	case sig := <-shutdown:
		logger.Sugar().Info("start shutdown ", sig)
		logger.Info("stopping grpc server")
		server.GracefulStop()
		logger.Info("closing grpc listener")
		_ = listener.Close()
		// ignored error
		logger.Info("grpc server gracefully shutdown")
	}
	return nil
}
