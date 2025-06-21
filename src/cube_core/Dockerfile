FROM cubejs/cube:v1.3.21

# Label the image
LABEL org.opencontainers.image.authors="Tuan Tran <tntuan0910@gmail.com>"

# Set environment variables
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

# Upgrade packages
USER root
RUN apt-get update && apt-get install -y \
    && rm -rf /var/lib/apt/lists/*

# Use the virtual environment automatically
ENV VIRTUAL_ENV=/opt/venv

# Copy the dependency file and install dependencies
ARG PYTHON_VERSION=3.12
COPY --from=ghcr.io/astral-sh/uv:0.6.10 /uv /uvx /bin/
COPY pyproject.toml .
RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    uv sync --no-cache --python ${PYTHON_VERSION} --active

#Place entry points in the environment at the front of the path
ENV PATH=${VIRTUAL_ENV}/bin:${PATH}

# Post cleanup
RUN apt-get remove -y gcc g++ build-essential && \
    apt-get autoremove -y && \
    rm -rf /tmp/* /var/tmp/* /var/lib/apt/lists/* /var/cache/apt/archives/*

# Switch back to non-root user
RUN chown -R 1000:1000 ${DAGSTER_HOME} ${VIRTUAL_ENV}
USER 1000
