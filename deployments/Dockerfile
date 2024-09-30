FROM golang:1-alpine AS builder

ARG APP_Name
ARG APP_Version
ARG BUILD_Region
ARG GO_ldflags

LABEL stage=${APP_Name}_builder

RUN if [ "$BUILD_Region" = "cn" ]; then \
    echo "==> BUILD_Region: $BUILD_Region"; \
    cp /etc/apk/repositories /etc/apk/repositories.bk; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
    go env -w GOPROXY="https://goproxy.cn,direct"; \
  fi

RUN apk --no-cache update && \
  apk --no-cache upgrade

WORKDIR /app

ADD ./main.go ./go.mod ./go.sum ./project.yaml ./
COPY ./internal ./internal
COPY ./migrations ./migrations
COPY ./pkg ./pkg
# COPY ./proto ./proto
COPY ./bin ./bin

COPY ./cache.local ./cache.local

RUN if [ -d ./cache.local/vendor ]; then \
    echo "==> Found cache.local/vendor"; \
    mv ./cache.local/vendor vendor; \
  fi

RUN if [ -d ./bin/swagger-go/cache.local/vendor ]; then \
    echo "==> Found ./bin/swagger-go/cache.local/vendor"; \
    mv ./bin/swagger-go/cache.local/vendor bin/swagger-go/vendor; \
  fi

# in alpine, date doesn't parse %:z
# mkdir -p target && go build -o target/main -ldflags="-w -s ${GO_ldflags}" main.go
RUN mkdir -p target

RUN go build -trimpath -ldflags="-w -s ${GO_ldflags}" -o target/main main.go

RUN cd bin/swagger-go && \
  go build -trimpath -ldflags="-w -s ${GO_ldflags}" -o ../../target/swagger main.go

####
FROM alpine:3

ARG APP_Name
ARG APP_Version
ARG BUILD_Time
ARG BUILD_Region

ENV TZ=${TZ:-Asia/Shanghai}
ENV APP_Name=${APP_Name}
ENV BUILD_Time=${BUILD_Time}
ENV BUILD_Region=${BUILD_Region}

RUN if [ "$BUILD_Region" = "cn" ]; then \
    echo "==> BUILD_Region: $BUILD_Region"; \
    cp /etc/apk/repositories /etc/apk/repositories.bk; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
  fi

RUN apk --no-cache update && \
  apk --no-cache upgrade && \
  apk --no-cache add tzdata
  # curl gcompat

RUN mkdir -p /app/target
WORKDIR /app

COPY --from=builder /app/target/ /app/target/
COPY --from=builder /app/cache.local/build.yaml /app/
# COPY ./migrations ./migrations

# http server
EXPOSE 9011
# internal server
EXPOSE 9015
# grpc server
EXPOSE 9016
# swagger server
EXPOSE 9017


CMD [./target/main, api, -release, -config=configs/prod.yaml, -http.addr=:9011, --grpc.addr=:9015, --internal.addr=:9016]
# CMD [./target/main, crons, -config=configs/crons.yaml]
# CMD [./target/main, swagger, -config=configs/swagger.yaml, -http.addr=:9017]
