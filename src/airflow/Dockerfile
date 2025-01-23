# Define the base image
ARG AIRFLOW_VERSION=2.10.4
ARG PYTHON_VERSION=3.11
FROM apache/airflow:slim-${AIRFLOW_VERSION}-python${PYTHON_VERSION}

# Label the image maintainer
LABEL maintainer="Tuan Tran < tntuan0910@gmail.com"

# Set this to avoid buffering stdout and stderr
ENV PYTHONUNBUFFERED=1

# Install additional system dependencies
USER root
RUN apt-get update && apt-get install -y libpq-dev gcc && apt-get clean

# Switch back to airflow user
USER airflow

# Install additional Python packages
COPY requirements.txt /requirements.txt
RUN pip install --upgrade pip && \
    pip install --no-cache-dir -r /requirements.txt

# Post cleanup
USER root
RUN apt-get remove -y gcc && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /requirements.txt

# Switch back to airflow user
USER airflow
