set -e
set -x

IMAGE=sidescanning

docker run -v $PWD:/headcp --rm --entrypoint cp ${IMAGE} /etc/datadog-agent/head  /headcp/

TAG=$(cat head)

docker tag ${IMAGE}:latest 601427279990.dkr.ecr.us-west-2.amazonaws.com/${IMAGE}:latest 
docker tag ${IMAGE}:latest 601427279990.dkr.ecr.us-west-2.amazonaws.com/${IMAGE}:${TAG}

aws-vault login dd-sandbox
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 601427279990.dkr.ecr.us-west-2.amazonaws.com

docker push 601427279990.dkr.ecr.us-west-2.amazonaws.com/${IMAGE}:latest
docker push 601427279990.dkr.ecr.us-west-2.amazonaws.com/${IMAGE}:${TAG}

