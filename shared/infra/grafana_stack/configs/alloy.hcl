otelcol.receiver.otlp "ingest" {
  grpc {
    endpoint = "0.0.0.0:4317"
  }

  output {
    logs   = [otelcol.processor.batch.default.input]
    traces = [otelcol.processor.batch.default.input]
  }
}

otelcol.processor.batch "default" {
  output {
    logs   = [otelcol.exporter.loki.to_loki.input]
    traces = [otelcol.exporter.otlp.to_tempo.input]
  }
}

otelcol.exporter.loki "to_loki" {
  // Send converted Loki entries to loki.write
  forward_to = [loki.write.default.receiver]
}
loki.write "default" {
  endpoint {
    url = sys.env("LOKI_URL")
  }
}

otelcol.exporter.otlp "to_tempo" {
  client {
    endpoint = sys.env("TEMPO_ENDPOINT")
    tls {
      insecure = true
      insecure_skip_verify = true
    }
  }
}
