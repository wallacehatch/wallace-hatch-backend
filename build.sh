env GOOS=linux go build -o main .


docker-machine create --driver amazonec2 --amazonec2-instance-type t2.medium --engine-install-url=https://web.archive.org/web/20170623081500/https://get.docker.com name


docker stop $(docker ps -a -q)
docker run -p 80:8090 -d --restart always demo-api-updated