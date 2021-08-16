# golang base image
# use alpine for lightweight image output
FROM golang:alpine as builder

# install git
RUN apk update && apk add --no-cache git

# set current working directory
WORKDIR /app

# copy go mod and go sum
COPY go.mod go.sum ./

# download dependencies
RUN go mod download

# copy sources
COPY . .

# build go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# begin multi stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# copy pre-built binary files, .env and firebaseServiceAccount.json
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/firebaseServiceAccount.json .

# expose port
EXPOSE 8080

# run executable
CMD ["./main"]

