
Building the image:

    $ docker build . -t sidescanning

The build will fail unless you also add a datadog.yml file in your repository.

The first run will take a while as it will checkout the agent's repository.

Login to ECR, tagging, and pushing can be done with the `deploy.sh` script.

When an image is pushed, it's tagged with `latest` and the short commit id of HEAD.

