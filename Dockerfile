# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.2
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build

WORKDIR /src

# Install Java for PKL
# im running pkl with java to make it work on aarch64
RUN apt-get update && \
    apt-get install -y default-jdk && \
    apt-get install -y ant && \
    apt-get clean;
    
# Fix certificate issues
RUN apt-get update && \
    apt-get install ca-certificates-java && \
    apt-get clean && \
    update-ca-certificates -f;
    
# Setup JAVA_HOME -- useful for docker commandline
ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64/
RUN export JAVA_HOME

ADD --chmod=111 "https://repo1.maven.org/maven2/org/pkl-lang/pkl-cli-java/0.26.3/pkl-cli-java-0.26.3.jar" /bin/pkl

RUN go install "github.com/air-verse/air@latest"
RUN go install "github.com/pressly/goose/v3/cmd/goose@latest"
# RUN go install "github.com/a-h/templ/cmd/templ@latest"

COPY go.mod go.sum ./

RUN go mod download && go mod verify

# COPY . ./

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/manage ./cmd/manage/

# RUN --mount=type=cache,target=/go/pkg/mod/ \
#     --mount=type=bind,target=. \
#     CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd/server/

CMD ["air"]
