# Builder

ARG GITHUB_PATH=github.com/hablof/logistic-package-api-facade

FROM golang:1.19-alpine AS builder

WORKDIR /home/${GITHUB_PATH}

RUN apk add --update make git protoc protobuf protobuf-dev curl

COPY . .
RUN make build

# facade

FROM alpine:latest as facade
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/facade .
# COPY --from=builder /home/${GITHUB_PATH}/config.yml .

RUN chown root:root facade

CMD ["./facade"]
