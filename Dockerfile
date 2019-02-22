FROM golang:1.11
#ARG DB_HOST
#ARG DB_USER
#ARG DB_NAME
#ARG DB_PASSWORD

#RUN mkdir -p /go/src/github.com/efkbook/blog-sample
#RUN mkdir -p /go/src/ginsample
#COPY src/ /go/src/
WORKDIR /go/src/ginsample

#RUN go get -v ./...
#RUN go install -v ./...
COPY gin.yml /go/bin/

#ENV DB_HOST ${DB_HOST}
#ENV DB_USER ${DB_USER}
#ENV DB_NAME ${DB_NAME}
#ENV DB_PASSWORD ${DB_PASSWORD}

CMD ["go", "run", "ginsample"]
EXPOSE 1323
