package service

import (
	"context"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/nats-io/nats.go"

	"github.com/cage1016/gae-custom-ws/internal/pkg/level"
)

const addTopic = "add"

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(AddService) AddService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
type AddService interface {
	// [method=post,expose=true]
	Sum(ctx context.Context, a int64, b int64) (res int64, err error)
	// [method=post,expose=true]
	Concat(ctx context.Context, a string, b string) (res string, err error)
}

// the concrete implementation of service interface
type stubAddService struct {
	logger log.Logger
	nc     *nats.Conn
}

// New return a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(nc *nats.Conn, logger log.Logger) (s AddService) {
	var svc AddService
	{
		svc = &stubAddService{nc: nc, logger: logger}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// Implement the business logic of Sum
func (ad *stubAddService) Sum(ctx context.Context, a int64, b int64) (res int64, err error) {
	res = a + b
	err = ad.nc.Publish(addTopic, []byte(strconv.FormatInt(res, 10)))
	if err != nil {
		level.Error(ad.logger).Log("method", "ad.nc.Publish", "value", strconv.FormatInt(res, 10), "err", err)
		return
	}

	return res, err
}

// Implement the business logic of Concat
func (ad *stubAddService) Concat(ctx context.Context, a string, b string) (res string, err error) {
	res = a + b
	err = ad.nc.Publish(addTopic, []byte(res))
	if err != nil {
		level.Error(ad.logger).Log("method", "ad.nc.Publish", "value", res, "err", err)
		return
	}

	return res, err
}
