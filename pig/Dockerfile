
FROM ubuntu:14.04
MAINTAINER chideat <chinaxiahaifeng@gmail.com>

RUN mkdir -p /opt/pig/logs
ADD conf /opt/pig/conf
ADD pig /opt/pig/

EXPOSE 7000

WORKDIR /opt/pig/

CMD ["./pig"]
