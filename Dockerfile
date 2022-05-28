FROM golang:1.14.3 as builder
WORKDIR /build/
COPY . /build/
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct
RUN make linux

FROM scratch
WORKDIR /project
COPY --from=builder /build/bin /project/bin
COPY ./ca-certificates.crt /etc/ssl/certs/