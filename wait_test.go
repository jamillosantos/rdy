package rdy_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/lestrrat-go/backoff/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jamillosantos/rdy"
)

var _ = Describe("Wait", func() {
	When("the API becomes ready", func() {
		It("should wait until the API is ready", func() {
			now := time.Now()
			waitTime := 50 * time.Millisecond
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if time.Since(now) < waitTime {
					w.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				log.Printf("Server ready at %s", time.Now().Format(time.RFC3339))
				w.WriteHeader(http.StatusOK)
			}))
			DeferCleanup(func() {
				server.Close()
			})

			ctx := context.Background()
			waitingSince := time.Now()
			err := rdy.Wait(ctx, rdy.WaitRequest{
				URL: server.URL,
				Backoff: backoff.Constant(
					backoff.WithInterval(time.Millisecond),
					backoff.WithMaxRetries(0),
				),
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(time.Since(waitingSince)).To(BeNumerically("~", waitTime, 5*time.Millisecond), "should wait approximately the expected time")
		})
	})

	When("the context times out before the API become ready", func() {
		It("should wait until the API is ready", func() {
			now := time.Now()
			waitTime := 50 * time.Millisecond
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if time.Since(now) < waitTime {
					w.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				log.Printf("Server ready at %s", time.Now().Format(time.RFC3339))
				w.WriteHeader(http.StatusOK)
			}))
			DeferCleanup(func() {
				server.Close()
			})

			ctx, cancel := context.WithTimeout(context.Background(), waitTime/2)
			defer cancel()

			waitingSince := time.Now()
			err := rdy.Wait(ctx, rdy.WaitRequest{
				URL: server.URL,
				Backoff: backoff.Constant(
					backoff.WithInterval(time.Millisecond),
					backoff.WithMaxRetries(0),
				),
			})
			Expect(err).To(MatchError(context.DeadlineExceeded))
			Expect(time.Since(waitingSince)).To(BeNumerically("~", waitTime/2, 5*time.Millisecond), "should wait approximately the expected time")
		})
	})
})
