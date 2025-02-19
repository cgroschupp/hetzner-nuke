# hetzner-nuke

[![license](https://img.shields.io/github/license/cgroschupp/hetzner-nuke.svg)](https://github.com/cgroschupp/hetzner-nuke/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/cgroschupp/hetzner-nuke.svg)](https://github.com/cgroschupp/hetzner-nuke/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/cgroschupp/hetzner-nuke)](https://goreportcard.com/report/github.com/cgroschupp/hetzner-nuke)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/cgroschupp/hetzner-nuke/total)
![GitHub Downloads (all assets, latest release)](https://img.shields.io/github/downloads/cgroschupp/hetzner-nuke/latest/total)

## Overview

Remove all resources from an Hetzner account.

## Example

### Example config

config.yaml:
```yaml
accounts:
  000000: {}
```

Replace "000000" with your project id. To get your project id run `hetzner-nuke project-info --hcloud-token <your-token>`

### Build locally
```sh
go install .
# Make sure the go path is inside the PATH
# export PATH=$(go env GOPATH)/bin:$PATH
hetzner-nuke run --hcloud-token <your-token> -c config.yaml
```

### Run it with podman or docker
```sh
# Add hcloud token
export HCLOUD_TOKEN=<your-token>
# docker run -ti -v $(pwd)/config.yaml:/config.yaml -e HCLOUD_TOKEN ghcr.io/cgroschupp/hetzner-nuke:v0.1.0 run
podman run -ti -v $(pwd)/config.yaml:/config.yaml -e HCLOUD_TOKEN ghcr.io/cgroschupp/hetzner-nuke:v0.1.0 run
```