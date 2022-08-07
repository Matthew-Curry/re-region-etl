FROM golang:latest

LABEL maintainer="Matthew Curry <matt.curry56@gmail.com>"

WORKDIR /app

# handle dependencies using go.mod and go.sum
COPY go.mod .
COPY go.sum .
RUN go mod download

# copy remaining source files
COPY . . 

# specify needed env vars. Select vars passed in as args during build
ARG RE_REGION_ETL_USER
ARG RE_REGION_ETL_PASSWORD
ARG RE_REGION_DB
ARG DB_PORT
ARG DB_HOST

ENV RE_REGION_ETL_USER $RE_REGION_ETL_USER
ENV RE_REGION_ETL_PASSWORD $RE_REGION_ETL_PASSWORD
ENV RE_REGION_DB $RE_REGION_DB
ENV DB_PORT $DB_PORT
ENV DB_HOST $DB_HOST

# build the app, cmd to run with help option
RUN go build

ENTRYPOINT ["./re-region-etl"]

CMD ["-h"]