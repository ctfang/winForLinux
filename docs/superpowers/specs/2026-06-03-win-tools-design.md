# Win-Tools Design Specification

> Go 语言实现的 Windows 命令行工具集，补充缺失的 Linux 命令，增强 Windows 原生命令。

## 概述

### 目标

1. 实现 Windows 缺失的常用 Linux 文件操作命令
2. 覆盖 Windows `ls` 命令，默认 Windows 风格，遇到 Linux 独有选项时自动切换
3. 每个命令编译为独立 `.exe`，可直接放入 PATH 使用
4. 安装时自动备份系统已有同名命令，支持恢复

### 命令清单

| 优先级 | 命令 | 复杂度 | 说明 |
|--------|------|--------|------|
| P0 | `ls` | 高 | 目录列表，双模式输出 |
| P0 | `grep` | 中 | 文本搜索，支持正则 |
| P0 | `cat` | 低 | 文件输出 |
| P0 | `which` | 低 | 查找可执行文件 |
| P1 | `head` | 低 | 文件头输出 |
| P1 | `tail` | 低 | 文件尾输出 |
| P1 | `sort` | 中 | 排序 |
| P1 | `uniq` | 低 | 去重 |
| P1 | `wc` | 低 | 统计 |
| P1 | `touch` | 低 | 创建/更新文件时间 |
| P1 | `basename` | 低 | 提取文件名 |
| P1 | `dirname` | 低 | 提取目录路径 |
| P2 | `sed` | 高 | 流编辑器 |
| P2 | `find` | 高 | 文件查找 |
| P2 | `diff` | 高 | 文件差异 |
| P2 | `xargs` | 中 | 参数传递 |
| P2 | `du` | 中 | 磁盘使用 |
| P2 | `df` | 中 | 磁盘空间 |

---

## 架构设计

### 方案：独立二进制 + 共享库（方案 C）

```
win-linux/
├── go.work                      # Go workspace 管理
├── go.mod                       # 根模块
│
├── cmd/                         # 可执行文件入口（只做入口调用 pkg）
│   ├── ls/main.go
│   ├── grep/main.go
│   ├── sed/main.go
│   ├── cat/main.go
│   ├── head/main.go
│   ├── tail/main.go
│   ├── find/main.go
│   ├── sort/main.go
│   ├── uniq/main.go
│   ├── wc/main.go
│   ├── diff/main.go
│   ├── touch/main.go
│   ├── which/main.go
│   ├── basename/main.go
│   ├── dirname/main.go
│   ├── xargs/main.go
│   ├── du/main.go
│   └── df/main.go
│
├── pkg/                         # 命令核心逻辑（可被其他项目 import）
│   ├── ls/
│   │   ├── ls.go               # 主逻辑
│   │   ├── windows.go          # Windows dir 格式输出
│   │   ├── linux.go            # Linux ls 格式输出
│   │   ├── color.go            # 颜色方案
│   │   └── ls_test.go
│   ├── grep/
│   │   ├── grep.go
│   │   └── grep_test.go
│   └── ...（每个命令类似结构）
│
├── shared/                      # 共享工具库
│   ├── osutil/
│   │   ├── fileinfo.go         # 跨平台文件信息
│   │   ├── permission.go       # Windows 权限映射
│   │   └── symlink.go          # 符号链接处理
│   ├── color/
│   │   ├── ansi.go             # ANSI 颜色码
│   │   └── terminal.go         # 终端检测
│   └── argparse/
│       ├── parser.go           # 参数解析（支持 - 和 --）
│       └── parser_test.go
│
├── test/                        # 集成测试
│   └── compat/                 # 与原生命令对比测试
│
├── build.ps1                    # 编译脚本
├── install.ps1                  # 安装脚本（带备份）
├── uninstall.ps1                # 卸载脚本（带恢复）
└── README.md
```

### 模块职责

| 目录 | 职责 | 依赖 |
|------|------|------|
| `cmd/` | 入口，解析参数，调用 `pkg/` | `pkg/` |
| `pkg/` | 命令核心逻辑 | `shared/` |
| `shared/` | 底层工具库，无业务逻辑 | 标准库 |
| `test/` | 集成测试 | `pkg/` |

