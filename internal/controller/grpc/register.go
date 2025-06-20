package grpc

import (
	v1 "github/smile-ko/go-template/internal/controller/grpc/v1"
	"github/smile-ko/go-template/pkg/logger"
	"github/smile-ko/go-template/pkg/postgres"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RegisterGRPCServices(app *grpc.Server, pg *postgres.Postgres, l logger.ILogger) {
	{
		v1.RegisterV1GRPC(app, pg, l)
	}

	reflection.Register(app)
}
