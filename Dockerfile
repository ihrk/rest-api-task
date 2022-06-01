FROM golang:1.18-alpine AS build_base

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go build ./cmd/rest-api-task/main.go


FROM alpine:latest 

RUN apk add --no-cache bash
RUN apk add --no-cache make
RUN apk add --no-cache git
COPY --from=build_base /app .

EXPOSE 8080

CMD ["./main"]
