# OneKeyAdminLauncher

一个使用 Go + Win32 原生控件编写的 Windows 程序启动管理器。程序自身以管理员权限运行，用于一键以管理员权限启动、检测和关闭一组常用工具。

## 功能

- 一键启动所有已启用的程序项。
- 一键关闭所有已启用项对应的进程。
- 支持普通 `.exe`、Windows 快捷方式 `.lnk` 和 UWP 启动项。
- 支持手动从当前进程列表中绑定检测/关闭用的进程名。
- 每秒自动刷新运行状态。
- 关闭时先尝试正常关闭窗口，失败后再使用 Win32 API 强制结束进程。
- 启动和关闭辅助流程不弹出 cmd 窗口。
- 支持最小化到系统托盘。
- 自动保存程序列表、启用状态、UWP 标记、进程绑定和窗口位置。
- 图标嵌入 exe，覆盖标题栏、任务栏和托盘。

## 配置文件

运行时配置保存在 exe 同级目录：

```text
launcher_config.json
```

如果检测到旧版配置：

```text
data/launcher_config.json
```

程序启动时会自动迁移到 exe 同级目录。

## 构建

构建环境需要：

- Windows
- Go
- `windres`
- Python + Pillow

构建命令：

```powershell
cd F:\Project\04_Desktop_Dev_Utilities\OneKeyAdminOpenPro
.\scripts\build.ps1
```

构建输出：

```text
dist/程序启动管理器.exe
```

构建脚本会将根目录本地素材 `ICO.png` 转为多尺寸 `ICO.ico`，再通过 `app.rc` 生成 `app.syso` 并嵌入最终 exe。图片/图标文件、`app.syso` 和 `dist/` 都不提交到仓库。

## 项目结构

```text
cmd/program-launch-manager/    程序入口、manifest 和资源脚本
internal/config/               配置加载、保存和迁移
internal/process/              启动、进程枚举和关闭逻辑
internal/winui/                Win32 UI
scripts/build.ps1              构建脚本
docs/                          规格和构建说明
```

## 备注

程序 manifest 使用 `requireAdministrator`。双击启动时会触发 UAC；如果拒绝管理员权限，程序会直接退出。
