# Define the base image
ARG PYTHON_VERSION=3.12
FROM python:${PYTHON_VERSION}-slim

# Build argument to determine environment (local or cloud)
ARG ENVIRONMENT=local
ENV DAGSTER_ENVIRONMENT=${ENVIRONMENT}

# Label the image
LABEL org.opencontainers.image.authors="Tuan Tran <tntuan0910@gmail.com>"
LABEL org.opencontainers.image.description="Rainbow Dagster Workspace - ${ENVIRONMENT} environment"

# Set environment variables
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    DAGSTER_WORKSPACE=/opt/dagster/app

# Set working directory
WORKDIR $DAGSTER_WORKSPACE

# Install additional system dependencies
USER root
RUN apt-get update && apt-get install -y \
    libpq-dev \
    gcc \
    g++ \
    build-essential \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Use the virtual environment automatically
ENV VIRTUAL_ENV=/opt/venv

# Copy the dependency file and install dependencies
COPY --from=ghcr.io/astral-sh/uv:0.6.10 /uv /uvx /bin/
COPY pyproject.toml .

# Install dependencies based on environment
RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    if [ "$ENVIRONMENT" = "cloud" ]; then \
        echo "Installing cloud dependencies..." && \
        uv sync --group cloud --no-dev --no-cache --python ${PYTHON_VERSION} --active; \
    elif [ "$ENVIRONMENT" = "local" ]; then \
        echo "Installing local dependencies..." && \
        uv sync --group local --no-cache --python ${PYTHON_VERSION} --active; \
    else \
        echo "Unknown environment: $ENVIRONMENT. Defaulting to local..." && \
        uv sync --group local --no-cache --python ${PYTHON_VERSION} --active; \
    fi

# Place entry points in the environment at the front of the path
ENV PATH=${VIRTUAL_ENV}/bin:${PATH}

# Post cleanup
RUN apt-get remove -y gcc g++ build-essential && \
    apt-get autoremove -y && \
    rm -rf /tmp/* /var/tmp/* /var/lib/apt/lists/* /var/cache/apt/archives/*

# Copy the dbt project and generate the manifest
COPY dbt dbt
RUN echo "Installing dbt packages for ${ENVIRONMENT} environment..." && \
    dbt deps --project-dir dbt/${ENVIRONMENT} && \
    echo "Generating dbt manifest for ${ENVIRONMENT} environment..." && \
    dbt parse --project-dir dbt/${ENVIRONMENT} --profiles-dir dbt/${ENVIRONMENT} && \
    echo "dbt manifest generated successfully"

# Copy all application code
COPY . .

# Switch back to non-root user
RUN chown -R 1000:1000 ${DAGSTER_WORKSPACE} ${VIRTUAL_ENV}
USER 1000

# Expose the port
EXPOSE 4000
