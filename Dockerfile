FROM golang:1.17.0 AS build

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /go/src/eventpool-api
COPY . .
RUN go mod download
RUN go build -v 

FROM alpine:3.14.0
COPY --from=build /go/src/eventpool-api/eventpool-api /usr/local/lib/eventpool-api/eventpool-api
COPY --from=build /go/src/eventpool-api/configs/ /usr/local/lib/eventpool-api/configs/
COPY --from=build /go/src/eventpool-api/templates/ /usr/local/lib/eventpool-api/templates/
COPY --from=build /go/src/eventpool-api/assets/ /usr/local/lib/eventpool-api/assets/

WORKDIR /usr/local/lib/eventpool-api
ENTRYPOINT ["/usr/local/lib/eventpool-api/eventpool-api"]