FROM golang:alpine3.20

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY cmd/ /app/cmd/
RUN go build -o tt-api-gateway ./cmd/tt-api-gateway

ENTRYPOINT [ "./tt-api-gateway" ] 
# I dont feel like multistage building today :)