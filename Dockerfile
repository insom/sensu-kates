FROM sstarcher/sensu

ENV RUNTIME_INSTALL=sensu-plugins-kubernetes
COPY kube.json /etc/sensu/check.d/
