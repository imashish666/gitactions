FROM golang:1.19-alpine AS build
WORKDIR /opt/code/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build main.go

FROM ubuntu:jammy-20230522
WORKDIR /opt/code/app
COPY --from=build /opt/code/app/main .
COPY ./prod_run.sh .
COPY ./filebeat.yml .
ARG region
ARG secret_name
ENV REGION_ENV=$region
ENV SECRET=$secret_name
ENV DEPLOYMENT="prod"
RUN apt-get update && \ 
    apt-get -y install curl && \ 
    curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-8.7.1-linux-x86_64.tar.gz && \
    tar xzvf filebeat-8.7.1-linux-x86_64.tar.gz && \ 
    rm filebeat-8.7.1-linux-x86_64.tar.gz && \ 
    mv filebeat-8.7.1-linux-x86_64 filebeat && \
    rm filebeat/filebeat.yml && \
    chmod go-w /opt/code/app/filebeat.yml && \
    cp filebeat.yml ./filebeat/ && \ 
    chmod +x /opt/code/app/prod_run.sh && \
    sed -i -e 's/\r$//' /opt/code/app/*.sh
EXPOSE 8080
ENTRYPOINT /opt/code/app/prod_run.sh