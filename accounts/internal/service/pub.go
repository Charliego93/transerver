package service

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/transerver/accounts/internal/biz"
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

func (g *PubService) RegisterHTTP(s *runtime.ServeMux) error {
	return acctspb.RegisterRsaServiceHandlerServer(context.Background(), s, g)
}

func (g *PubService) Routers() ([]string, []string) {
	return nil, nil
}

func (g *PubService) PublicKey(_ context.Context, req *acctspb.RsaRequest) (*wrapperspb.BytesValue, error) {
	requestId := fmt.Sprintf("%s:%s", req.GetAction(), req.GetUnique())
	obj, err := g.usecase.FetchObj(requestId)
	if err != nil {
		return nil, err
	}

	return &wrapperspb.BytesValue{Value: obj.Public}, nil
}

func (g *PubService) Unique(context.Context, *emptypb.Empty) (*wrapperspb.StringValue, error) {
	uniqueId, err := g.usecase.FetchUniqueId(time.Minute * 10)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{Value: uniqueId}, nil
}
