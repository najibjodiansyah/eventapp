FROM golang:1.17

RUN mkdir /app

WORKDIR /app

COPY ./ /app

RUN cd ./app && go build -o app

CMD cd ./app && ./app