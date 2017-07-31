# centos:7.3.1611

## 性能测试

`test.log` 体积 | 直接用文件收集日志的速度 | 用 `s6-log` 收集日志的速度
--- | --- | ---
80 MB | 94.5 MB/s | 209 MB/s
160 MB | 91.4 MB/s | 212 MB/s
320 MB | 90.5 MB/s | 103 MB/s
640 MB | 234 MB/s | 99.6 MB/s
1280 MB | 225 MB/s | 95.7 MB/s
2560 MB | 212 MB/s | 91.8 MB/s

> 测试直接用文件收集日志的速度的指令:
>
> ```
> docker run -it --entrypoint /usr/bin/dd -v /root/tmp:/root/tmp laincloud/centos:7.3.1611 if=/root/tmp/test.log of=/root/tmp/test.out
> ```
>
> 测试用 `s6-log` 收集日志的速度的指令:
>
> ```
> docker run -it -v /root/tmp:/root/tmp -v /lain/logs:/lain/logs -e 'S6_LOGGING_SCRIPT=n3 s268435455' laincloud/centos:7.3.1611 dd if=/root/tmp/test.log
> ```
>
> `test.log` 是从生产环境里复制的日志.

当 `test.log` 的体积较小的时候，`dd` 直接写入文件的速度较慢是因为此时 `dd` 的潜力
还没有充分发挥出来；当 `test.log` 体积较大的时候，`dd` 直接写入文件的速度就可以达到
200 MB/s 了。

与此相反，当 `test.log` 的体积较小的时候，用 `s6-log` 收集日志的速度能
达到 200 MB/s，说明 `s6-log` 与直接写入文件相比性能几乎没有损失；当 `test.log` 的
体积大于 256 MB 的时候，日志文件发生了 rotate，导致速度下降到 100 MB/s 左右；继续增大
`test.log` 的体积，`s6-log` 的速度基本稳定在 90 MB/s 左右。鉴于只有日志文件的体积到达
256 MB 的时候才会触发 rotate，即 rotate 操作不会太频繁，`s6-log` 能满足生产环境的
收集日志需求。

## FAQ

### 为什么有 `rootfs/usr/bin/s6-log` 文件?

因为 `s6-overlay` 中的 `s6-log` 并不是最新版，日志文件最大只能为 16 MB，会频繁触发
rotate，降低收集日志的速度。而最新（2017-07-28 22:54:25）版 `s6-log`(2.6.0.0) 的文件
最大体积可以达到 268435455 Byte，有利于减少 rotate 次数，从而有利于提高收集日志的
速度。`rootfs/usr/bin/s6-log` 即为 `s6-log-2.6.0.0`。
