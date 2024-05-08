# Copyright © Koroutine Ltd

# //////////////////////////////////////////////////////////////////////////////////////////////////////////////////// #

# Building the binary of the App
FROM golang:latest AS base

# Set the working directory to /app
WORKDIR /usr/src/app

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# //////////////////////////////////////////////////////////////////////////////////////////////////////////////////// #

# Moving the binary to the 'final Image' to make it smaller
FROM golang:alpine as development

# Install gcc make and libc-dev to allow Makefile compilation
RUN apk update && apk add --no-cache gcc make libc-dev

# Set the working directory to /app
WORKDIR /usr/src/app

# Ensure staticcheck is executable and in the PATH
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

# Ensure go-critic is executable and in the PATH
RUN go install github.com/go-critic/go-critic/cmd/gocritic@latest

# Set the GOPATH environment variable
ENV $PATH:$(go env GOPATH)/bin

# Install szh shell:
RUN apk add --no-cache zsh git wget

# Install oh-my-zsh:
# Uses "Spaceship" theme with some customization. Uses some bundled plugins and installs some more from github
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.1.5/zsh-in-docker.sh)" -- \
    -t https://github.com/denysdovhan/spaceship-prompt \
    -a 'SPACESHIP_PROMPT_ADD_NEWLINE="false"' \
    -a 'SPACESHIP_PROMPT_SEPARATE_LINE="false"' \
    -p git \
    -p ssh-agent \
    -p https://github.com/zsh-users/zsh-autosuggestions \
    -p https://github.com/zsh-users/zsh-completions

# Copy across all the files
COPY . /usr/src/app

# //////////////////////////////////////////////////////////////////////////////////////////////////////////////////// #