FROM golang:1.9.0

# for testing
RUN git config --global user.name "Foobert"
RUN git config --global user.email "foo@bar.com"

ARG GITHUB_TOKEN

WORKDIR /go/src/github.com/figome/semantic-changelog
RUN go get github.com/aktau/github-release

COPY ./kitt/kitt /deps/kitt
