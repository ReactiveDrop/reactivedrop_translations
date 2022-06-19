FROM golang:1.17

# install packages we use
RUN --mount=type=cache,target=/var/lib/apt apt update && apt -y install --no-install-recommends dos2unix

# setup git for line endings
RUN git config --global core.autocrlf true

# copy over entrypoint
COPY utils/docker/entrypoint.sh /

ENTRYPOINT [ "/entrypoint.sh" ]