#!/bin/bash
env GOOS=linux go build -o main .
docker build -t wallace-hatch-backend .
docker tag wallace-hatch-backend:latest 145054867171.dkr.ecr.us-east-1.amazonaws.com/wallace-hatch-backend:latest
docker push 145054867171.dkr.ecr.us-east-1.amazonaws.com/wallace-hatch-backend:latest

export NAME=kops.wallacehatch.com
export KOPS_STATE_STORE=s3://wallace-hatch-kubernetes
kops export kubecfg ${NAME}

/usr/local/bin/kubectl set image deployment/client-deployment client-container=145054867171.dkr.ecr.us-east-1.amazonaws.com/wallace-hatch-backend:latest