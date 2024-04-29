# system-server
[![](https://github.com/beclab/system-server/actions/workflows/build_main.yaml/badge.svg?branch=main)](https://github.com/beclab/system-server/actions/workflows/build_main.yaml)

## Description
As a part of system runtime frameworks, system-server provides a mechanism for security calls between apps.

## How to build
1. Install Custom Resources
```sh
kubectl apply -f config/crds
```

2. Generate the code
```sh
make update-codegen
```

3、Build binary
```sh
make system-server
```

or

```sh
make linux
```
For Linux version
