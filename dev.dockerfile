FROM golang:1.19-alpine AS build
WORKDIR /opt/code/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build main.go

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM ubuntu:jammy-20230522
WORKDIR /opt/code/app
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /opt/code/app/main .
COPY ./dev_run.sh .
ARG region
ARG secret_name
ENV REGION_ENV=${region}
ENV SECRET=${secret_name}
ENV DEPLOYMENT="dev"
EXPOSE 8080
ENTRYPOINT /opt/code/app/dev_run.sh