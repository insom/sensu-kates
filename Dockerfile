FROM sstarcher/sensu

ENV RUNTIME_INSTALL=kubernetes
COPY kube.json /etc/sensu/conf.d/
