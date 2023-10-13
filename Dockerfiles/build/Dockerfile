

FROM ubuntu:latest as build

RUN apt-get update && apt-get install -qy   \
        golang-1.20                         \
        python3                             \
        python3-pip                         \
        git                                 \
        clang                               \
        cmake                               \
        vim                                 \
        openssh-server

RUN pip install invoke

RUN git clone https://github.com/DataDog/datadog-agent.git

WORKDIR /datadog-agent

ENV PATH=$PATH:/usr/lib/go-1.20/bin/

RUN invoke install-tools

# This is just to trigger the rest of the build
# last_build_date is populated by push.sh
COPY last_build_date /

RUN git pull && git checkout jinroh/ebs-sbom
RUN invoke agent.build --build-exclude=systemd

RUN git rev-parse --short HEAD > head

FROM ubuntu:latest

RUN apt-get update && apt-get install -qy   \
        ca-certificates

RUN mkdir -p /etc/datadog-agent/conf.d/

COPY --from=build /datadog-agent/bin/agent/agent /bin/
COPY --from=build /datadog-agent/bin/agent/dist/conf.d /etc/datadog-agent/conf.d/
COPY --from=build /datadog-agent/dev/lib/libdatadog-agent-rtloader.so.0.1.0 /lib/libdatadog-agent-rtloader.so.1
COPY --from=build /datadog-agent/head /etc/datadog-agent/
COPY datadog.yaml /etc/datadog-agent/

CMD ["agent", "run"]
