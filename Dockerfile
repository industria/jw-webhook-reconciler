FROM golang:1.17-bullseye AS build

WORKDIR /build

COPY go.mod .
COPY *.go ./

RUN go mod download

RUN go build -o /reconsile

FROM debian:bullseye-slim AS final

RUN adduser --disabled-password --gecos "" reconsiler

USER reconsiler

WORKDIR /home/reconsiler

COPY --from=build /reconsile /home/reconsiler/reconsile

#ENTRYPOINT ["/home/reconsiler/reconsile"]