### Go Workspace 结构

```go
// go.work
use (
    .
    ./cmd/ls
    ./cmd/grep
    ./cmd/cat
    // ... 每个命令独立 module
)
```

每个 `cmd/*/` 有独立 `go.mod`，通过 `go.work` 统一管理依赖。

---

## ls 命令详细设计

### 行为逻辑

```
用户输入 ls [args...]
    │
    ├─ 无参数 → Windows dir 风格输出（当前目录）
    │
    ├─ 存在任何 Linux 独有选项（-l, -a, -h, -t, -R, -r, --color, -1, -d, -F, -S, -i, --all）？
    │   └─ 是 → Linux ls 风格输出（忽略 Windows 选项）
    │
    └─ 仅包含 Windows 选项或路径 → Windows dir 风格输出
```

**模式切换规则**：
- 只要参数中出现任何一个 Linux 独有选项，整个命令切换到 Linux 模式
- `ls -a` → Linux 模式（`-a` 是 Linux 选项）
- `ls /W` → Windows 模式（`/W` 是 Windows 选项）
- `ls -a /W` → Linux 模式（`-a` 触发 Linux 模式，`/W` 被忽略）
- `ls C:\Users` → Windows 模式（纯路径）
- `ls -l C:\Users` → Linux 模式（`-l` 触发）
- 未知的 `-` 开头选项视为 Linux 选项，触发 Linux 模式

### Windows 支持的选项（透传给 dir 逻辑）

| 选项 | 含义 |
|------|------|
| 路径 | `ls C:\Users` 或 `ls .` |
| `/W` | 宽列表格式 |
| `/B` | 裸格式（只有文件名） |
| `/S` | 递归子目录 |
| `/O` | 排序方式 |
| `/A` | 显示特定属性的文件 |

### Linux 独有选项（触发 Linux 模式）

| 选项 | 含义 | 实现方式 |
|------|------|----------|
| `-l` | 长格式 | 读取文件属性，格式化输出 |
| `-a` / `--all` | 显示隐藏文件 | 遍历时包含 `.` 开头文件 |
| `-h` | 人类可读大小 | 转换字节为 KB/MB/GB |
| `-t` | 按时间排序 | 使用 `os.Stat()` 的 ModTime |
| `-R` | 递归 | 递归遍历目录 |
| `-r` | 反序 | 反转排序结果 |
| `--color` | 颜色输出 | ANSI 颜色码 |
| `-1` | 每行一个文件 | 简单列表 |
| `-d` | 目录本身 | 不展开目录内容 |
| `-F` | 文件类型标记 | 追加 `/` `*` `@` 等 |
| `-S` | 按大小排序 | 使用文件大小 |
| `-i` | 显示 inode | 显示文件索引号 |

### 长格式 (-l) 输出

```
drwxr-xr-x  1 user  group  4096  Jun  3 12:00  Documents/
-rw-r--r--  1 user  group  1.2M  Jun  3 11:30  report.pdf
-rwxr-xr-x  1 user  group  8.5K  Jun  3 10:15  script.sh*
```

**Windows 权限映射**：

| Linux 权限 | Windows 映射 |
|------------|-------------|
| 读 (r) | `FILE_GENERIC_READ` |
| 写 (w) | `FILE_GENERIC_WRITE` |
| 执行 (x) | `FILE_GENERIC_EXECUTE` 或 `.exe/.bat/.cmd/.ps1` |
| 用户/组 | 从 Windows 安全描述符获取 |

### 颜色方案

| 类型 | 颜色 | ANSI 码 |
|------|------|---------|
| 目录 | 蓝色 | 34 |
| 可执行文件 | 绿色 | 32 |
| 符号链接 | 青色 | 36 |
| 压缩文件 | 红色 | 31 |
| 图片文件 | 洋红 | 35 |
| 普通文件 | 默认色 | 0 |

### Windows dir 风格输出

