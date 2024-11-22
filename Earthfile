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
	COPY +compile-scraper/ /
	RUN apk add --no-cache firefox
	RUN apk add --no-cache geckodriver
	RUN mkdir /secrets
	# only do this because the path is currently hardcoded
	#RUN ln  /snap/bin/geckodriver /usr/bin/geckodriver
	ENTRYPOINT ["/scraper", "-file", "/secrets/.secrets"]
	SAVE IMAGE flatech/scraper:$tag

