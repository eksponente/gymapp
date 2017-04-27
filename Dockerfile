FROM golang

ADD . /go/src/gymapp

#restore golang dependencies
WORKDIR /go/src/gymapp
RUN go get github.com/tools/godep && godep restore


RUN go get github.com/revel/cmd/revel
# RUN revel build gymapp dev
WORKDIR /go
ENTRYPOINT revel run gymapp dev 3000




#configure the Database -- supposedly done by amazon
# Add PostgreSQL's repository. It contains the most recent stable release
#     of PostgreSQL, ``9.3``.
# RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list

# Install ``python-software-properties``, ``software-properties-common`` and PostgreSQL 9.3
#  There are some warnings (in red) that show up during the build. You can hide
#  them by prefixing each apt-get statement with DEBIAN_FRONTEND=noninteractive
# RUN apt-get update && apt-get install -y postgresql postgresql-contrib
#
# USER postgres
#
# RUN    /etc/init.d/postgresql start &&\
#     psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
#     createdb -O docker docker


# ENTRYPOINT revel run gymapp dev 3000

# EXPOSE 3000


# RUN go install gymapp

# ENTRYPOINT /go/bin/gymapp

# FROM golang:1.4.1-onbuild
# RUN mkdir -p /go/src/gymapp
#
# RUN go get github.com/tools/godep
# RUN godep restore
#
# ENTRYPOINT revel run gymapp dev 9000
# EXPOSE 3000
