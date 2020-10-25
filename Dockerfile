FROM golang:1.15

RUN mkdir -p $GOPATH/src/github.com/ncostamagna/axul_contact
WORKDIR $GOPATH/src/github.com/ncostamagna/axul_contact
COPY ./axul_contact .
RUN ls

RUN go get -d -v ./... 
RUN go install -v ./... 

EXPOSE 4000

CMD ["axul_contact"]