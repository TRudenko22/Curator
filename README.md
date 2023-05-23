# curator
Curator (/ˈkyo͝orˌādər/) is an air-gapped infrastructure consumption analysis project for OpenShift Container Platform.

## Description
Operator Curator is an air-gapped infrastructure consumption analysis tool for the Red Hat OpenShift Container Platform. Curator retrieves infrastructure utilization for the OpenShift Platform using Operator koku-metrics and provides users the ability to query the infrastructure utilization based on time period, namespace, and infrastructure parameters.

Users can generate periodic standard and custom reports on infrastructure utilization, which are optionally delivered through automated emails. Curator also provides APIs to query the information utilization data that is stored in a database in the OpenShift cluster and it can also be used to feed data collected to any infrastructure billing system or business intelligence system. Additionally, Curator also provides administrators of the OpenShift cluster the option to back up their cluster infrastructure consumption data to S3-compatible storage.

## Getting Started
You'll need to have administrator access to an OpenShift v.4.5+ cluster to deploy Operator Curator. 
For more information on the prerequisites, please view the Requirements section. 
Once deployed, all the authorized users and systems will be able to view the infrastructure utilization of OpenShift.

**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
#### Pre-requisites
- Admin access to OpenShift v.4.5+ 
- [Koku Metrics Operator](https://dev.operatorhub.io/operator/koku-metrics-operator) installed on the cluster

1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/curator:tag
```
	
3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/curator:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)


