
#ADD kubectl
FROM bitnami/kubectl:latest as kubectl


FROM docker:dind
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache curl bash
ADD start.sh /
RUN chmod +x /start.sh

# For AMD64 / x86_64
RUN curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.24.0/kind-linux-amd64 \
    && chmod +x ./kind \
    && mv ./kind /usr/local/bin/kind
ADD kind-maker.sh /

COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/kubectl

# 初始化脚本，挂载覆盖
ADD ingress-controller-host.yaml /kind-init/ingress-controller-host.yaml

ENTRYPOINT [""]
CMD ["bash","/start.sh"]