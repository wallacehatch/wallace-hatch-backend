env GOOS=linux go build -o docker/main .

cp -a email-templates/ docker/email-templates/


# docker-machine create --driver amazonec2 --amazonec2-instance-type t2.micro --engine-install-url=https://web.archive.org/web/20170623081500/https://get.docker.com wallace-hatch-backend


# docker stop $(docker ps -a -q)
# docker run -p 80:8090 -d --restart always demo-api-updated


export KUBE_CA_CERT=$(kubectl config view --flatten --output=json \
       | jq --raw-output '.clusters[2] .cluster ["certificate-authority-data"]')
export KUBE_ENDPOINT=$(kubectl config view --flatten --output=json \
       | jq --raw-output '.clusters[2] .cluster ["server"]')
export KUBE_ADMIN_CERT=$(kubectl config view --flatten --output=json \
       | jq --raw-output '.users[2] .user ["client-certificate-data"]')
export KUBE_ADMIN_KEY=$(kubectl config view --flatten --output=json \
       | jq --raw-output '.users[2] .user ["client-key-data"]')
export KUBE_USERNAME=$(kubectl config view --flatten --output=json \
       | jq --raw-output '.users[2] .user ["username"]')