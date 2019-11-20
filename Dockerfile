FROM scratch
WORKDIR $GOPATH/src/gin-blog
#COPY go.mod .
#COPY go.sum .
#RUN GO111MODULE=on go mod download
COPY . $GOPATH/src/gin-blog
#RUN GO111MODULE=on go build -o .
#RUN GO111MODULE=on GOPROXY=https://goproxy.io go build .
EXPOSE 8000
CMD ["./gin-blog"]