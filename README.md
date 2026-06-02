# frplib

`frplib` 是官方 frp 的极简 Android AAR 内核绑定库。

只负责在 Android 中启动、停止、重启 `frpc` / `frps`。配置直接使用官方 TOML，不提供 UI、通知栏、保活、权限申请、订阅解析、配置编辑器，也不支持 JitPack。

## 下载

到 GitHub Releases 下载 AAR：

```text
frplib-arm64-v8a.aar
frplib-armeabi-v7a.aar
frplib-x86_64.aar
frplib-universal.aar
```

普通业务推荐直接使用：

```text
frplib-universal.aar
```

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

## 使用

单实例：

```kotlin
val err = Frplib.StartClient(frpcToml)
if (err.isNotEmpty()) {
    // 处理错误
}

Frplib.StopClient()
```

服务端：

```kotlin
Frplib.StartServer(frpsToml)
Frplib.StopServer()
```

多实例：

```kotlin
Frplib.StartClientWithID("client-a", frpcTomlA)
Frplib.StartClientWithID("client-b", frpcTomlB)
Frplib.StopClientWithID("client-a")
```

日志：

```kotlin
Frplib.SetLogCallback(object : FrpLogCallback {
    override fun onLog(instanceId: String, type: String, level: String, message: String) {
        // type = client / server
        // level = debug / info / warn / error
    }
})
```

## reload 说明

它们等价于 safe restart：

```text
先校验新 TOML，校验通过后停止旧实例，再用新 TOML 启动。
```

## 支持 ABI

支持：

```text
arm64-v8a
armeabi-v7a
x86_64
universal
```

## 发布

在 GitHub Actions 手动运行 `发布 AAR`，输入发布 tag，例如：

```text
v0.69.1
```

流水线会使用同一个 tag 拉取上游 frp：

```text
github.com/fatedier/frp@v0.69.1
```

然后构建并上传 AAR 到对应 GitHub Release。
