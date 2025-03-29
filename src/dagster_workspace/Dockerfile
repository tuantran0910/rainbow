# Define the base image
ARG PYTHON_VERSION=3.12
FROM python:${PYTHON_VERSION}-slim

# Label the image
LABEL org.opencontainers.image.authors="Tuan Tran <tntuan0910@gmail.com>"

# Set environment variables
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    DAGSTER_WORKSPACE=/opt/dagster/app

# Set working directory
WORKDIR $DAGSTER_WORKSPACE

# Install additional system dependencies
USER root
RUN apt-get update && apt-get install -y libpq-dev gcc \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Use the virtual environment automatically
ENV VIRTUAL_ENV=/opt/venv

# Copy the dependency file and install dependencies
COPY --from=ghcr.io/astral-sh/uv:0.6.10 /uv /uvx /bin/
COPY pyproject.toml .
RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    uv sync --no-cache --python ${PYTHON_VERSION} --active

# Place entry points in the environment at the front of the path
ENV PATH=${VIRTUAL_ENV}/bin:${PATH}

# Post cleanup
RUN apt-get remove -y gcc && \
    apt-get autoremove -y && \
    rm -rf /tmp/* /var/tmp/* /var/lib/apt/lists/* /var/cache/apt/archives/*

# Copy Dagster instance YAML
COPY definitions.py .

# Switch back to non-root user
USER 1000

# Expose the port
EXPOSE 4000

# Run the Dagster GRPC server
CMD ["dagster", "api", "grpc", "-h", "0.0.0.0", "-p", "4000", "-f", "definitions.py"]
