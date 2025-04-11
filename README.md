# kubesystem

A project with examples about how to run a Kubernetes locally with different services deployed.

## Before getting started

Most of the commands in this tutorial, specially the ones to install different tools, are using MacOS example.

## Kubernetes

Install [king](https://kind.sigs.k8s.io/) to run your Kubernetes cluster.

```sh
brew install kind
```

Create cluster

```sh
kind create cluster --config kubernetes/config.yaml
```
