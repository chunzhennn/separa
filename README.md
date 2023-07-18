# parascan

一个疯狂借鉴的扫描器

```text
Usage: separa scan

Flags:
  -h, --help                  Show context-sensitive help.

      --debug                 Debug mode to get more detail output.
  -t, --target=STRING         Target to scan, supports CIDR.
  -f, --target-file=STRING    Target file to scan, split each target line by
                              line with '\n'.
  -c, --config-file="config.yaml"
                              Config file to load, default is config.yaml in
                              current dir.
  -o, --output-file="output.json"
                              Output file to save, default is output.json in
                              current dir.
  -p, --port=STRING           Port to scan, default is TOP 1000. you can use
                              ',' to split or '-' to range, like '80,443,22' or
                              '1-65535'
  -d, --delay=5               Delay between each request
  -n, --top=1000              Top N ports to scan, default is 1000
```

## Dev

1. install go v1.20
2. `go mod tidy`

## Deploy

```sh
docker compose up
```

该 docker 环境默认直接扫描所给的 40 个 cidr。

## Example

```sh
./separa.exe -t 16.163.13.0
```