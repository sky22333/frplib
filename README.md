# frplib

`frplib` 是官方 frp 的极简 Android AAR 核心绑定库。

它只负责在 Android 进程内启动、停止、重启 `frpc` / `frps`。配置直接使用官方 TOML，不包含 ForegroundService、通知栏、UI、数据库、权限申请、开机自启、保活、VPN Service、配置编辑器、订阅系统、账号系统或 JitPack 发布。

## AAR

从 GitHub Releases 下载：

```text
frplib-arm64-v8a.aar
frplib-armeabi-v7a.aar
frplib-x86_64.aar
frplib-universal.aar
```

普通业务 App 推荐使用：

```text
frplib-universal.aar
```

支持 ABI：

```text
arm64-v8a
armeabi-v7a
x86_64
universal
```

不构建 `x86`。

## 导入

把 AAR 放到业务 App：

```text
app/libs/frplib-universal.aar
```

Gradle 引入：

```kotlin
dependencies {
    implementation(files("libs/frplib-universal.aar"))
}
```

生成的 Java 包名：

```text
io.github.sky22333.frplib
```

## API

所有启动、停止、重载方法都返回字符串：

```text
""           表示成功
"CODE: ..." 表示错误
```

默认单实例：

```kotlin
Frplib.StartClient(frpcToml)
Frplib.StopClient()
Frplib.ReloadClient(frpcToml)
Frplib.IsClientRunning()

Frplib.StartServer(frpsToml)
Frplib.StopServer()
Frplib.ReloadServer(frpsToml)
Frplib.IsServerRunning()
```

多实例：

```kotlin
Frplib.StartClientWithID("client-a", frpcTomlA)
Frplib.StopClientWithID("client-a")
Frplib.ReloadClientWithID("client-a", frpcTomlB)
Frplib.IsClientRunningWithID("client-a")

Frplib.StartServerWithID("server-a", frpsTomlA)
Frplib.StopServerWithID("server-a")
Frplib.ReloadServerWithID("server-a", frpsTomlB)
Frplib.IsServerRunningWithID("server-a")
```

停止全部实例：

```kotlin
Frplib.StopAll()
```

查看实例：

```kotlin
Frplib.ListInstances()
```

返回格式为每行一个实例：

```text
type:id:state
type:id:state:lastError
```

## 行为

重复启动同一个 `client/server + id` 会返回：

```text
ALREADY_RUNNING: ...
```

重复停止不存在、已停止、停止中或失败的实例是安全 no-op，返回空字符串。

错误 TOML 会返回明确错误，例如：

```text
INVALID_TOML: parse frpc TOML failed: ...
```

实例状态由 frp 服务的 `Run(ctx)` 真实退出驱动。`Stop` 会发出取消信号，必要时关闭底层资源，并等待真实退出或超时。

## Reload

`ReloadClient` / `ReloadServer` 使用 safe restart，不承诺热更新：

```text
先验证新 TOML。
验证失败：不停止旧实例。
验证成功：停止旧实例，再用新 TOML 启动新实例。
```

如果新 TOML 能解析，但新实例因为端口占用、权限、网络或其他运行时原因启动失败，库不承诺自动恢复旧实例。

## 日志

设置日志回调：

```kotlin
Frplib.SetLogCallback(object : FrpLogCallback {
    override fun onLog(instanceID: String, type: String, level: String, message: String) {
        // type: client / server / frp
        // level: trace / debug / info / warn / error
    }
})
```

日志分两类：

```text
client/server 生命周期日志：包含精确 instanceID 和 type。
frp 官方内部日志：type 为 frp，instanceID 为空，不伪造实例归属。
```

## 发布

在 GitHub Actions 手动运行发布工作流，输入上游 frp tag，例如：

```text
v0.69.1
```

CI 使用输入的 tag 拉取上游 frp：

```text
github.com/fatedier/frp@v0.69.1
```

构建产物只上传 GitHub Releases，不提交到源码仓库。
