// ==============================================
// Grafana Alloy OpenTelemetry Collector Config
// Production-Grade Baseline
//
// Pipelines:
// - OTLP receiver → memory limiter → batch processors
// - Metrics pipeline has transform to promote selected resource attributes
// - Exporters for Loki (logs), Tempo (traces), Mimir (metrics)
//
// Key Notes:
// - Memory limiter prevents OOM crashes
// - Batch processors tuned separately for logs/traces vs metrics
// - Only essential resource attributes promoted (avoid cardinality explosion)
// ==============================================


// ----------------------------------------------
// Receiver: OTLP gRPC
// Entry point for telemetry from instrumented services
// ----------------------------------------------
otelcol.receiver.otlp "ingest" {
  grpc {
    endpoint = "0.0.0.0:4317"
  }

  output {
    logs    = [otelcol.processor.memory_limiter.production.input]
    traces  = [otelcol.processor.memory_limiter.production.input]
    metrics = [otelcol.processor.memory_limiter.production.input]
  }
}


// ----------------------------------------------
// Processor: Memory limiter
// Prevents Alloy from crashing due to OOM
// ----------------------------------------------
otelcol.processor.memory_limiter "production" {
  limit          = "256MiB"   // Max memory Alloy can use
  spike_limit    = "64MiB"    // Extra buffer for short spikes
  check_interval = "2s"  // How often memory is checked

  output {
    logs    = [otelcol.processor.batch.default.input]
    traces  = [otelcol.processor.batch.default.input]
    metrics = [otelcol.processor.batch.metrics.input]
  }
}


// ----------------------------------------------
// Processor: Batch (logs/traces)
// Groups into batches before export
// Shorter timeout → lower latency
// ----------------------------------------------
otelcol.processor.batch "default" {
  timeout             = "200ms"
  send_batch_size     = 512
  send_batch_max_size = 1024

  output {
    logs   = [otelcol.exporter.loki.to_loki.input]
    traces = [otelcol.exporter.otlp.to_tempo.input]
  }
}


// ----------------------------------------------
// Processor: Batch (metrics)
// Larger batch size and longer timeout for efficiency
// ----------------------------------------------
otelcol.processor.batch "metrics" {
  timeout             = "5s"
  send_batch_size     = 1024
  send_batch_max_size = 2048

  output {
    metrics = [otelcol.processor.transform.resource_attributes.input]
  }
}


// ----------------------------------------------
// Processor: Transform
// Promote selected resource attributes to metric labels
// Avoids high cardinality by limiting to essentials
// ----------------------------------------------
otelcol.processor.transform "resource_attributes" {
  error_mode = "ignore"

  metric_statements {
    context = "datapoint"
    statements = [
      `set(datapoint.attributes["service_name"], resource.attributes["service.name"])`,
      `set(datapoint.attributes["service_version"], resource.attributes["service.version"])`,
      `set(datapoint.attributes["service_namespace"], resource.attributes["service.namespace"])`,
      `set(datapoint.attributes["deployment_environment"], resource.attributes["deployment.environment"])`,
      // `set(datapoint.attributes["host_name"], resource.attributes["host.name"])`,
    ]
  }

  output {
    metrics = [otelcol.exporter.prometheus.to_mimir.input]
  }
}


// ----------------------------------------------
// Exporter: Loki
// Sends logs to Grafana Loki
// ----------------------------------------------
otelcol.exporter.loki "to_loki" {
  forward_to = [loki.write.default.receiver]
}
loki.write "default" {
  endpoint {
    url = sys.env("LOKI_URL")
  }
}


// ----------------------------------------------
// Exporter: Tempo
// Sends traces to Grafana Tempo
// Currently insecure (dev only)
// ----------------------------------------------
otelcol.exporter.otlp "to_tempo" {
  client {
    endpoint = sys.env("TEMPO_ENDPOINT")
    tls {
      insecure             = true
      insecure_skip_verify = true
    }
  }
}


// ----------------------------------------------
// Exporter: Prometheus → Mimir
// Exposes metrics as Prometheus remote_write stream
// Attributes already promoted to labels above
// ----------------------------------------------
otelcol.exporter.prometheus "to_mimir" {
  forward_to = [prometheus.remote_write.mimir.receiver]
}

prometheus.remote_write "mimir" {
  endpoint {
    url = sys.env("MIMIR_URL")
  }
}

