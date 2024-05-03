const {
  DiagConsoleLogger,
  DiagLogLevel,
  diag,
} = require('@opentelemetry/api');
const { OTLPMetricExporter } = require('@opentelemetry/exporter-metrics-otlp-grpc');
const { Resource } = require('@opentelemetry/resources');
const { PeriodicExportingMetricReader } = require('@opentelemetry/sdk-metrics');
const { NodeSDK } = require('@opentelemetry/sdk-node');
const { SemanticAttributes } = require('@opentelemetry/semantic-conventions');

const otelSdk = new NodeSDK({
  resource: Resource.default().merge(
    new Resource({
      [SemanticAttributes.SERVICE_NAME]: 'nodejs app service',
      [SemanticAttributes.SERVICE_VERSION]: '1.0.0',
    }),
  ),
  metricReader: new PeriodicExportingMetricReader({
    exporter: new OTLPMetricExporter(),
    exportIntervalMillis: 10000,
  }),
});

process.on('SIGTERM', () => {
  otelSdk
    .shutdown()
    .then(
      () => console.log('otelSdk shutdown successful'),
      (err) => console.log('otelSdk shutdown error', err),
    )
    .finally(() => process.exit(0));
});
diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.VERBOSE);
otelSdk.start();
