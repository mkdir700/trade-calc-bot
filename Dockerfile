FROM golang:1.22.2 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GIN_MODE=release
    
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# 指定OS等，并go build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o capital_calculator_tgbot .

FROM alpine

WORKDIR /app

COPY --from=builder /app/capital_calculator_tgbot .

RUN touch .env

ENV GIN_MODE=release 

ENV ENV=production

CMD ["./capital_calculator_tgbot"]
