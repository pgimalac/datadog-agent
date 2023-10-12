set -e
set -x

IMAGE=sidescanning
REGISTRY=601427279990.dkr.ecr.us-west-2.amazonaws.com

docker build -t ${IMAGE} .

docker run -v $PWD:/headcp --rm --entrypoint cp ${IMAGE} /etc/datadog-agent/head  /headcp/

TAG=$(cat head)

docker tag ${IMAGE}:latest ${REGISTRY}/${IMAGE}:latest
docker tag ${IMAGE}:latest ${REGISTRY}/${IMAGE}:${TAG}

unset AWS_VAULT
aws-vault login dd-sandbox
aws-vault exec dd-sandbox -- aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin ${REGISTRY}


docker push ${REGISTRY}/${IMAGE}:latest
docker push ${REGISTRY}/${IMAGE}:${TAG}

