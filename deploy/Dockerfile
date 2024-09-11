FROM golang:alpine AS builder

WORKDIR /jwt-auth-service

ADD ../go.mod .

COPY . .

RUN go build -o ./build/auth.exe ./cmd/main.go

FROM alpine

WORKDIR /jwt-auth-service

COPY ./schema /jwt-auth-service/schema
COPY ./configs/.env /jwt-auth-service/configs/.env
COPY --from=builder /jwt-auth-service/build /jwt-auth-service

CMD ["/jwt-auth-service/auth.exe"]