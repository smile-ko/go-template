package v1

import (
	userv1 "github/smile-ko/go-template/docs/proto/user/v1/gen"
	"github/smile-ko/go-template/pkg/logger"
	"github/smile-ko/go-template/pkg/postgres"

	"google.golang.org/grpc"
)

func RegisterV1GRPC(app *grpc.Server, pg *postgres.Postgres, l logger.Interface) {
	h := NewHandler(pg)

	userv1.RegisterUserServiceServer(app, h)
}
