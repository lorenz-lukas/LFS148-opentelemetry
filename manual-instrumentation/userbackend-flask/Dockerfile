FROM python:3.12-slim
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENV POETRY_VERSION=1.8.3 \
    POETRY_HOME=/opt/poetry \
    POETRY_NO_INTERACTION=1 \
    POETRY_VIRTUALENVS_CREATE=false \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

RUN python -m pip install --upgrade pip \
    && python -m pip install --no-cache-dir "poetry==${POETRY_VERSION}"

COPY pyproject.toml /app/pyproject.toml
COPY requirements.txt /app/requirements.txt
RUN poetry install --no-root --only main

COPY . /app

EXPOSE 7000

CMD ["sh", "-c", "cd src && poetry run flask --app app run --host=0.0.0.0 --port=7000"]
