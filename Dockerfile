FROM golang:1.19

WORKDIR $GOPATH/bin

COPY main .

EXPOSE 8081

CMD ["main"]