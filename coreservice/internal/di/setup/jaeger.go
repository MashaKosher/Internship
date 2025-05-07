package setup

import (
	"coreservice/internal/di"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func mustJaeger(conf di.ConfigType) di.JaegerType {
	cfg := &config.Configuration{
		ServiceName: conf.Jaeger.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			LocalAgentHostPort:  conf.Jaeger.Host + ":" + conf.Jaeger.Port,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)
	return di.JaegerType{Tracer: tracer, Closer: closer}
}

func deferJaeger(closer io.Closer) {
	closer.Close()
}
