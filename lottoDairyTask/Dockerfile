FROM golang:1.18


WORKDIR /go/src/project

RUN go env GO111MODULE=on
RUN go mod init lotto
RUN go mod tidy
RUN go get github.com/jasonlvhit/gocron
RUN go get github.com/gocolly/colly
RUN go get github.com/PuerkitoBio/goquery
RUN go get gorm.io/gorm
RUN go get gorm.io/driver/mysql
RUN go install -v github.com/ramya-rao-a/go-outline@latest
RUN go install -v golang.org/x/tools/gopls@latest
ADD main.go /go/src/project/
CMD ["go", "run", "main.go"]