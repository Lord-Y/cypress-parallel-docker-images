ARG CYPRESS_DOCKER_IMAGE_VERSION="7.2.0"
FROM cypress/included:${CYPRESS_DOCKER_IMAGE_VERSION}
ARG CYPRESS_PARALLEL_CLI="v0.0.1"

RUN echo CYPRESS_PARALLEL_CLI $CYPRESS_PARALLEL_CLI CYPRESS_DOCKER_IMAGE_VERSION $CYPRESS_DOCKER_IMAGE_VERSION
RUN wget https://github.com/Lord-Y/cypress-parallel-cli/releases/download/${CYPRESS_PARALLEL_CLI}/cypress-parallel-cli_linux_amd64.tar.gz
RUN tar xzf cypress-parallel-cli_linux_amd64.tar.gz
RUN rm -f cypress-parallel-cli_linux_amd64.tar.gz
RUN mv cypress-parallel-cli_linux_amd64 /usr/local/bin/cypress-parallel-cli
RUN chmod +x /usr/local/bin/cypress-parallel-cli

ENTRYPOINT ["/usr/local/bin/cypress-parallel-cli"]
