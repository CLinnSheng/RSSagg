FROM golang:1.23.6-bookworm

# Install wget for downloading wait-for-it
RUN apt-get update && apt-get install -y wget

# Download wait-for-it script
RUN wget -O /usr/local/bin/wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh \
    && chmod +x /usr/local/bin/wait-for-it.sh

# Create code src directory
RUN mkdir -p /usr/src

# Copy all the source code from host to container
COPY ./src /usr/src

# Change the working directory
WORKDIR /usr/src/

# Download the dependencies
RUN go get -d -v ./...

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build the go code
RUN go build

# EXPOSE the port
EXPOSE 8080

# ENV
ENV SCHEMA_DIR=/usr/src/sql/schema

# Wait for database to be ready, then run migrations and start the app
CMD wait-for-it.sh go_db:5432 -t 60 -- sh -c 'goose -dir $SCHEMA_DIR postgres "postgres://postgres:postgres@go_db:5432/rssagg" up && ./RSSagg'