FROM golang:1.16
ENV APP_HOME go/src/github.com/syned13/ticket-support-back

RUN mkdir -p $APP_HOME
ADD . $APP_HOME
RUN mkdir build
WORKDIR $APP_HOME

# RUN go mod download
RUN go build -o build/main ./cmd/main.go

RUN chmod +x ./build/main

CMD ["./build/main"]