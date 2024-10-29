# How to

1. Install:
Linux:
```sh
go install github.com/xaionaro-go/vodmover@latest
```
Windows:
Go to [https://github.com/xaionaro-go/vodmover/releases](https://github.com/xaionaro-go/vodmover/releases) and download the file.

2. Create config:
```yaml
obs:
  address: 127.0.0.1:4455
  password: myPasswordHere
move_vods:
- pattern_wildcard: "*2024*"
  destination: /tmp/2024/
```

3. Run:
```sh
vodmover --log-file path/to/log --config path/to/config
```
(or `vodmover.exe` if you are a Windows user)