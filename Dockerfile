FROM golang:latest

COPY payment /
ENTRYPOINT ["/payment", "8000"]