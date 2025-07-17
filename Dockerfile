# Use the official Golang image as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Command to run the executable
CMD ["go","build","-x","-v","-o","aka","./cmd/aka"]
