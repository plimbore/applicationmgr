# applicationmgr


### Bootstrap

Setup local development environment as instructed in [./docs/dev/local-env-setup.md](./docs/dev/local-env-setup.md) doc.


- Bootstrap application code

```sh
kubebuilder init --domain applicationmgr.io --repo github.com/plimbore/applicationmgr
```

This shall create following directory structure

```
aplicationmgr
  |+ .devcontainer
  |+ .github
  |+ cmd
  |+ config
  |+ hack
  |+ test
  |- .dockerignore
  |- .gitignore
  |- Dockerfile
  |- go.mod
  |- go.sum
  |- Makefile
  |- PROJECT
  |- README.md
```


Create api

```sh
# almc: Application Lifecycle Manager Controller
kubebuilder create api --group almc --version v1 --kind Application
```

This will create following files/directories

```
aplicationmgr
  |+ api
  |+ bin
  |+ internal
  |- .golangcdi.yml
```

Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:

```sh
make manifests
```

- Create basic helm chart

```sh
helm create applicationmgr
# Renamed the parent directory to helm-chart for easy understanding
mv applicationmgr helm-chart
```

Add file [./helm-chart/values-local.yaml](./helm-chart/values-local.yaml) for local deployment

- Add examples

Create [./examples](./examples/) directory

