FROM golang:1.15

RUN mkdir -p $GOPATH/src/github.com/ncostamagna/axul_contact
WORKDIR $GOPATH/src/github.com/ncostamagna/axul_contact

COPY . .
RUN ls

ARG DATABASE_HOST
ARG DATABASE_USER
ARG DATABASE_PASSWORD 
ARG DATABASE_NAME
ARG DATABASE_PORT
ARG DATABASE_DEBUG
ARG DATABASE_MIGRATE
ARG APP_PORT
ARG APP_URL

ENV DATABASE_HOST $DATABASE_HOST
ENV DATABASE_USER $DATABASE_USER
ENV DATABASE_PASSWORD $DATABASE_PASSWORD
ENV DATABASE_NAME $DATABASE_NAME
ENV DATABASE_PORT $DATABASE_PORT
ENV DATABASE_DEBUG $DATABASE_DEBUG
ENV DATABASE_MIGRATE $DATABASE_MIGRATE
ENV APP_PORT $APP_PORT
ENV APP_URL $APP_URL

RUN go get -d -v ./... 
RUN go install -v ./...
EXPOSE 8080

CMD ["axul_contact"]