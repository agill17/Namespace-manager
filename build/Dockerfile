FROM alpine

ENV HELM_VERSION="v2.10.0"

RUN apk add --no-cache ca-certificates \
    && wget -q https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz -O - | tar -xzO linux-amd64/helm > /usr/local/bin/helm \
    && chmod +x /usr/local/bin/helm


ENV OPERATOR=/usr/local/bin/namespace-manager \
    USER_UID=1001 \
    USER_NAME=namespace-manager

# install operator binary
COPY build/_output/bin/namespace-manager ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
