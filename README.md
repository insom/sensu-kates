# sensu-kates

Playing with Sensu on Kubernetes

---

Starting with [sstarcher][]'s basically canonical [upstream Docker packaging of
Sensu][upstream] (which is excellent), I converted the example Docker Compose
file to Kubernetes resources using [kompose][].

This spits out five services and five deployments, one of each for:

- redis
- sensu-server
- sensu-api
- sensu-client
- uchiwa

Now, I don't think that Sensu client needs a service so I've deleted that. One
reason that you _might_ want that would be if you have nodes reporting their
status via the check submission API, and you're running a shared client to
allow that. I would probably look at using a DaemonSet for that, though.

With a couple of small changes to Uchiwa's deployment (to override the
`SENSU_HOSTNAME` to `api` -- matching the DNS that Kubernetes will set up) and
to Sensu client (you must provide `CLIENT_ADDRESS` or it doesn't work), it all
works out of the box.

The steps that you need to go from this repository to a working basic Sensu set
up in your default namespace are:

```bash
# Deployments
kubectl apply -f redis-deployment.yaml
kubectl apply -f api-deployment.yaml
kubectl apply -f server-deployment.yaml
kubectl apply -f client-deployment-vanilla.yaml # See below
kubectl apply -f uchiwa-deployment.yaml

# Services
kubectl apply -f redis-service.yaml
kubectl apply -f api-service.yaml
kubectl apply -f server-service.yaml
kubectl apply -f uchiwa-service.yaml
```

Of course, you can't _see_ anything, so I've imported a plain
`nginx-ingress-controller` (YAML included for completeness) and exposed it with
an `Ingress` resource.

```bash
# Ingress Controller
kubectl apply -f nginx-ingress-controller.yaml
# ... now customer expose-uchi.yaml with your desired hostname etc. ...
kubectl apply -f expose-uchi.yaml
```

That will give you a full stack, with one client which isn't monitoring much of anything.

---

## Boring

It would be at least marginally useful if we could monitor some Kubernetes stats. Of course, [there's a Sensu plugin for that][sensu-plugin], and it's able to use the service account exposed to the container if you're running it from within Kubernetes. The command will be:

```bash
check-kube-pods-runtime.rb --in-cluster
```

We'll also need to make sure that the required gem is installed inside the Sensu client, first. `sstarcher` has thought of that though: if you supply `RUNTIME_INSTALL` to the `sensu` container it will install `build-essential` and some other dependencies and install the gem for you onto the embedded Ruby's `GEM_PATH`.

We'll also need to bundle a new config file. It's pretty small:

```json
{
 "checks": {
    "check_pods_runtime": {
      "command": "check-kube-pods-runtime.rb --in-cluster",
      "subscribers": [],
      "standalone": true,
      "interval": 60,
      "ttl": 600
    }
  }
}
```

The easiest bet was to create a new Dockerfile that just imports the original, tweaks the environment and copys the file into `/etc/sensu/check.d/`:

```
FROM sstarcher/sensu

ENV RUNTIME_INSTALL=kubernetes
COPY kube.json /etc/sensu/conf.d/
```

Build that Dockerfile and push it somewhere:

```bash
docker built -t gcr.io/insom-161401/test/sensu .
gcloud docker push gcr.io/insom-161401/test/sensu

# I ended up making a project for this on Google Container Registry
# 
# Because my Kubernetes install is in my house and I don't fully understand
# ACLs yet, I made the image public. Here's how, if you want to do that:
#
# gsutil defacl ch -u AllUsers:R gs://artifacts.insom-161401.appspot.com
# gsutil acl ch -r -u AllUsers:R gs://artifacts.insom-161401.appspot.com
# gsutil acl ch -u AllUsers:R gs://artifacts.insom-161401.appspot.com
#
# Just be aware that makes THE WHOLE BUCKET PUBLIC.
```

The version of the Sensu client deployment which uses _this_ container is `client-deployment.yaml`:

```bash
kubectl apply -f client-deployment.yaml
```

---

### Postscript

There's room for improvement here. I'm only using headless services because
they are what [kompose][] gave me. I'm pretty sure I would better off involving
`kube-proxy` and setting these services up with actual ports.

[sensu-plugin]: https://github.com/sensu-plugins/sensu-plugins-kubernetes/
[kompose]: https://github.com/kubernetes-incubator/kompose
[sstarcher]: https://github.com/sstarcher
[upstream]: https://hub.docker.com/r/sstarcher/sensu/
