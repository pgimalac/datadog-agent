#!/bin/sh

set -e 
set -x

if [ $# -ne 1 ]; then
    echo "Usage: $0 tag"
    echo ""
    echo "Will deploy the side scanning container with tag short_git_commit_id"
    exit 1
fi;

TAG=$1

CLUSTER=SideScanning
REGION=us-west-2
SANDBOX=dd-sandbox #sso-sandbox-account-admin

aws-vault exec ${SANDBOX} -- aws --region=${REGION} ecr describe-images --repository-name=sidescanning --image-ids '[{"imageTag":"'${TAG}'"}]'

TMP=$(mktemp)
TMP2=$(mktemp)
TMP3=$(mktemp)
TMP4=$(mktemp)
TMP5=$(mktemp)
TMP6=$(mktemp)

aws-vault exec ${SANDBOX} -- aws --region ${REGION} ecs describe-task-definition --task-definition SideScanner > ${TMP}

cat ${TMP} | jq '.taskDefinition | {
  containerDefinitions:    .containerDefinitions,
  cpu:                     .cpu,
  executionRoleArn:        .executionRoleArn,
  family:                  .family,
  memory:                  .memory,
  networkMode:             .networkMode,
  placementConstraints:    .placementConstraints,
  requiresCompatibilities: .requiresCompatibilities,
  runtimePlatform:         .runtimePlatform,
  taskRoleArn:             .taskRoleArn,
  volumes:                 .volumes,
}' > ${TMP2}

IMAGE=$(cat ${TMP2}|jq -r '.containerDefinitions[0].image')
NEW_IMAGE=$(echo ${IMAGE} | sed -e "s/:[^:]*$/:${TAG}/")

cat ${TMP2} |jq '.containerDefinitions[0].image="'${NEW_IMAGE}'"' > ${TMP3}

aws-vault exec ${SANDBOX} -- aws --region ${REGION} ecs register-task-definition --cli-input-json file://${TMP3} > ${TMP4}

TASK_DEFINITION_ARN=$(cat ${TMP4} | jq -r .taskDefinition.taskDefinitionArn)

aws-vault exec ${SANDBOX} -- aws --region ${REGION} ecs run-task --task-definition ${TASK_DEFINITION_ARN}  --network-configuration 'awsvpcConfiguration={subnets=[subnet-19e8f831,subnet-461bbf0f,subnet-58d35e3f],securityGroups=[sg-69f1aa11]}' --capacity-provider-strategy capacityProvider=FARGATE --cluster SideScanning > ${TMP5}

TASK_ARN=$(cat ${TMP5} | jq -r .tasks[].taskArn)

aws-vault exec ${SANDBOX} -- aws --region ${REGION} ecs list-tasks --cluster ${CLUSTER} > ${TMP6}

for task in $(cat ${TMP6} | jq -r '.taskArns[]' | grep -v ${TASK_ARN}); do

    aws-vault exec ${SANDBOX} -- aws --region ${REGION} ecs stop-task --cluster ${CLUSTER} --task ${task}

done

rm ${TMP}  ${TMP2} ${TMP3}
rm ${TMP4} ${TMP5} ${TMP6}

