
# Buliding Side Scanning

## Build the image

This command will build the agent and push it to the AWS registry:


    $ sh push.sh


The build requires a datadog.yml file in your repository.

The first run will take a while as it will checkout the agent's repository.

## Deploy the image

This command will update the ECS cluster running on the Sandbox AWS account with the provided container tag:


    $ sh deploy.sh 74f3bb9e29

Login to ECR, tagging, and pushing can be done with the `deploy.sh` script.

When an image is pushed, it's tagged with `latest` and the short commit id of HEAD.