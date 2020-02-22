package srv

import (
	"context"
	"sync"

	"github.com/micro/go-micro/server"
)

// WaitGroupWrapper is a handler wrapper which adds a handler to a sync.WaitGroup
func WaitGroupWrapper(wg *sync.WaitGroup) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			wg.Add(1)
			defer wg.Done()
			return h(ctx, req, rsp)
		}
	}
}

