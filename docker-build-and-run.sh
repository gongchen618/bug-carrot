#!/bin/sh

docker stop bug-carrot
docker rm bug-carrot
docker rmi gongchen0618/bug-carrot:carrot0.0.1
docker rmi bug-carrot

docker build -t bug-carrot .
docker tag bug-carrot gongchen0618/bug-carrot:carrot0.0.1
docker images|grep none|awk '{print $3}'|xargs docker rmi
docker run -d --network host --name bug-carrot bug-carrot
docker logs bug-carrot