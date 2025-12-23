FROM golang:1.25-alpine AS build

WORKDIR /app

RUN apk update
RUN apk add make

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN make build-binary app_name=bootstrap

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/bootstrap ./bootstrap

ENTRYPOINT ["./bootstrap"]
