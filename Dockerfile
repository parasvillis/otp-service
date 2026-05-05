# get the base image, aliased as builder
FROM golang:1.25-alpine AS builder
# defile working directory inside the container
WORKDIR /app
# copy go mod and sum files
COPY go.mod go.sum ./
# download dependencies
RUN go mod download
# copy the source code
COPY . .
# build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .
# RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest  && chmod +x /go/bin/migrate

# use a distroless image for the final container -> optimized for running Go applications
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
# only copy the compiled binary from the builder stage, no source code or build tools -> lightweight and secure
COPY --from=builder /app/main /main
COPY ./migrations /migrations
# COPY --from=builder /go/bin/migrate /migrate
# COPY --from=builder /app/migrations /migrations
# run the application as a non-root user for better security -> attacker will not have root acces to do some crazyy stuff
USER nonroot:nonroot
ENTRYPOINT [ "/main" ]