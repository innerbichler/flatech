VERSION 0.8



compile-scraper:
	FROM golang:1.21.0
	WORKDIR /build
	COPY ./scraper /build
	RUN go mod download
	RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scraper
	SAVE ARTIFACT /build/scraper AS LOCAL binaries/scraper

scraper:
	ARG tag="latest"
	FROM alpine:3.20.3
	COPY +compile-scraper/scraper /
	RUN apk add --no-cache firefox
	RUN apk add --no-cache geckodriver
	ENTRYPOINT ["/scraper"]
	SAVE IMAGE flatech-scraper:$tag

