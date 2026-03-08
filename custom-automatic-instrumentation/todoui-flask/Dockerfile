FROM python:3.12-slim AS builder
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt /app/requirements.txt

RUN python -m pip install --upgrade pip \
    && python -m pip install --no-cache-dir --prefix=/install -r /app/requirements.txt \
    && PYTHONPATH=/install/lib/python3.12/site-packages PATH=/install/bin:${PATH} \
       opentelemetry-bootstrap --action=install

FROM python:3.12-slim AS runtime
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --system app && useradd --system --gid app --create-home --home-dir /home/app app

COPY --from=builder --chown=app:app /install /install
COPY --chown=app:app . /app

ENV PATH="/install/bin:${PATH}" \
    PYTHONPATH="/install/lib/python3.12/site-packages" \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

USER app

EXPOSE 7000

CMD ["opentelemetry-instrument","python","app.py"]
