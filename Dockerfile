# syntax=docker/dockerfile:1


# Getting ready image of Golang from the Docker libraries
FROM golang:1.22.6


# Creating directory inside the image that you're building, in order to instruct
# the Docker to use this directory as the default destination for all subsequent commands.
WORKDIR /app

# Usually all modules are already exist in the image of the Golang that we've already got,
# however if you have your own modules installed that are not in the standart modules pakege
# you should copy it to your WORKDIR

# First parameter is what to copy, second parameter where ("./" is quite sensative)
COPY go.mod go.sum ./

RUN go mod download


# Copying source code to the image
COPY . .

# Compiling the application
RUN go build -o rsshub ./cmd

# Telling Docker what command to run when the image is used to start the container
CMD ["./rsshub"]