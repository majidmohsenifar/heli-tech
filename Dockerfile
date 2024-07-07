FROM golang:1.22.4-bullseye AS BUILD
ARG SRV
WORKDIR /app
#RUN go env -w GOPROXY=
RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY $SRV/go.mod $SRV/go.sum ./
RUN go mod download -x
COPY $SRV .
RUN go build -o main .
FROM base as builder

FROM debian:bullseye AS FINAL
RUN apt update && apt install -y ca-certificates 
WORKDIR /app
RUN groupadd -g 1001 -r heli && \
        useradd -u 1001 -r -s /bin/false -d /app -g heli heli && \
        chown -R heli:heli /app
USER heli:heli
COPY --from=BUILD --chown=heli:heli /app/main /app
CMD ["/app/main"]
