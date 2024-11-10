FROM docker.io/library/ubuntu:24.04

ENV STOPS_DATABASE_PATH=/app/data/stops.db

WORKDIR /app

COPY database_generator* /app

RUN apt-get update && apt install -y python3 && rm -rf /var/lib/apt/lists/*

CMD ["/bin/bash", "/app/database_generator.sh"]