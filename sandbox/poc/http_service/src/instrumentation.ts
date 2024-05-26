import { OTLPMetricExporter } from '@opentelemetry/exporter-metrics-otlp-grpc';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { PeriodicExportingMetricReader } from '@opentelemetry/sdk-metrics';
import { NodeSDK } from '@opentelemetry/sdk-node';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { HttpInstrumentation } from '@opentelemetry/instrumentation-http';
import { ExpressInstrumentation } from '@opentelemetry/instrumentation-express';
import { DiagConsoleLogger, DiagLogLevel, diag } from '@opentelemetry/api';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { NestInstrumentation } from '@opentelemetry/instrumentation-nestjs-core';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';

const traceExporter = new OTLPTraceExporter({
  timeoutMillis: 30000,
  url: `http://${process.env.OTLP_SERVICE_NAME}:4317`,
});

const metricExporter = new PeriodicExportingMetricReader({
  exporter: new OTLPMetricExporter({
    url: `http://${process.env.OTLP_SERVICE_NAME}:4317`,
  }),
  exportIntervalMillis: 10000,
});

const sdk = new NodeSDK({
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: process.env.SERVICE_NAME,
    [SemanticResourceAttributes.SERVICE_VERSION]: process.env.SERVICE_VERSION,
    [SemanticResourceAttributes.SERVICE_NAMESPACE]:
      process.env.OTLP_SERVICE_NAME,
    [SemanticResourceAttributes.DEPLOYMENT_ENVIRONMENT]:
      process.env.ENV_PLATFORM,
  }),
  traceExporter: traceExporter,
  metricReader: metricExporter,
  spanProcessor: new BatchSpanProcessor(traceExporter, {
    scheduledDelayMillis: 10,
  }),
  instrumentations: [getNodeAutoInstrumentations()],
});

diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.DEBUG);
sdk.start();
