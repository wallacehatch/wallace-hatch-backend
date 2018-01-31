env GOOS=linux go build -o docker/main .


docker-machine create --driver amazonec2 --amazonec2-instance-type t2.micro --engine-install-url=https://web.archive.org/web/20170623081500/https://get.docker.com wallace-hatch-backend


docker stop $(docker ps -a -q)
docker run -p 80:8090 -d --restart always demo-api-updated