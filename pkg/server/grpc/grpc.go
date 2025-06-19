package grpcapp

import (
	"github.com/Sanchir01/currency-wallet/pkg/logger"
	walletsv1 "github.com/Sanchir01/wallets-proto/gen/go/wallets"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type App struct {
	api walletsv1.ExchangeServiceClient
	log *slog.Logger
}

func NewClientGRPC[T any](
	l *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
	clientFactory func(grpc.ClientConnInterface) T,
) (T, error) {
	const op = "grpc.client.new"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithBackoff(grpcretry.BackoffLinear(500 * time.Millisecond)),
		grpcretry.WithCodes(codes.Aborted, codes.Unavailable, codes.NotFound, codes.DeadlineExceeded),
		grpcretry.WithPerRetryTimeout(timeout),
		grpcretry.WithMax(uint(retriesCount)),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
			grpclog.UnaryClientInterceptor(logger.InterceptorsLogger(l), logOpts...),
		),
	)
	if err != nil {
		var zero T
		return zero, err
	}

	client := clientFactory(cc)

	return client, nil
}
