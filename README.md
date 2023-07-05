# parascan

一个疯狂抄袭的扫描器

```text
Usage: separa

A simple scanner for Web Security

Flags:
  -h, --help                  Show context-sensitive help.
      --debug                 Debug mode to get more detail output.
  -t, --target=STRING         Target to scan, supports CIDR.
  -f, --target-file=STRING    Target file to scan, split each target line by
                              line with '\n'.
  -c, --config-file="config.yaml"
                              Config file to load, default is config.yaml in
                              current dir.
  -o, --output-file="output.txt"
                              Output file to save, default is output.txt in
                              current dir.
  -p, --port=STRING           Port to scan, default is TOP 1000. you can use
                              ',' to split or '-' to range, like '80,443,22' or
                              '1-65535'
```

## Dev

1. install go v1.20
2. `go mod download`
3. `go build`
4. `./separa.exe --help`

## Example

```sh
./separa.exe -t 16.163.13.0
```