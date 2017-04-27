FROM golang

ADD . /go/src/gymapp

#restore golang dependencies
WORKDIR /go/src/gymapp
RUN go get github.com/tools/godep && godep restore


RUN go get github.com/revel/cmd/revel
WORKDIR /go
ENTRYPOINT revel run gymapp dev 3000

EXPOSE 3000
