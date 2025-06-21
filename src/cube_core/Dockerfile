FROM cubejs/cube:v1.3.21

# Label the image
LABEL org.opencontainers.image.authors="Tuan Tran <tntuan0910@gmail.com>"

# Set environment variables
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

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

# Switch back to non-root user
RUN chown -R 1000:1000 ${VIRTUAL_ENV}
USER 1000
