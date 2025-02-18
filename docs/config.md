# Config

The configuration is the user supplied configuration that is used to drive the nuke process. The configuration is a YAML file that is loaded from the path specified by the --config flag.


## Simple Example

```yaml
presets:
  common:
    filters:
      Server:
        - type: dateOlderThanNow
          property: Created
          value: -24h
          invert: true
accounts:
  000000:
    presets:
      - common

resource-types:
  includes:
    - Server
```

This sample config removes **only** servers that are older than 24 hours.