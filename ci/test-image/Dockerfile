FROM golang:1.9.0

# for testing
RUN git config --global user.name "Foobert"
RUN git config --global user.email "foo@bar.com"

WORKDIR /go/src/github.com/figome/semantic-changelog
RUN go get github.com/golang/dep/cmd/dep

COPY ./kitt/kitt /deps/kitt
