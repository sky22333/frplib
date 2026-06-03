# frplib

`frplib` 是 frp 的 Android AAR 绑定库。它接收官方 TOML 配置字符串，用于启动、停止、重启 `frpc` / `frps`。

## 导入

把 AAR 放到 App：

```text
app/libs/frplib-universal.aar
```

在 Gradle 中引入：

```kotlin
dependencies {
    implementation(files("libs/frplib-universal.aar"))
}
```

生成的 Java 包名：

```text
io.github.sky22333.frplib
```

## 基础用法

App 启动后先设置私有临时目录：

```kotlin
val tempErr = Frplib.setTempDir(context.cacheDir.absolutePath)
if (tempErr.isNotEmpty()) {
    // 处理临时目录错误，不要继续启动 frp
}
```

启动客户端：

```kotlin
val err = Frplib.startClient(frpcToml)
if (err.isNotEmpty()) {
    // 处理错误
}

Frplib.reloadClient(newFrpcToml)
Frplib.stopClient()
```

启动服务端：

```kotlin
Frplib.startServer(frpsToml)
Frplib.reloadServer(newFrpsToml)
Frplib.stopServer()
```

返回值：

```text
""           成功
"CODE: ..." 错误
```

常见错误：

```text
ALREADY_RUNNING: ...
INVALID_TEMP_DIR: ...
INVALID_TOML: ...
START_FAILED: ...
STOP_FAILED: ...
RELOAD_FAILED: ...
```

## 多实例

```kotlin
Frplib.startClientWithID("client-a", frpcTomlA)
Frplib.startClientWithID("client-b", frpcTomlB)
Frplib.stopClientWithID("client-a")

Frplib.startServerWithID("server-a", frpsTomlA)
Frplib.reloadServerWithID("server-a", newFrpsTomlA)
Frplib.stopServerWithID("server-a")
```

辅助方法：

```kotlin
Frplib.isClientRunning()
Frplib.isClientRunningWithID("client-a")
Frplib.isServerRunning()
Frplib.isServerRunningWithID("server-a")

Frplib.stopAll()
Frplib.listInstances()
```

`listInstances()` 每行返回一个实例：

```text
type:id:state
type:id:state:lastError
```

## 日志

```kotlin
Frplib.setLogCallback(object : FrpLogCallback {
    override fun onLog(instanceID: String, type: String, level: String, message: String) {
        // type: client / server / frp
        // level: trace / debug / info / warn / error
    }
})
```

生命周期日志包含实例 ID。frp 内部日志的 `type` 为 `frp`，实例 ID 为空。

日志回调不保证在主线程。需要更新 UI 时，请切回 Android 主线程。

## 配置文件路径

传入的是 TOML 字符串。TOML 中如果引用证书、密钥、include 等文件，建议使用 App 私有目录下的绝对路径。

## Reload

`reloadClient` 和 `reloadServer` 使用安全重启：

```text
验证新 TOML -> 停止旧实例 -> 启动新实例
```

如果新 TOML 验证失败，旧实例会继续运行。

<details>
<summary>已知说明</summary>

`StopServer` / `stopServerWithID` 可能因上游 frps `Run()` 未返回而超时，即使端口已释放也不代表 frplib 已收到完整退出确认。

</details>
