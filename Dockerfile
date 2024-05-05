FROM golang:1.22

WORKDIR /home

COPY ./ /home

RUN cd /home && go build -o app ./cmd/app/main.go

CMD [ "/home/app" ]
