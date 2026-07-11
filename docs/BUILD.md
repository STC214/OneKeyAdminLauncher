# 构建说明

## 目录

`OneKeyAdminOpenPro` 是独立的 Go + Win32 重构工程，位于：

```text
F:\Project\04_Desktop_Dev_Utilities\OneKeyAdminOpenPro
```

当前结构：

```text
OneKeyAdminOpenPro/
  ICO.ico
  cmd/program-launch-manager/
    main.go
    app.rc
    app.manifest
    app.syso
  internal/
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
dist/程序启动管理器.exe
```

运行时配置保存在 exe 同级目录：

```text
dist/launcher_config.json
```

如果直接从其他目录运行 exe，配置也会保存在该 exe 所在目录。

## 管理员权限

最终 exe 嵌入 `requireAdministrator` manifest。双击启动时 Windows 会弹出 UAC。用户拒绝 UAC 时程序退出。

## 构建命令

```powershell
cd F:\Project\04_Desktop_Dev_Utilities\OneKeyAdminOpenPro
.\scripts\build.ps1
```

构建脚本会：

- 运行 `go test ./...`。
- 使用仓库内的 `ICO.ico`；如果本地存在更新的 `ICO.png`，则使用 Python + Pillow 重新生成多尺寸图标。
- 使用 `windres --codepage=65001` 将 `app.rc` 编译为 `app.syso`。
- 构建 `dist/程序启动管理器.exe`。

如果 `ICO.ico` 已存在且比 `ICO.png` 更新，则直接复用。

## 验证项

- `go test ./...` 通过。
- `go vet ./...` 通过。
- `dist/程序启动管理器.exe` 生成。
- exe 启动时请求管理员权限。
- 标题栏、任务栏和托盘图标显示嵌入式 ico。
- 首次保存配置后，在 exe 同级目录生成 `launcher_config.json`。
- 程序列表超出高度时滚动条可用。
- 窗口最小宽度为当前屏幕宽度的一半。