```
 Volume in drive C is Windows
 Volume Serial Number is XXXX-XXXX

 Directory of D:\github.cocm\win-linux

2024/06/03  12:00    <DIR>          .
2024/06/03  12:00    <DIR>          ..
2024/06/03  11:30         1,258,496  report.pdf
2024/06/03  10:15             8,704  script.sh
               2 File(s)      1,267,200 bytes
               2 Dir(s)   50,123,456,789 bytes free
```

---

## 其他命令设计

### grep

```go
type Options struct {
    Pattern     string   // 正则表达式
    Files       []string // 文件列表（空=stdin）
    IgnoreCase  bool     // -i
    LineNumber  bool     // -n
    Count       bool     // -c
    Recursive   bool     // -r / -R
    Invert      bool     // -v
    OnlyMatch   bool     // -o
    WordRegexp  bool     // -w
    ExtendedReg bool     // -E (扩展正则)
    FixedString bool     // -F (固定字符串)
    Color       string   // --color=auto/always/never
}
```

### sed

```go
type Options struct {
    Scripts    []string // sed 脚本表达式
    InPlace    bool     // -i (原地编辑)
    Files      []string // 文件列表
    Quiet      bool     // -n (静默)
    Expression string   // -e 表达式
    Separator  byte     // 分隔符（默认 /）
}
```

**支持的 sed 命令**：

| 命令 | 说明 | 示例 |
|------|------|------|
| `s/old/new/` | 替换 | `s/foo/bar/g` 全局替换 |
| `d` | 删除行 | `/^$/d` 删除空行 |
| `p` | 打印行 | `-n '/pattern/p'` |
| `i\text` | 插入文本 | 在匹配行前插入 |
| `a\text` | 追加文本 | 在匹配行后追加 |
| `c\text` | 替换整行 | 将匹配行替换为文本 |
| `y/src/dst/` | 字符转换 | `y/abc/ABC/` |
| `=` | 打印行号 | |
| `q` | 退出 | `10q` 打印10行后退出 |

**地址范围**：
- 单行：`5d`（第5行）
- 范围：`1,10d`（第1-10行）
- 模式：`/start/,/end/d`
- 步进：`0~2d`（每隔一行）

不追求完整 sed 实现，覆盖最常用场景。

### find

```go
type Options struct {
    Paths      []string // 搜索路径（支持 Windows 路径 C:\, UNC \\server\share）
    Name       string   // -name 模式（支持 * 和 ? 通配符）
    Type       string   // -type f/d/l（l=symlink/junction）
    Size       string   // -size +1M（支持 K/M/G 后缀）
    MTime      string   // -mtime -7（天数）
    Exec       string   // -exec 命令 {}（{} 替换为找到的文件路径）
    MaxDepth   int      // -maxdepth
    MinDepth   int      // -mindepth
    Delete     bool     // -delete
    Print      bool     // -print (默认)
    Iname      string   // -iname（不区分大小写的 name）
    IRegexp    string   // -iregexp（不区分大小写的正则）
}
```

**Windows 特殊处理**：
- 路径分隔符：统一使用 `\`，但也接受 `/`
- UNC 路径：支持 `\\server\share` 格式
- 长路径：支持 `\\?\` 前缀突破 260 字符限制
- 隐藏文件：Windows 隐藏属性通过 `FILE_ATTRIBUTE_HIDDEN` 检测
- 符号链接：Windows 的 symlink 和 junction 都视为 `l` 类型

### 简单命令

| 命令 | 核心逻辑 |
|------|----------|
| `cat` | 读取文件 → 输出到 stdout，支持 `-n` 行号 |
| `head` | 读取前 N 行 (`-n`) 或前 N 字节 (`-c`) |
| `tail` | 读取后 N 行，支持 `-f` 跟踪 |
| `sort` | 排序，支持 `-r` 反序、`-n` 数值、`-u` 去重、`-k` 指定列 |
| `uniq` | 去重相邻行，支持 `-c` 计数、`-d` 只显示重复 |
| `wc` | 统计 `-l` 行、`-w` 词、`-c` 字节、`-m` 字符 |
| `touch` | 创建文件或更新时间戳，支持 `-a` 访问时间、`-m` 修改时间 |
| `which` | 在 PATH 中查找可执行文件，支持 `-a` 显示所有 |
| `basename` | 提取文件名，支持去除后缀 |
| `dirname` | 提取目录路径 |
| `xargs` | 从 stdin 读取参数，传递给命令，支持 `-I{}` 替换、`-n` 每次参数数、`-P` 并行 |
| `du` | 磁盘使用，支持 `-h` 人类可读、`-s` 汇总、`--max-depth` |
| `df` | 磁盘空间，支持 `-h` 人类可读 |
| `diff` | 文件差异，支持 `-u` unified 格式、`-r` 递归 |

---

## 构建与安装系统

### build.ps1

```powershell
param(
    [string]$OutputDir = ".\dist",
    [string]$Version = "1.0.0",
    [switch]$Clean
)

