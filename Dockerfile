FROM golang:1.17-bullseye AS build

WORKDIR /build

COPY go.mod .
COPY *.go ./

RUN go mod download

RUN go build -o /reconcile

FROM debian:bullseye-slim AS final

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN adduser --disabled-password --gecos "" reconciler

USER reconciler

WORKDIR /home/reconciler

COPY --from=build /reconcile /home/reconciler/reconcile

#ENTRYPOINT ["/home/reconsiler/reconsile"]
