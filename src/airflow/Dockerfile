# Define the base image
ARG AIRFLOW_VERSION=2.10.4
ARG PYTHON_VERSION=3.11
FROM apache/airflow:slim-${AIRFLOW_VERSION}-python${PYTHON_VERSION}

# Label the image
LABEL org.opencontainers.image.authors="Tuan Tran <tntuan0910@gmail.com>"

# Set this to avoid buffering stdout and stderr
ENV PYTHONUNBUFFERED=1

# Install additional system dependencies
USER root
RUN apt-get update && apt-get install -y libpq-dev gcc \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Use the virtual environment automatically
ENV VIRTUAL_ENV=/opt/venv

# Copy the dependency file and install dependencies
COPY pyproject.toml .
RUN uv pip install --no-cache --python ${PYTHON_VERSION} -r pyproject.toml

# Switch back to airflow user
USER airflow

# Copy the rest of the files
COPY . .

# Place entry points in the environment at the front of the path
ENV PATH=${VIRTUAL_ENV}/bin:${PATH}

# Post cleanup
USER root
RUN apt-get remove -y gcc && \
    apt-get autoremove -y && \
    rm -rf /tmp/* /var/tmp/* /var/lib/apt/lists/* /var/cache/apt/archives/*

# Switch back to airflow user
USER airflow
