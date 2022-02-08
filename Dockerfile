# Start from golang base image
FROM golang:alpine

# Add Maintainer info
LABEL maintainer="Alankrit Joshi <alankritjoshi@gmail.com>"

# Whois.
RUN apk add whois

## Install git.op
## Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Create folder app
RUN mkdir /app

# Copy the source from the current directory to the working Directory inside the container
ADD . /app

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 go build -v -o main .

# Expose port 3000 to the outside world
EXPOSE 3000

#Command to run the executable
CMD ["/app/main"]
