package rdy

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/lestrrat-go/backoff/v2"
)

var (
	ErrTimeout  = errors.New("timeout")
	errNotReady = errors.New("not ready")
)

type WaitRequest struct {
	URL      string
	Backoff  backoff.Policy
	Reporter Reporter
}

func do(ctx context.Context, req WaitRequest) error {
	req.Reporter.L1(ctx, ">> Request: GET %s", req.URL)
	res, err := http.Get(req.URL)
	if err != nil {
		req.Reporter.L1(ctx, "failed performing request: %s", err.Error())
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	req.Reporter.L1(ctx, "<< Response: %s", res.Status)
	if content, err := io.ReadAll(res.Body); err == nil {
		req.Reporter.L2(ctx, "<< Response Body")
		req.Reporter.L2(ctx, string(content))
	}
	if res.StatusCode != http.StatusOK {
		return errNotReady
	}
	return nil
}

// Wait blocks until the given URL returns a 200 OK status code or the given ctx is Done.
func Wait(ctx context.Context, req WaitRequest) error {
	req.Reporter = &safeReporter{req.Reporter}
	rdy := make(chan struct{})
	go func() {
		b := req.Backoff.Start(ctx)
		for backoff.Continue(b) {
			select { // if context canceled.
			case <-ctx.Done():
				break
			default:
			}
			err := do(ctx, req)
			if err != nil {
				continue
			}
			close(rdy)
			break
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-rdy:
		return nil
	}
}
