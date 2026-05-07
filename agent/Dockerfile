FROM python:3.12-alpine
WORKDIR /app
COPY src ./src
ENV PYTHONPATH=/app/src
EXPOSE 8090
CMD ["python", "-m", "agent_core.server"]
