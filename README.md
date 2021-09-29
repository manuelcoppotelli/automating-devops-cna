Automating DevOps in a Cloud-Native way
=======================================

## Authors

Manuel Coppotelli [@manuel_coppo](https://twitter.com/manuel_coppo)

Lino Telera [@linotelera](https://twitter.com/linotelera)


## Setup

### Requirements:

```sh
brew install helm helmfile
helm plugin install https://github.com/databus23/helm-diff
```

### Deploy infrastructure

```sh
cd infrastructure

helmfile apply
kubectl apply -f tekton
kubectl apply -f flux.yaml

cd ..
```


## Configuration

### Configure GitOps

```sh
kubectl apply -f gitops
```


### Configure Tekton Pipeline

```sh
kubectl apply -f pipelines
```
