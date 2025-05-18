package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HandlerProps struct {
	Method string
	Path   string
	Func   http.HandlerFunc
}

type TemplateData struct {
	Title              string
	ExternalCallURL    string
	ExternalCallMethod string
	RequestNameOTEL    string
	Content            string
	OTELTracer         trace.Tracer
}

type WebServer struct {
	Router       chi.Router
	Handlers     []HandlerProps
	TemplateData *TemplateData
}

func NewWebServer(service string) *WebServer {
	newWebServer := WebServer{
		Router:   chi.NewRouter(),
		Handlers: make([]HandlerProps, 0),
	}

	if service == "a" {
		newWebServer.AddHandler(http.MethodPost, "/", newWebServer.ZipCode)
	} else {
		newWebServer.AddHandler(http.MethodGet, "/", newWebServer.ZipCodeAndTemperature)
	}

	return &newWebServer
}

func (s *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	s.Handlers = append(s.Handlers, HandlerProps{
		Method: method,
		Path:   path,
		Func:   handler,
	})
}

func (s *WebServer) Start(ser string) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdown, err := initProvider(fmt.Sprintf("api_%s", ser), "otel-collector:4317")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %v", err)
		}
	}()

	tracer := otel.Tracer("microservice-tracer")
	s.TemplateData = &TemplateData{
		Title:              fmt.Sprintf("api-%s", ser),
		ExternalCallURL:    fmt.Sprintf("http://api-%s:8080", ser),
		ExternalCallMethod: "GET",
		RequestNameOTEL:    fmt.Sprintf("request-%s", ser),
		OTELTracer:         tracer,
	}

	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Logger)
	s.Router.Handle("/metrics", promhttp.Handler())
	for _, h := range s.Handlers {
		s.Router.Method(h.Method, h.Path, h.Func)
	}

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: s.Router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && http.ErrServerClosed != err {
			log.Printf("Error starting the server.")
			return
		}
	}()

	<-stop

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Cold not gracefully shutdown the server: %v\n", err)
	}

	println("Stop Serve")
}

func initProvider(serviceName, collectorUrl string) (func(context.Context) error, error) {
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create resource: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, collectorUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create gRPC connection to collector: %v", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("Failed to create trace exporter: %v", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return traceProvider.Shutdown, nil
}
