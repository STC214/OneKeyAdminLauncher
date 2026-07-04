# 构建说明

## 目录

`program-launch-manager` 是独立的 Go + Win32 重构工程，放在当前项目目录下，不依赖 APK 逆向文件。

计划结构：

```text
program-launch-manager/
  ICO.png
  ICO.ico
  cmd/program-launch-manager/
    main.go
    app.rc
    app.manifest
    app.syso
  internal/
    app/
    config/
    process/
    winui/
  docs/
    SPEC.md
    BUILD.md
  scripts/
    build.ps1
  go.mod
```

## 构建目标

输出文件：

```text
program-launch-manager/dist/程序启动管理器.exe
```

运行时配置：

```text
program-launch-manager/dist/launcher_config.json
```

## 管理员权限

最终 exe 需要嵌入 `requireAdministrator` manifest。双击启动时 Windows 会弹出 UAC。用户拒绝 UAC 时程序退出。

## 构建命令

首版实现完成后使用：

```powershell
cd F:\Project\04_Desktop_Dev_Utilities\OneKeyAdminOpenPro
.\scripts\build.ps1
```

构建脚本会使用 Python + Pillow 将根目录 `ICO.png` 转为多尺寸 `ICO.ico`；若 `ICO.ico` 已存在且比 `ICO.png` 更新，则直接复用。

验证项：

- `go test ./...` 通过。
- `dist/程序启动管理器.exe` 生成。
- exe 启动时请求管理员权限。
- `ICO.png` 会转换为根目录 `ICO.ico`，并通过 `app.rc`/`app.syso` 嵌入最终 exe。
- 首次启动会创建 `dist/launcher_config.json`。
