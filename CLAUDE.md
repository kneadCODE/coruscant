# CLAUDE.md

This file is for Claude.

The objective of this project is to be able to learn and apply industry best practices at an enterprise grade/scale, but using OSS/Free resources as much as possible. I am a Solution Architect, Principal Engineer & Engineering Manager, but all my office work is private and I am lot of times limited to office policies which restricts me from trying out/applying the industry practices as of Aug 2025. This project will address that and also act as a showcase as well. I am experienced in Go, backend systems and infra. But am also interested to increase depth and breadth in all aspects.

## Observability Stack Decisions & Learnings

### OpenTelemetry Best Practices
- **RELY ON AUTO-INSTRUMENTATION**: OTEL libraries handle most attributes automatically. Avoid manual attribute addition unless absolutely necessary.
- **Follow OTEL Semantic Conventions**: Stick to standard conventions rather than custom implementations.
- **Log-Trace Correlation**: OTEL slog bridge automatically handles trace context correlation - no manual trace ID injection needed.

### Grafana Stack Configuration  
- **Traces Navigation**: Added dedicated Tempo sidebar navigation for direct trace access in Grafana.
- **Dashboard Location**: Dashboards mounted at `/var/lib/grafana/dashboards` with provisioning config.
- **Logs-Traces Correlation**: Configure Loki derived fields with regex `(?:trace_id=|traceID=)([a-fA-F0-9]{32})` for flexible trace ID formats.

### Local Development Settings
- **No Compression**: Skip OTLP compression for localhost collectors (no network benefit, adds CPU overhead).
- **Port Configuration**: Alloy (4317) → Tempo OTLP (4319), avoid conflicts with internal gRPC (9095).
- **Environment Variables**: Use template variables in Grafana configs instead of hardcoded URLs.

### HTTP Tracing Key Points
- **Missing Middleware**: Always add `httpserver.WithTracing()` to enable distributed tracing.
- **Span Events**: Use `trace.WithTimestamp(time.Now())` for absolute timestamps, not relative.
- **Health Endpoint Filtering**: Always filter `/_/ping`, `/_/ready`, `/_/health`, `/_/metrics` from tracing.

### Testing & Validation
- **Build Commands**: `go build ./...` for compilation, `golangci-lint run --fix` for linting.
- **Test Coverage**: Achieved comprehensive coverage across telemetry package with improved Measure function tests.
- **Error Patterns**: DNS resolver errors in tests are expected for local OTLP endpoints.
- **Test Cleanup**: Ensure cleanup functions receive context parameter: `defer cleanup(ctx)` not `defer cleanup()`.

### Measure Function API Design
- **Flexible Types**: Enhanced Measure function to accept `...any` instead of `...string` for attributes.
- **Type Safety**: Keys must be strings, values support: string, int, int64, float64, bool, with fallback string conversion.
- **Usage Pattern**: `ctx, finish := Measure(ctx, "operation", "key1", value1, "key2", value2)`
- **Error Handling**: Invalid keys (non-strings) are safely skipped, odd number of attributes ignored.

### Grafana Configuration Fixes
- **TraceQL Syntax**: Use `{resource.service.name=~"service|names"}` for valid TraceQL queries, not `{service.name=~"..."}`
- **Tempo Traces to Logs**: Configure `tracesToLogs` with service tag mapping and time shifts for context correlation.
- **Tempo Traces to Metrics**: Set up `tracesToMetrics` with span metrics queries for request rate, error rate, and duration.
- **Derived Fields URLs**: Use proper Grafana explore URLs for trace/span navigation from log entries.
- **Sidebar Navigation**: Configured Traces section in grafana.ini navigation for direct access to trace explorer.
