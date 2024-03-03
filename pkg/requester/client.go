package requester

import (
	"context"
	"crypto/tls"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"
)

type (
	Client interface {
		Get(url string) (*http.Response, error)
	}

	requester struct {
		client *http.Client
		ctx    context.Context
		tracer trace.Tracer
	}
)

func NewRequester(ctx context.Context, tracer trace.Tracer, timeout int) Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(timeout) * time.Millisecond}
	return &requester{client: client, ctx: ctx, tracer: tracer}
}

func (r *requester) Get(url string) (*http.Response, error) {
	ctx, spam := r.tracer.Start(r.ctx, url)
	r.ctx = ctx
	defer spam.End()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(r.ctx, propagation.HeaderCarrier(req.Header))
	fmt.Println(req.Header)
	return r.client.Do(req)
}
