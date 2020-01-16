# OAM Builder

OAM Builder generates OAM types and controller code scaffold based on k8s operator generator. It generates OAM (Workload, Trait,Scope) Operators automatically and lets developers focus on buisness logic.

Prerequisites:

*	[Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)

## Usage

Install:

```
make
export PATH=$PWD/bin:$PATH
```

Initialize domain:

```
kubebuilder init --domain test
```

Create Workload:

```
oambuilder workload --group demo --version v1 --kind DemoWorkload
```

Create Trait:

```
oambuilder trait --group demo --version v1 --kind DemoTrait
```

Create Exchange:

```
oambuilder exchange --group demo --version v1 --kind DemoExchange
```

