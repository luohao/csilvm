FROM centos:7

RUN yum install -y gcc-4.8.5 gcc-c++-4.8.5 make git lvm2-devel util-linux

ENV GOLANG_VERSION 1.9.2
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 de874549d9a8d8d8062be05808509c09a88a248e77ec14eb77453530829ac02b

RUN rm -rf /usr/local/go && \
      curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz && \
      echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - && \
      tar -C /usr/local -xzf golang.tar.gz && \
      rm -f golang.tar.gz

ENV GOPATH /go
ENV PATH /go/bin:$PATH
ENV PATH /usr/local/go/bin:$PATH

RUN mkdir -p /go/src/github.com/alecthomas && \
    cd /go/src/github.com/alecthomas && \
    git clone https://github.com/alecthomas/gometalinter.git --branch=v1.2.1 && \
    go install -v github.com/alecthomas/gometalinter && \
    gometalinter --install && \
    go get -u golang.org/x/tools/cmd/goimports && \
    mkdir -p /go/src/github.com/mesosphere/csilvm

WORKDIR /go/src/github.com/mesosphere/csilvm
