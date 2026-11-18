# No third-party dependencies, so the image needs only the interpreter and the package.
FROM python:3.14-slim
WORKDIR /app
COPY worker ./worker
ENV PYTHONUNBUFFERED=1
USER nobody
ENTRYPOINT ["python", "-m", "worker"]
