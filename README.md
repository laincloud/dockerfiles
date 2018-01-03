# dockerfiles

安装了 `s6` 的 Dockerfiles.

## s6 是什么？

[s6](https://skarnet.org/software/s6/index.html) 是为 Unix 设计的进程管理工具，
并且提供了 [s6-log](https://skarnet.org/software/s6/s6-log.html) 等日志收集
工具。

## 为什么要使用 `s6`？

在 Linux 系统里，PID 为 1 的进程是特殊的：

- PID 为 1 的进程需要在子进程退出后收割子进程
- PID 为 1 的进程收到 `SIGTERM` 或 `SIGINT` 的信号时不会被杀掉

假如容器里 PID 为 1 的进程没有正确收割子进程，会造成僵尸进程；假如 PID 为 1
的进程没有正确处理信号，会造成无法停止容器。而应用开发者一般不会，也不需要关心
这两件事。因此，我们需要为镜像设置一个默认的 PID 为 1 的进程来管理子进程和信号。
这样，当应用开发者不关心自己的进程 PID 为几时，由此 PID 为 1 的进程来管理
开发者的进程、传递信号，并在应用不响应信号时做适当的操作；当应用开发者
希望使用自己的 PID 为 1 的进程时，可以覆盖此默认设置。

另一方面，在使用 docker 的过程中，我们发现 docker 的日志收集有性能瓶颈：假如容器
将日志打到了标准输出，根据 `log-driver` 的设置，docker 会将容器的标准输出收集到
`syslog` 或 `json-file` 等；而这个过程是中心化的，所有容器的标准输出都会先经过
docker daemon，docker daemon 会成为收集过程的瓶颈。我们测试后得到以下结果：

log-driver | 日志收集速度
--- | ---
`syslog` | 14.9 MB/s
`json-file` | 37.9 MB/s

假如不是将日志打到标准输出，而是直接写入文件，速度则为 220 MB/s。可见 docker daemon
收集标准输出里的日志时性能较差。尤其严重的是，假如日志过多，会造成 docker daemon
卡死，从而影响宿主机上所有的容器。

`s6` 能同时解决上述 3 个问题。`s6-svscan` 和 `s6-supervise` 负责管理子进程，可以在
子进程退出后收割子进程，并正确的处理信号；`s6-log` 可以收集 stdout/stderr 并自动 rotate，
并且能达到 200 MB/s 的速度。所以我们使用 `s6`。

> `s6-log` 的性能测试参见 [centos/7.3.1611/README.md#性能测试](centos/7/README.md#性能测试)。

## `s6-overlay` 与 `s6` 的关系

[s6-overlay](https://github.com/just-containers/s6-overlay) 提供了一系列初始化脚本和工具，
方便了 `s6` 在 docker 镜像中的使用。

## 我们对 `s6-overlay` 的定制

我们定义了默认的 [s6 service](https://skarnet.org/software/s6/servicedir.html)，启动容器时：

- [app-init](centos/7.3.1611/rootfs/app-init) 将 Dockerfile 或命令行里传进来的 CMD 写进
  `/etc/services.log/app/run`，然后 `s6-supervise` 会自动启动 `/etc/services.log/app/run`
- `s6-supervise` 自动启动 [/etc/services.log/app/log/run](centos/7.3.1611/rootfs/etc/services.d/app/log/run)，
  收集 CMD 的 stdout/stderr，并写入 `/lain/logs/default/current`
- `s6-log` 在 `/lain/logs/default/current` 达到 268435455 Byte 时，会自动 rotate，并保留 3 份
  历史日志

CMD 退出时，`s6-supervise` 会自动执行 [/etc/services.log/app/finish](centos/7.3.1611/rootfs/etc/services.d/app/finish)，
然后容器也跟着退出。

另外，我们还将 `s6-overlay` 里的 `s6-log` 升级到了 `2.6.0.0`，原因请参见
[centos/7.3.1611/README.md#为什么有-rootfsusrbins6-log-文件](centos/7/README.md#为什么有-rootfsusrbins6-log-文件)

## 使用示例

```
# 使用 s6 管理 docker CMD 并收集 CMD 的 stdout/stderr
docker run -d -v /lain/logs/default:/lain/logs/default laincloud/centos:7.3.1611 ${CMD}
# 这时可以在宿主机上通过 tail -f /lain/logs/default/current 查看 ${CMD} 的标准输出

# 不使用 s6 管理 docker CMD
docker run -d --entrypoint ${CMD} laincloud/centos:7.3.1611
# 这时可以在宿主机上通过 docker logs -f ${CONTAINER_ID} 查看 ${CMD} 的标准输出
```

## dockerfiles 的构建

因为 Dockerfile 较多，如果每次更改都要触发全部构建的话速度将会很慢，所以我们编写
了 `dockerfiles` 命令行工具，只重新构建发生了变化的镜像以及依赖于变更镜像的镜
像。

### 编译

```
go get -u github.com/golang/dep/cmd/dep
dep ensure
go test ./...
go install
```

### 运行

```
dockerfiles build  # 构建受到 (origin/master, HEAD] 影响的镜像
dockerfiles pull  # 下载受到 (orgin/master, HEAD] 影响的镜像
dockerfiles push  # 推送受到 (origin/master, HEAD] 影响的镜像
dockerfiles retag --old-registry-host ${oldRegistryHost} --old-organization ${oldOrganization} --new-registry-host ${newRegistryHost} --new-organization ${newOrganization} # 为受到 (origin/master, HEAD] 影响的镜像重新打标签
```

> - dockerfiles 的运行依赖于：
>     - docker
>     - git
>     - make
> - dockerfiles 构建 ${a}/${b} 目录下的 Dockerfile 时：
>     - 如果 Dockerfile 首行有 `# TAGS ${tag1} ${tag2} ... ${tagN}`，则将其打标签为 `laincloud/${a}:${tag1}, laincloud/${a}:${tag2}，...，laincloud/${a}:${tagN}`
>     - 如果没有，则将其打标签为 `laincloud/${a}:${b}`
