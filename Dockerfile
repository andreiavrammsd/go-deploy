FROM ubuntu:16.04

RUN apt-get update

RUN apt-get install git wget nano -y

RUN wget https://storage.googleapis.com/golang/go1.7.linux-amd64.tar.gz
RUN tar -xvf go1.7.linux-amd64.tar.gz
RUN mv go /usr/local

RUN echo "export GOROOT=/usr/local/go" >> ~/.profile
RUN echo "export GOPATH=/deploy" >> ~/.profile
RUN echo "export PATH=\$GOPATH/bin:\$GOROOT/bin:\$PATH" >> ~/.profile

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN mkdir -p ~/.ssh
RUN ssh-keyscan -H bitbucket.org >> ~/.ssh/known_hosts
RUN ssh-keyscan -H github.com >> ~/.ssh/known_hosts

ENTRYPOINT ["/entrypoint.sh"]
