# pyright: reportMissingTypeStubs=false, reportUnknownParameterType=false, reportMissingParameterType=false, reportUnknownArgumentType=false, reportUnknownMemberType=false, reportAttributeAccessIssue=false

import json
import logging
from os import getenv
import time

import requests
from client import ChaosClient, FakerClient
from flask import Flask, Response, make_response, request
from header_span_processor import request_headers_ctx
from logging_utils import handler
from metric_utils import (
    create_meter,
    create_request_instruments,
    create_resource_instruments,
)
from prometheus_client import start_http_server
from opentelemetry import context, trace as trace_api
from opentelemetry.propagate import extract, inject
from opentelemetry.semconv.trace import SpanAttributes
from opentelemetry.trace import get_current_span
from opentelemetry.trace.status import Status, StatusCode
from trace_utils import create_tracer

logging.basicConfig(level=logging.INFO)
tracer = create_tracer("userbackend-flask", "0.1")
meter = create_meter("userbackend-flask", "0.1")

# global variables
app = Flask(__name__)
logger = logging.getLogger()
logger.addHandler(handler)
PORT = getenv("PORT", 7000)
db = ChaosClient(client=FakerClient())
workload_instruments = {}
rc_instruments = {}


def init_observability() -> None:
    global workload_instruments, rc_instruments, meter
    start_http_server(
        9464, addr="0.0.0.0"
    )  # prometheus metrics scrap endpoint :9464/metrics
    workload_instruments = create_request_instruments(meter)
    rc_instruments = create_resource_instruments(meter)


init_observability()


@app.before_request
def attach_context_with_trace_header():
    ctx = extract(request.headers)
    previous_ctx_token = context.attach(ctx)
    request.environ["previous_ctx_token"] = previous_ctx_token
    headers_map = {k.lower(): v for k, v in request.headers.items()}
    request.environ["headers_ctx_token"] = request_headers_ctx.set(headers_map)


@app.teardown_request
def restore_context_on_teardown(err):
    previous_ctx_token = request.environ.get("previous_ctx_token", None)
    if previous_ctx_token:
        context.detach(previous_ctx_token)
    headers_ctx_token = request.environ.get("headers_ctx_token", None)
    if headers_ctx_token:
        request_headers_ctx.reset(headers_ctx_token)


@app.before_request
def before_request():
    if workload_instruments:
        workload_instruments["traffic_volume"].add(
            1, attributes={"http.route": request.path}
        )
    request.environ["request_start"] = time.time_ns()


@app.after_request
def after_request(response: Response) -> Response:
    if workload_instruments:
        workload_instruments["request_latency"].record(
            amount=(time.time_ns() - request.environ["request_start"]) / 1_000_000_000,
            attributes={
                "http.request.method": request.method,
                "http.route": request.path,
                "http.response.status_code": response.status_code,
            },
        )
    return response


@app.route("/users", methods=["GET"])
@tracer.start_as_current_span("/users")
def get_user():
    user, status = db.get_user(123)
    logging.info(f"Found user {user!s} with status {status}")
    data = {}
    if user is not None:
        data = {"id": user.id, "name": user.name, "address": user.address}
    else:
        logging.warning(f"Could not find user with id {123}")
    logging.debug(f"Collected data is {data}")
    response = make_response(data, status)
    logging.debug(f"Generated response {response}")
    return response


@tracer.start_as_current_span("do_stuff")
def do_stuff():
    headers = {}
    inject(headers)
    time.sleep(0.1)
    url = "http://echo:80/"
    try:
        response = requests.get(url, headers=headers, timeout=3)
        logging.info("Headers included in outbound request:")
        logging.info(json.dumps(response.json()["request"]["headers"], indent=2))
        return response, response.status_code
    except requests.RequestException as exc:
        span = trace_api.get_current_span()
        span.record_exception(exc)
        span.set_status(Status(StatusCode.ERROR, str(exc)))
        return None, 500
    except (ValueError, KeyError, TypeError):
        app.logger.warning("Failed to parse downstream headers from echo response")
        return response, response.status_code


@app.route("/")
@tracer.start_as_current_span("/")
def index():
    status_code = 200
    downstream_response, downstream_status = do_stuff()
    if downstream_response is None:
        status_code = 500

    span = get_current_span()
    span.set_attributes(
        {
            SpanAttributes.HTTP_REQUEST_METHOD: request.method,
            SpanAttributes.URL_PATH: request.path,
            SpanAttributes.HTTP_RESPONSE_STATUS_CODE: status_code,
            "downstream.http.status_code": downstream_status,
        }
    )
    if downstream_response is None:
        return make_response({"error": "Failed to reach downstream service"}, 500)

    logging.info("Info from the index function")
    current_time = time.strftime("%a, %d %b %Y %H:%M:%S", time.gmtime())
    return f"Hello, World! It's currently {current_time}"


@app.route("/health")
@tracer.start_as_current_span("/health")
def health():
    app.logger.info("Health check")
    return "OK", 200


@app.route("/error")
@tracer.start_as_current_span("/error")
def error():
    app.logger.info("Error endpoint called")
    return "Something went wrong", 500


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=PORT, debug=True)
