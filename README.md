# applicationmgr

Application Lifecycle Manager Controller for kubernetes

NOTICE: Helm chart is work in progress

## TODO

- [ ] Helm chart
- [ ] ingress https support

## Description

A Kubernetes controller that manages the lifecycle of applications defined by a new Custom Resource Definition (CRD) called Application.

When an Application resource is created, the controller should automatically provision and manage all necessary Kubernetes components for that application, including a Deployment, a Service, and an Ingress.

When the Application resource is deleted, the controller should clean up all associated Kubernetes resources it created.

The statuses are updated properly in the CR and required default settings are patched if required.

## Getting Started

### Prerequisites

- go version v1.23.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
- Refer [./docs/dev/local-env-setup.md](./docs/dev/local-env-setup.md) & [./HISTORY.md](./HISTORY.md) for more info

### To Deploy on the cluster

Note: You can refer [./docs/assets/asciinema/ cast files for build steps](./docs/assets/asciinema/).

- Build & test (For windows users: use WSL)

    ```sh
    asciinema play docs/assets/asciinema/dev-build-test.cast
    ```

- Apply Application CRD

    ```sh
    asciinema play docs/assets/asciinema/dev-apply.cast
    ```

- Cleanup

    ```sh
    asciinema play docs/assets/asciinema/dev-cleanup.cast
    ```


**Build and push your image to the location specified by `IMG`:**

```sh
# To build binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager cmd/main.go

# To build locally
# Note: If IMG is not provided, default value would be controller:latest
make docker-build IMG=applicationmgr:latest
# To build and push
make docker-build docker-push IMG=<some-registry>/applicationmgr:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/applicationmgr:tag
# To deploy local image
make deploy IMG=applicationmgr:latest
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample or examples directory:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall

**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

    ```sh
    make build-installer IMG=<some-registry>/applicationmgr:tag
    ```

    **NOTE:** The makefile target mentioned above generates an 'install.yaml'
    file in the dist directory. This file contains all the resources built
    with Kustomize, which are necessary to install this project without its
    dependencies.

2. Using the installer

    Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
    the project, i.e.:

    ```sh
    kubectl apply -f https://raw.githubusercontent.com/<org>/applicationmgr/<tag or branch>/dist/install.yaml
    ```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

    ```sh
    kubebuilder edit --plugins=helm/v1-alpha
    ```

2. See that a chart was generated under 'dist/chart', and users can obtain this solution from there.

    **NOTE:** If you change the project, you need to update the Helm Chart
    using the same command above to sync the latest changes. Furthermore,
    if you create webhooks, you need to use the above command with
    the '--force' flag and manually ensure that any custom configuration
    previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
    is manually re-applied afterwards.

## Contributing

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
