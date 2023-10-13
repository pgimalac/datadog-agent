set -e
set -x

IMAGE=sidescanning
REGISTRY=601427279990.dkr.ecr.us-west-2.amazonaws.com
SANDBOX=dd-sandbox #sso-sandbox-account-admin

date > last_build_date
docker build -t ${IMAGE} .

docker run -v $PWD:/headcp --rm --entrypoint cp ${IMAGE} /etc/datadog-agent/head  /headcp/

TAG=$(cat head)

docker tag ${IMAGE}:latest ${REGISTRY}/${IMAGE}:latest
docker tag ${IMAGE}:latest ${REGISTRY}/${IMAGE}:${TAG}

unset AWS_VAULT
aws-vault login ${SANDBOX}
aws-vault exec ${SANDBOX} -- aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin ${REGISTRY}


docker push ${REGISTRY}/${IMAGE}:latest
docker push ${REGISTRY}/${IMAGE}:${TAG}

