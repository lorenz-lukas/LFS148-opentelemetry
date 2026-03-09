from os import getenv

from opentelemetry.sdk._logs import LoggerProvider, LoggingHandler
from opentelemetry.sdk._logs.export import ConsoleLogExporter, SimpleLogRecordProcessor
from opentelemetry.sdk.resources import Resource
from opentelemetry.exporter.otlp.proto.grpc._log_exporter import OTLPLogExporter


logger_provider = LoggerProvider(
    resource=Resource.create(
        {
            "service.name": "api",
        }
    ),
)

# logger_provider.add_log_record_processor(SimpleLogRecordProcessor(exporter=ConsoleLogExporter()))
endpoint = getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://otelcol:4317")
logger_provider.add_log_record_processor(
    SimpleLogRecordProcessor(exporter=OTLPLogExporter(endpoint=endpoint, insecure=True))
)
handler = LoggingHandler(logger_provider=logger_provider)
