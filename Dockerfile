FROM golang:1.18


WORKDIR /go/src/project

RUN go env GO111MODULE=on
RUN go mod init lotto
RUN go mod tidy
RUN go get github.com/jasonlvhit/gocron
RUN go get github.com/gocolly/colly
RUN go get github.com/PuerkitoBio/goquery
ADD main.go /go/src/project/
CMD ["go", "run", "main.go"]