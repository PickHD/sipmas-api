FROM golang:1.16

WORKDIR /home/sipmas_api
COPY . .

RUN go get -u github.com/cosmtrek/air
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["server"]