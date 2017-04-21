# Phusion's baseimage seems a good starting point
# See https://github.com/phusion/baseimage-docker
FROM phusion/baseimage:0.9.19

# Go
FROM golang:1.7

# Use baseimage-docker's init system.
CMD ["/sbin/my_init"]

# Add basic dependencies
RUN apt-get update && apt-get install -y \
		ca-certificates curl gcc libc6-dev make \
		bzr git mercurial \
		g++ \
		curl \
        zip unzip \
        python2.7 \
		--no-install-recommends

# NodeJS
RUN curl -sL https://deb.nodesource.com/setup_6.x | bash - && \
    apt-get install -y nodejs 

# Typescript
RUN npm install -g typescript 
RUN npm install -g ts-loader 

# Webpack
RUN npm install -g webpack

# Glide (glide.sh)
RUN curl https://glide.sh/get | bash

# Go bindata
RUN go get -u github.com/jteeuwen/go-bindata/...

# AWS CLI
#RUN mkdir -p /root/.aws
#COPY .aws/config /root/.aws
#COPY .aws/credentials /root/.aws
#RUN curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip" \
  #&& unzip awscli-bundle.zip \
  #&& ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws \
  #&& rm -rf ./awscli-bundle \
  #&& rm awscli-bundle.zip

# UPX
ADD https://github.com/upx/upx/releases/download/v3.92/upx-3.92-amd64_linux.tar.xz /usr/local
RUN apt-get update && apt-get install -y xz-utils && \
    xz -d -c /usr/local/upx-3.92-amd64_linux.tar.xz | tar -xOf - upx-3.92-amd64_linux/upx > /bin/upx && \
    chmod a+x /bin/upx

# Setup working directory
WORKDIR /mnt/ur
