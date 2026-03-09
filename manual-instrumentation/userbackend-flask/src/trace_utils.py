from os import getenv

from opentelemetry import trace as trace_api
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor, ConsoleSpanExporter
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from header_span_processor import HeaderSpanProcessor
from resource_utils import create_resource


def create_tracing_pipeline() -> BatchSpanProcessor:
    # exporter = ConsoleSpanExporter()
    endpoint = getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://otelcol:4317")
    exporter = OTLPSpanExporter(endpoint=endpoint, insecure=True)
    span_processor = BatchSpanProcessor(exporter)
    return span_processor


def create_tracer(name: str, version: str) -> trace_api.Tracer:
    rc = create_resource(name, version)
    processor = create_tracing_pipeline()
    provider = TracerProvider(resource=rc)
    provider.add_span_processor(HeaderSpanProcessor())
    provider.add_span_processor(processor)
    trace_api.set_tracer_provider(provider)
    tracer = trace_api.get_tracer(name, version)
    return tracer
