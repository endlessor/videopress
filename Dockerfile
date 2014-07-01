FROM ubuntu:12.04
MAINTAINER Andrew Allen <andrew@allenan.com>

# Mercurial
RUN echo 'deb http://ppa.launchpad.net/mercurial-ppa/releases/ubuntu precise main' > /etc/apt/sources.list.d/mercurial.list
RUN echo 'deb-src http://ppa.launchpad.net/mercurial-ppa/releases/ubuntu precise main' >> /etc/apt/sources.list.d/mercurial.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 323293EE

RUN apt-get update
RUN apt-get install -y wget git bzr mercurial

RUN wget -qO- http://golang.org/dl/go1.3.linux-amd64.tar.gz | tar -C /usr/local -xzf -

ENV PATH  /usr/local/go/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin
ENV GOPATH  /go
ENV GOROOT  /usr/local/go

RUN go get github.com/allenan/videopress

WORKDIR /go/src/github.com/allenan/videopress
ADD . /go/src/github.com/allenan/videopress

RUN go get
RUN go build main.go

EXPOSE 8000

ENTRYPOINT ./main
