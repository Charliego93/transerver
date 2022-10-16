package data

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/transerver/accounts/internal/biz"
	"github.com/transerver/commons/errors"
	"go.uber.org/zap"
	"time"
)

var _ biz.PubRepo = (*pubRepo)(nil)

type pubRepo struct {
	*Data
	logger *zap.Logger
}

func NewRsaRepo(data *Data, logger *zap.Logger) biz.PubRepo {
	return &pubRepo{Data: data, logger: logger}
}

func (g *pubRepo) FetchRsaObj(ctx context.Context, requestId string) (*biz.RsaObj, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	cmd := g.redis.Get(ctx, requestId)
	if cmd.Err() == redis.Nil {
		return nil, nil
	}
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	var rsaObj = &biz.RsaObj{}
	err := cmd.Scan(rsaObj)
	return rsaObj, err
}

func (g *pubRepo) StoreRsaObj(ctx context.Context, requestId string, expiration time.Duration, obj *biz.RsaObj) error {
	status, err := g.redis.Set(ctx, requestId, obj, expiration).Result()
	if err != nil {
		return err
	}
	if "OK" != status {
		return errors.New("store rsa obj error")
	}
	return nil
}

func (g *pubRepo) UniqueIdExists(ctx context.Context, uniqueId string) bool {
	r, err := g.redis.Exists(ctx, uniqueId).Result()
	if err != nil {
		return false
	}
	return r == 1
}

func (g *pubRepo) StoreUniqueId(ctx context.Context, uniqueId string, ttl time.Duration) error {
	r, err := g.redis.Set(ctx, uniqueId, "", ttl).Result()
	if err != nil {
		return err
	}
	if r != "OK" {
		return errors.New("store unique id error")
	}
	return nil
}
