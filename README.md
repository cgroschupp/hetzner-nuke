# hetzner-nuke

[![license](https://img.shields.io/github/license/cgroschupp/hetzner-nuke.svg)](https://github.com/cgroschupp/hetzner-nuke/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/cgroschupp/hetzner-nuke.svg)](https://github.com/cgroschupp/hetzner-nuke/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/cgroschupp/hetzner-nuke)](https://goreportcard.com/report/github.com/cgroschupp/hetzner-nuke)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/cgroschupp/hetzner-nuke/total)
![GitHub Downloads (all assets, latest release)](https://img.shields.io/github/downloads/cgroschupp/hetzner-nuke/latest/total)



## Overview

Remove all resources from an Hetzner account.

## Example

### Build locally
```sh
go install .
# Make sure the go path is inside the PATH
# export PATH=$(go env GOPATH)/bin:$PATH
touch config.yaml
hetzner-nuke run --hcloud-token <your-token>
```

### Run it with podman or docker
```sh
# Create dummy config
touch config.yaml
# Add hcloud token
export HCLOUD_TOKEN=<your-token>
# docker run -ti -v $(pwd)/config.yaml:/config.yaml -e HCLOUD_TOKEN ghcr.io/cgroschupp/hetzner-nuke:v0.1.0 run
podman run -ti -v $(pwd)/config.yaml:/config.yaml -e HCLOUD_TOKEN ghcr.io/cgroschupp/hetzner-nuke:v0.1.0 run
```