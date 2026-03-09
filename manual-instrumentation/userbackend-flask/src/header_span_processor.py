from contextvars import ContextVar
from typing import Mapping, Optional

from opentelemetry.sdk.trace import ReadableSpan, Span, SpanProcessor

request_headers_ctx: ContextVar[Optional[Mapping[str, str]]] = ContextVar(
    "request_headers_ctx", default=None
)


class HeaderSpanProcessor(SpanProcessor):
    def on_start(self, span: Span, parent_context=None) -> None:
        headers = request_headers_ctx.get() or {}
        traceparent = headers.get("traceparent")
        tracestate = headers.get("tracestate")
        if not traceparent:
            sc = span.get_span_context()
            flags = int(sc.trace_flags)
            traceparent = f"00-{sc.trace_id:032x}-{sc.span_id:016x}-{flags:02x}"

        span.set_attribute("http.request.header.traceparent", traceparent)
        if tracestate:
            span.set_attribute("http.request.header.tracestate", tracestate)

    def on_end(self, span: ReadableSpan) -> None:
        pass

    def shutdown(self) -> None:
        pass

    def force_flush(self, timeout_millis: int = 30000) -> bool:
        return True
