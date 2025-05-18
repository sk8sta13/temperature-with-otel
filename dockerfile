FROM golang:1.24 AS build
WORKDIR /app
ARG FOLDER=api_b
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api cmd/${FOLDER}/main.go

FROM scratch
WORKDIR /
COPY --from=build /api /api
ENTRYPOINT ["/api"]