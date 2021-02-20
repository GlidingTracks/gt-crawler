FROM golang:latest

COPY . /gt-crawler

WORKDIR /gt-crawler

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait

RUN chmod +x /wait
RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]

ENTRYPOINT /wait && CompileDaemon -include="*tmpl" -log-prefix=false -build="go build -o gt-crawler" -command="./gt-crawler"

