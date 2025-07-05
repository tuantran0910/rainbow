FROM cubejs/cube:v1.3.21

# Label the image
LABEL org.opencontainers.image.authors="Tuan Tran <tntuan0910@gmail.com>"

# Set environment variables
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONHTTPSVERIFY=0

# Switch to root to install packages
USER root

# Install git and clean up to keep image small
RUN apt-get update && \
    apt-get install -y git inotify-tools && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install uv
COPY --from=ghcr.io/astral-sh/uv:0.6.10 /uv /uvx /bin/

# Install Python dependencies directly in CubeJS's Python environment
RUN uv pip install --system --no-cache-dir --break-system-packages \
    cube_dbt>=0.6.2 \
    google-cloud-storage>=3.1.1

# Switch back to non-root user
USER 1000
