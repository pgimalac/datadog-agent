from gcip import Artifacts, Cache, CacheKey, Need, Pipeline, Job, Rule

pipeline = (
    Pipeline()
    .initialize_image(
        "486234852809.dkr.ecr.us-east-1.amazonaws.com/ci/datadog-agent-buildimages/deb_x64$DATADOG_AGENT_BUILDIMAGES_SUFFIX:$DATADOG_AGENT_BUILDIMAGES"
    )
    .initialize_tags("arch:amd64")
    .initialize_variables(
        DATADOG_AGENT_BUILDIMAGES_SUFFIX="",
        DATADOG_AGENT_BUILDIMAGES="v23168861-fcc48431",
        S3_ARTIFACTS_URI="s3://dd-ci-artefacts-build-stable/$CI_PROJECT_NAME/$CI_PIPELINE_ID",
    )
    .prepend_scripts("source /root/.bashrc")
)
setup_agent_version = Job(
    stage="setup",
    script=[
        "inv -e agent.version --version-cached",
        "aws s3 cp $S3_CP_OPTIONS $CI_PROJECT_DIR/agent-version.cache $S3_ARTIFACTS_URI/agent-version.cache",
    ],
    name="agent_version",
)
go_deps = Job(
    stage="fetch",
    script=[
        "if [ -f modcache.tar.gz ]; then exit 0; fi",
        "inv -e deps --verbose",
        "cd $GOPATH/pkg/mod/ && tar czf $CI_PROJECT_DIR/modcache.tar.gz .",
    ],
    name="go_deps",
    rules=[
        Rule(if_statement="$CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH", variables={"POLICY": "pull-push"}),
        Rule(if_statement="$CI_COMMIT_BRANCH != $CI_DEFAULT_BRANCH", variables={"POLICY": "pull"}),
    ],
    needs=[setup_agent_version],
    artifacts=Artifacts("$CI_PROJECT_DIR/modcache.tar.gz", expire_in="1 day"),
    cache=Cache(
        cache_key=CacheKey(prefix="go_deps", files=[".**/go.mod]"]),
        paths=["modcache.tar.gz"],
    ),
)

pipeline.add_children(setup_agent_version, go_deps)
pipeline.write_yaml()