if ($Clean) { Remove-Item -Recurse -Force $OutputDir -EA SilentlyContinue }
New-Item -ItemType Directory -Force $OutputDir | Out-Null

$commands = @("ls","grep","sed","cat","head","tail","find","sort","uniq",
              "wc","diff","touch","which","basename","dirname","xargs","du","df")

foreach ($cmd in $commands) {
    Write-Host "Building $cmd..." -ForegroundColor Cyan
    go build -o "$OutputDir\$cmd.exe" -ldflags "-s -w -X main.version=$Version" ".\cmd\$cmd"
}
Write-Host "Build complete! Output: $OutputDir" -ForegroundColor Green
```

### install.ps1（带备份）

```powershell
param(
    [string]$InstallDir = "$env:LOCALAPPDATA\win-tools",
    [string]$BackupDir  = "$env:LOCALAPPDATA\win-tools-backup",
    [switch]$AddToPath,
    [switch]$Force,
    [switch]$Restore
)

# 恢复模式
if ($Restore) {
    # 从 BackupDir 恢复原始命令
    # 恢复后删除已安装的命令
}

# 安装流程
# 1. 扫描要安装的命令
# 2. 检查 PATH 中是否有同名文件，有则备份
# 3. 备份元信息（原始路径、版本、日期）→ metadata.json
# 4. 安装新命令
# 5. 可选添加到 PATH
```

### 备份目录结构

```
%LOCALAPPDATA%\win-tools-backup\
├── grep\
│   ├── grep.exe.original
│   └── metadata.json
├── ls\
│   ├── ls.exe.original
│   └── metadata.json
└── restore-log.json
```

### 恢复流程

```powershell
# 恢复所有备份
.\install.ps1 -Restore

# 恢复特定命令
.\install.ps1 -Restore -Command grep
```

---

## 错误处理

| 场景 | 处理方式 |
|------|----------|
| 文件不存在 | stderr 输出错误，退出码 1 |
| 权限不足 | 友好错误信息，提示管理员运行 |
| 参数错误 | 显示用法帮助，退出码 2 |
| 管道中断 | 忽略 SIGPIPE，正常退出 |
| 编码问题 | 默认 UTF-8，自动检测 BOM |

## 共存策略

当 PATH 中存在同名命令（如 Git Bash 的 `grep.exe`）：

1. **安装时**：备份原有命令，显示警告
2. **运行时**：检测到被覆盖的命令时显示提示
3. **卸载时**：恢复原始命令

## 测试策略

```
每个命令的测试：
├── 单元测试（pkg/*/）     # 核心逻辑测试
├── 集成测试（test/）       # 端到端测试
└── 兼容性测试              # 与原生命令输出对比
```

**测试方法**：
- 对比测试：运行我们的命令和原生命令，比较输出
- 边界测试：空文件、大文件、特殊字符、长路径
- 管道测试：与其他命令组合使用
- Windows 特性测试：UNC 路径、长文件名、Unicode 文件名

## 版本信息

每个命令支持 `--version`：

```
$ grep --version
win-tools grep 1.0.0
windows/amd64, built with Go 1.22
```
