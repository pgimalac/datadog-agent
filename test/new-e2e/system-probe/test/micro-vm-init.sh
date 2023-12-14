#!/bin/bash
set -eEuxo pipefail

retry_count=$1
pkgs_run_config_file=$2
docker_dir=/kmt-docker

# Add provisioning steps here !
## Start docker
systemctl start docker
## Load docker images
if [[ -d "${docker_dir}" ]]; then
  find "${docker_dir}" -maxdepth 1 -type f -exec docker load -i {} \;
fi
# VM provisioning end !

apt install -y strace sysstat || true
yum install -y strace sysstat || true

# Start tests
code=0
/test-runner -retry "${retry_count}" -packages-run-config "${pkgs_run_config_file}" || code=$?

cp /job_env.txt /ci-visibility/junit/
tar -C /ci-visibility/testjson -czvf /ci-visibility/testjson.tar.gz .
tar -C /ci-visibility/junit -czvf /ci-visibility/junit.tar.gz .

exit ${code}
