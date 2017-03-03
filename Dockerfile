FROM golang:latest

# Set GOPATH/GOROOT environment variables
RUN mkdir -p /go
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

# go get all of the dependencies
RUN go get github.com/labstack/echo
RUN go get github.com/labstack/echo/middleware

# Set up app
ADD . /app
WORKDIR /app

EXPOSE 3000

CMD ["go", "run", "/app/main.go"]
