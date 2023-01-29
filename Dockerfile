FROM golang:1.19

RUN mkdir -p $GOPATH/src/github.com/ncostamagna/axul_contact
WORKDIR $GOPATH/src/github.com/ncostamagna/axul_contact

COPY . .
RUN ls

RUN go get -d -v ./... 
RUN go install -v ./...
EXPOSE 8080

CMD ["axul_contact"]