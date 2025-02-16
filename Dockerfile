FROM golang:1.22

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY ./bin /app

WORKDIR /app

EXPOSE 8080
EXPOSE 9000
VOLUME /data

CMD ["./shop", "-conf", "/data/config.yaml"]
