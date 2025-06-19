package v1

import (
	"context"
	userv1 "github/smile-ko/go-template/api/proto/user/v1/gen"
	"github/smile-ko/go-template/pkg/postgres"
)

type Handler struct {
	userv1.UnimplementedUserServiceServer
	pg *postgres.Postgres
}

func NewHandler(pg *postgres.Postgres) *Handler {
	return &Handler{
		pg: pg,
	}
}

func (h *Handler) GetUserById(ctx context.Context, req *userv1.GetUserByIdReq) (*userv1.PublicUserInfoResp, error) {
	return nil, nil
}
