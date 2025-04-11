FROM public.ecr.aws/docker/library/golang:1.24 AS builder

COPY go.mod /app/go.mod
COPY go.sum /app/go.sum

WORKDIR /app
RUN go mod download

FROM builder AS builder2

COPY . /app
WORKDIR /app/cmd/app

RUN CGO_ENABLED=0 go build -o app_binary

FROM debian:bullseye-slim

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get upgrade -y && apt-get install -y git

RUN groupadd --gid 5001 nonroot && useradd --create-home --system --uid 5001 --gid 5001 nonroot

USER nonroot:nonroot
COPY --from=builder2 --chown=nonroot:nonroot /app/cmd/app/app_binary /app/app_binary

WORKDIR /app
ENTRYPOINT ["/app/app_binary"]
