# OneKeyAdminLauncher Modern UI

这是 `OneKeyAdminLauncher` 的视觉优化子项目。代码从稳定版完整复制而来，功能逻辑保持一致，本分支只调整 Win32 UI 的视觉效果、间距、字体和配色。

## 功能

- 一键启动所有已启用的程序项。
- 一键关闭所有已启用项对应的进程。
- 支持普通 `.exe`、Windows 快捷方式 `.lnk` 和 UWP 启动项。
- 支持手动从当前进程列表中绑定检测和关闭使用的进程名。
- 每秒自动刷新运行状态。
- 关闭程序时先尝试正常关闭窗口，失败后再使用 Win32 API 强制结束进程。
- 启动和关闭辅助流程不弹出 cmd 窗口。
- 支持最小化到系统托盘。
- 自动保存程序列表、启用状态、UWP 标记、进程绑定和窗口位置/大小。
- 图标嵌入 exe，覆盖标题栏、Windows 任务栏和系统托盘。
- 列表使用虚拟行控件池，程序项较多时只渲染可见行，滚动更稳且不易卡顿。
- 窗口最小宽度按当前屏幕宽度的一半计算，同时保留 Modern UI 布局所需的最低宽度，避免控件被压坏。
- Modern UI 版本使用更现代的深色调、Segoe UI 字体、扁平按钮和右侧固定对齐布局。
- 顶部操作和行内操作按钮使用彩色自绘样式，UWP 标记文字使用更亮的前景色。
- 标题栏使用 DWM 深色标题栏属性，保留系统原生拖动、最小化、最大化和关闭按钮。

## UI 行为

- 主界面保持紧凑的深色 Win32 工具风格。
- 顶部工具栏固定显示，不随列表滚动。
- 程序列表超出窗口高度后使用右侧滚动条。
- 滚动时复用可见行控件，避免为所有配置项创建大量 Win32 子窗口。
- 主列表使用 Modern UI 自绘细滚动条，支持滚轮、轨道点击和滑块拖动。
- 窗口缩放时会擦除旧绘制区域，避免控件边框残影；滚动时使用轻量刷新，减少闪动。
- 最小化按钮会隐藏主窗口并进入托盘；双击托盘图标恢复。

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
cd .\OneKeyAdminOpenPro_ModernUI
.\scripts\build.ps1
```

构建脚本会自动执行测试并生成：

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
