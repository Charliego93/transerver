package data

import (
	"context"
	"github.com/transerver/accounts/internal/biz"
	"github.com/transerver/accounts/internal/ent"
	"github.com/transerver/accounts/internal/ent/region"
	"go.uber.org/zap"
)

var _ biz.RegionRepo = (*regionRepo)(nil)

type regionRepo struct {
	data   *Data
	logger *zap.Logger
}

func NewRegionRepo(data *Data, logger *zap.Logger) biz.RegionRepo {
	return &regionRepo{data: data, logger: logger}
}

func (g *regionRepo) FindByCode(ctx context.Context, code string) (*ent.Region, error) {
	return g.data.ent.Region.Query().Select(region.FieldCode, region.FieldArea).Where(region.Code(code)).First(ctx)
}

func (g *regionRepo) All(ctx context.Context, lang string) (ent.Regions, error) {
	return g.data.ent.Region.Query().All(ctx)
}
