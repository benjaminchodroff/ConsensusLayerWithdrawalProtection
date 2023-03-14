FROM golang:1.19.4-bullseye

WORKDIR /opt/rescuer

COPY . .

RUN go mod download

RUN go build .
RUN go get clwp.rescuer

CMD ["/opt/rescuer/clwp.rescuer"]
