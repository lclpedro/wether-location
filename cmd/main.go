package main

import (
	"context"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/lclpedro/weather-location/configs"
	"github.com/lclpedro/weather-location/internal/scaffold/services"
	"github.com/lclpedro/weather-location/internal/scaffold/views"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func initProvider(serviceName, collectorEndpoint string) *sdktrace.TracerProvider {
	ctx := context.Background()
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(serviceName),
	))

	if err != nil {
		log.Fatal("Error to mount resource service", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		collectorEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	if err != nil {
		log.Fatal("Error to connect grpc service", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {

	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider
}

func init() {
	configs.InitConfigs()
}

func main() {

	tp := initProvider("weather-location-app", "otel-collector:4317")
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal("Error to shutdown services otlp", err)
		}
	}()

	trace := otel.Tracer("weather-location-tracer")

	app := fiber.New()
	app.Use(otelfiber.Middleware())

	allServices := services.NewAllServices(trace)
	app = views.NewAllHandlerViews(app, trace, allServices)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Error to up service", err)
	}
}
