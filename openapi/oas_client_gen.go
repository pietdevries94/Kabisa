// Code generated by ogen, DO NOT EDIT.

package openapi

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/otelogen"
	"github.com/ogen-go/ogen/uri"
)

func trimTrailingSlashes(u *url.URL) {
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawPath = strings.TrimRight(u.RawPath, "/")
}

// Invoker invokes operations described by OpenAPI v3 specification.
type Invoker interface {
	// CreateNewQuoteGame invokes createNewQuoteGame operation.
	//
	// The quote game returns three quotes and three authors. In `PUT /quote-game/:id`, the player can
	// respond with their answer. There is a deadline of five minutes.
	//
	// POST /quote-game
	CreateNewQuoteGame(ctx context.Context) (CreateNewQuoteGameRes, error)
	// GetRandomQuote invokes getRandomQuote operation.
	//
	// Returns a random quote.
	//
	// GET /quote
	GetRandomQuote(ctx context.Context) (GetRandomQuoteRes, error)
	// SubmitAnswerForQuoteGame invokes submitAnswerForQuoteGame operation.
	//
	// This request expects an answer from the user and will return if the answer was correct and what
	// the correct answer should be.
	//
	// POST /quote-game/{id}/answer
	SubmitAnswerForQuoteGame(ctx context.Context, request []QuoteGameAnswer, params SubmitAnswerForQuoteGameParams) (SubmitAnswerForQuoteGameRes, error)
}

// Client implements OAS client.
type Client struct {
	serverURL *url.URL
	baseClient
}

var _ Handler = struct {
	*Client
}{}

// NewClient initializes new Client defined by OAS.
func NewClient(serverURL string, opts ...ClientOption) (*Client, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	trimTrailingSlashes(u)

	c, err := newClientConfig(opts...).baseClient()
	if err != nil {
		return nil, err
	}
	return &Client{
		serverURL:  u,
		baseClient: c,
	}, nil
}

type serverURLKey struct{}

// WithServerURL sets context key to override server URL.
func WithServerURL(ctx context.Context, u *url.URL) context.Context {
	return context.WithValue(ctx, serverURLKey{}, u)
}

func (c *Client) requestURL(ctx context.Context) *url.URL {
	u, ok := ctx.Value(serverURLKey{}).(*url.URL)
	if !ok {
		return c.serverURL
	}
	return u
}

// CreateNewQuoteGame invokes createNewQuoteGame operation.
//
// The quote game returns three quotes and three authors. In `PUT /quote-game/:id`, the player can
// respond with their answer. There is a deadline of five minutes.
//
// POST /quote-game
func (c *Client) CreateNewQuoteGame(ctx context.Context) (CreateNewQuoteGameRes, error) {
	res, err := c.sendCreateNewQuoteGame(ctx)
	return res, err
}

func (c *Client) sendCreateNewQuoteGame(ctx context.Context) (res CreateNewQuoteGameRes, err error) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("createNewQuoteGame"),
		semconv.HTTPRequestMethodKey.String("POST"),
		semconv.HTTPRouteKey.String("/quote-game"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(elapsedDuration)/float64(time.Millisecond), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, CreateNewQuoteGameOperation,
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/quote-game"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeCreateNewQuoteGameResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// GetRandomQuote invokes getRandomQuote operation.
//
// Returns a random quote.
//
// GET /quote
func (c *Client) GetRandomQuote(ctx context.Context) (GetRandomQuoteRes, error) {
	res, err := c.sendGetRandomQuote(ctx)
	return res, err
}

func (c *Client) sendGetRandomQuote(ctx context.Context) (res GetRandomQuoteRes, err error) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("getRandomQuote"),
		semconv.HTTPRequestMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/quote"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(elapsedDuration)/float64(time.Millisecond), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, GetRandomQuoteOperation,
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/quote"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeGetRandomQuoteResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// SubmitAnswerForQuoteGame invokes submitAnswerForQuoteGame operation.
//
// This request expects an answer from the user and will return if the answer was correct and what
// the correct answer should be.
//
// POST /quote-game/{id}/answer
func (c *Client) SubmitAnswerForQuoteGame(ctx context.Context, request []QuoteGameAnswer, params SubmitAnswerForQuoteGameParams) (SubmitAnswerForQuoteGameRes, error) {
	res, err := c.sendSubmitAnswerForQuoteGame(ctx, request, params)
	return res, err
}

func (c *Client) sendSubmitAnswerForQuoteGame(ctx context.Context, request []QuoteGameAnswer, params SubmitAnswerForQuoteGameParams) (res SubmitAnswerForQuoteGameRes, err error) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("submitAnswerForQuoteGame"),
		semconv.HTTPRequestMethodKey.String("POST"),
		semconv.HTTPRouteKey.String("/quote-game/{id}/answer"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(elapsedDuration)/float64(time.Millisecond), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, SubmitAnswerForQuoteGameOperation,
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [3]string
	pathParts[0] = "/quote-game/"
	{
		// Encode "id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			if unwrapped := string(params.ID); true {
				return e.EncodeValue(conv.StringToString(unwrapped))
			}
			return nil
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	pathParts[2] = "/answer"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeSubmitAnswerForQuoteGameRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeSubmitAnswerForQuoteGameResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}
