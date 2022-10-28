FROM golang:1.19
# ENV GOPROXY="https://proxy.golang.org,direct"
# ENV HTTP_PROXY="192.168.5.8:3128"
# ENV HTTPS_PROXY="192.168.5.8:3128"
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /server
CMD ["/server"]