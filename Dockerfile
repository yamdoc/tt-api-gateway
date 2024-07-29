FROM golang:alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY cmd/ /app/cmd/
RUN go build -ldflags '-extldflags "-static"' -tags netgo,osusergo -o tt-api-gateway ./cmd/tt-api-gateway


FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/tt-api-gateway /app/

ENTRYPOINT [ "./tt-api-gateway" ] 