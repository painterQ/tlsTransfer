## 备忘录
1. 配置、证书、数据默认在当前目录查找，如果没有在当前目录，则查找~/.tls_transfer，都没有则报错
2. 配置文件默认在~/.tls_transfer/conig.toml
3. 数据默认在~/.tls_transfer/data
4. 提供docker镜像
5. 支持使用server酱推送，支持使用钉钉机器人
6. TLS双向认证

```
go install  golang.org/x/tools/cmd/stringer@latest
go install github.com/davidrjenni/reftools/cmd/fillstruct@latest
```