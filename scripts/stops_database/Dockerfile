FROM docker.io/library/python:3.10

ENV STOPS_DATABASE_PATH=/app/data/stops.db

WORKDIR /app

COPY database_generator* /app

CMD ["/bin/bash", "/app/database_generator.sh"]