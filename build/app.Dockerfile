FROM golang:1.23.0-alpine3.20 as build
WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download && go mod verify
COPY cmd /src/cmd
COPY internal /src/internal
RUN go build -o /app/app /src/cmd/auth/auth.go

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/app /app/app
EXPOSE 8000
CMD ["/app/app"]

