FROM amd64/alpine:latest
# FROM scratch

WORKDIR /app

COPY goapp ./
ENV TZ=Asia/Bangkok

ENTRYPOINT ["/app/goapp"]