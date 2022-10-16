package service

import (
	"context"
	"fmt"
	"github.com/transerver/accounts/internal/biz"
	"github.com/transerver/commons/errors"
	"github.com/transerver/protos/acctspb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

type PubService struct {
	acctspb.UnimplementedRsaServiceServer

	usecase *biz.PubUsecase
	logger  *zap.Logger
}

func NewRsaService(usecase *biz.PubUsecase, logger *zap.Logger) *PubService {
	return &PubService{usecase: usecase, logger: logger}
}

func (g *PubService) RegisterGRPC(s *grpc.Server) {
	acctspb.RegisterRsaServiceServer(s, g)
}

func (g *PubService) PublicKey(ctx context.Context, req *acctspb.RsaRequest) (*wrapperspb.BytesValue, error) {
	if req.G {
		uniqueId, err := g.Unique(ctx, nil)
		if err != nil {
			return nil, err
		}
		req.Unique = uniqueId.Value
	}

	err := req.Validate()
	if err != nil {
		return nil, errors.ErrorArgument(ctx, err)
	}

	if !req.G {
		if err = g.usecase.ValidateUniqueId(ctx, req.GetUnique()); err != nil {
			return nil, err
		}
	}

	requestId := fmt.Sprintf("%s:%s", req.GetAction(), req.GetUnique())
	obj, err := g.usecase.FetchObj(ctx, requestId)
	if err != nil {
		return nil, err
	}

	return &wrapperspb.BytesValue{Value: obj.Public}, nil
}

func (g *PubService) Unique(ctx context.Context, _ *emptypb.Empty) (*wrapperspb.StringValue, error) {
	uniqueId, err := g.usecase.FetchUniqueId(ctx, time.Minute*10)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{Value: uniqueId}, nil
}
