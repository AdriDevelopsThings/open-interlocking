FROM golang
WORKDIR /build
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM busybox

WORKDIR /dist
COPY --from=0 /build/main /dist/main


ENV OPEN_INTERLOCKING_HOST=0.0.0.0:80
EXPOSE 80
CMD ["/dist/main"]