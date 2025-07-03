FROM apache/superset:v5.0.0

USER root

RUN apt-get update && apt-get install -y \
    python3-dev \
    default-libmysqlclient-dev \
    build-essential \
    pkg-config

RUN pip install --no-cache \
    Authlib==1.5.2 \
    psycopg2-binary==2.9.10

USER superset
