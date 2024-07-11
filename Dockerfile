ARG DOCKER_ENTRYPOINT

# Build stage
FROM golang:1.22-bookworm as builder
ARG DOCKER_ENTRYPOINT

RUN apt-get update && \
    apt-get install -y --no-install-recommends make ca-certificates && \
    update-ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /builder
ADD ./ ./

RUN make clean && make


# Final stage
FROM debian:bookworm-slim
ARG DOCKER_ENTRYPOINT

RUN groupadd --gid=999 --system app && \
    useradd --uid=999 --no-log-init --create-home --system --gid app app && \
    mkdir -p /data && \
    chown -R app:app /data && \
    apt-get update && \
    apt-get install -y --no-install-recommends gosu ca-certificates && \
    update-ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /data
ENV DATA_DIR /data

COPY --from=builder /builder/build/* /usr/local/bin/.
COPY --from=builder /builder/docker/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN ln -s /usr/local/bin/${DOCKER_ENTRYPOINT} /usr/local/bin/start
RUN chmod ugo+x /usr/local/bin/${DOCKER_ENTRYPOINT}
RUN chmod ugo+x /usr/local/bin/docker-entrypoint.sh
RUN chmod ugo+x /usr/local/bin/start

ENTRYPOINT ["docker-entrypoint.sh"]
