package handler

import (
	"context"
	"time"

	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/cache/proto"
	"github.com/micro/services/pkg/cache"
)

type Cache struct{}

func (c *Cache) Get(ctx context.Context, req *pb.GetRequest, rsp *pb.GetResponse) error {
	if len(req.Key) == 0 {
		return errors.BadRequest("cache.get", "missing key")
	}

	var value string

	if err := cache.Context(ctx).Get(req.Key, &value); err != nil {
		return errors.InternalServerError("cache.get", err.Error())
	}

	rsp.Key = req.Key
	// set the value
	rsp.Value = value

	return nil
}

func (c *Cache) Set(ctx context.Context, req *pb.SetRequest, rsp *pb.SetResponse) error {
	if len(req.Key) == 0 {
		return errors.BadRequest("cache.set", "missing key")
	}

	// max size 1mb e.g byte * 1024 * 1024
	if len(req.Value) > 1e6 {
		return errors.BadRequest("cache.set", "value is too big")
	}

	ttl := time.Time{}

	if req.Ttl > 0 {
		ttl = time.Now().Add(time.Duration(req.Ttl) * time.Second)
	}

	if err := cache.Context(ctx).Set(req.Key, req.Value, ttl); err != nil {
		return errors.InternalServerError("cache.set", err.Error())
	}

	rsp.Status = "ok"

	return nil
}

func (c *Cache) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	if len(req.Key) == 0 {
		return errors.BadRequest("cache.delete", "missing key")
	}
	if err := cache.Context(ctx).Delete(req.Key); err != nil {
		return errors.InternalServerError("cache.delete", err.Error())
	}

	rsp.Status = "ok"

	return nil
}

func (c *Cache) Increment(ctx context.Context, req *pb.IncrementRequest, rsp *pb.IncrementResponse) error {
	if len(req.Key) == 0 {
		return errors.BadRequest("cache.increment", "missing key")
	}

	// increment the value
	v, err := cache.Context(ctx).Increment(req.Key, req.Value)
	if err != nil {
		return errors.InternalServerError("cache.increment", err.Error())
	}

	// set the response value
	rsp.Key = req.Key
	rsp.Value = v

	return nil
}

func (c *Cache) Decrement(ctx context.Context, req *pb.DecrementRequest, rsp *pb.DecrementResponse) error {
	if len(req.Key) == 0 {
		return errors.BadRequest("cache.decrement", "missing key")
	}

	v, err := cache.Context(ctx).Decrement(req.Key, req.Value)
	if err != nil {
		return errors.InternalServerError("cache.decrement", err.Error())
	}

	// set the response value
	rsp.Key = req.Key
	rsp.Value = v

	return nil
}
