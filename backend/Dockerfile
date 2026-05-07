FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o /out/jijin-backend ./cmd/server

FROM alpine:3.20
COPY --from=build /out/jijin-backend /usr/local/bin/jijin-backend
EXPOSE 8080
CMD ["jijin-backend"]
