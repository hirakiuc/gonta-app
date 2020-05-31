![workflow:Go](https://github.com/hirakiuc/gonta-app/workflows/Go/badge.svg)

# gonta-app

`gonta-app` is a slack app written in golang, as a Cloud Function.

# Config

| Name | Description |
|:-----|:-----------:|
| GO111MODULE | on |
| GCP\_PROJECT | `your gcp project` |

# HowToDev

```
$ make run
```

# HowToDeploy

1. Create `env.yaml` file from `sample.env.yml`.
2. `make deploy`
