from flask import Flask, render_template, request, jsonify, redirect, url_for
from functools import wraps
import logging
import requests
from os import getenv

# Import the trace API
from opentelemetry import trace

# Acquire a tracer
tracer = trace.get_tracer("todo.tracer")


def trace_span(span_name, attribute_extractors=None):
    """
    Decorator parametrizável para adicionar tracing a uma função.

    Args:
        span_name: Nome do span
        attribute_extractors: Dict com {nome_atributo: callable} para extrair valores

    Exemplo:
        @trace_span("add", {
            "todo.value": lambda: request.form.get('todo'),
            "backend.url": lambda: app.config['BACKEND_URL']
        })
    """

    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            with tracer.start_as_current_span(span_name) as span:
                # Extrai e seta os atributos
                if attribute_extractors:
                    for attr_name, extractor in attribute_extractors.items():
                        try:
                            value = extractor()
                            span.set_attribute(attr_name, value)
                        except Exception as e:
                            logging.debug(f"Erro ao extrair atributo {attr_name}: {e}")

                return func(*args, **kwargs)

        return wrapper

    return decorator


app = Flask(__name__)
logging.getLogger(__name__)
logging.basicConfig(
    format="%(levelname)s:%(name)s:%(module)s:%(message)s", level=logging.INFO
)

# Set a default external API URL
# Override the default URL if an environment variable is set
PORT = getenv("PORT", 7000)
app.config["BACKEND_URL"] = getenv("BACKEND_URL", f"http://localhost:{PORT}/todos/")


@app.route("/")
def index():

    backend_url = app.config["BACKEND_URL"]
    response = requests.get(backend_url)

    logging.info("GET %s/todos/", backend_url)
    if response.status_code == 200:
        # Print out the response content
        # print(response.text)
        logging.info("Response: %s", response.text)
        todos = response.json()

    return render_template("index.html", todos=todos)


@app.route("/add", methods=["POST"])
@trace_span(
    "add",
    {
        "todo.value": lambda: request.form.get("todo"),
        "backend.url": lambda: app.config["BACKEND_URL"],
    },
)
def add():
    if request.method == "POST":
        new_todo = request.form["todo"]
        logging.info("POST  %s/todos/%s", app.config["BACKEND_URL"], new_todo)
        response = requests.post(app.config["BACKEND_URL"] + new_todo)
    return redirect(url_for("index"))


@app.route("/delete", methods=["POST"])
@trace_span(
    "delete",
    {
        "todo.value": lambda: request.form.get("todo"),
        "backend.url": lambda: app.config["BACKEND_URL"],
    },
)
def delete():
    if request.method == "POST":
        delete_todo = request.form["todo"]
        logging.info("DELETE %s/todos/%s", app.config["BACKEND_URL"], delete_todo)
        print(delete_todo)
    response = requests.delete(app.config["BACKEND_URL"] + delete_todo)
    return redirect(url_for("index"))


@app.route("/actuator/health")
@trace_span("health")
def health():
    try:
        # Tentar conectar ao backend para validar saúde completa
        backend_url = app.config["BACKEND_URL"]
        response = requests.get(backend_url, timeout=2)
        if response.status_code == 200:
            return jsonify(status="UP", checks={"backend": "UP"}), 200
        else:
            return jsonify(status="DOWN", checks={"backend": "DOWN"}), 503
    except Exception as e:
        logging.warning(f"Health check falhou: {e}")
        # Ainda consideramos UP se o serviço está respondendo
        return jsonify(status="UP", message="Flask is running"), 200


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=PORT)
    # app.run(debug=True) # doesn't work with auto instrumentation
