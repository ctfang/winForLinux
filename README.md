# Win-Tools

Go 语言实现的 Windows 命令行工具集，补充缺失的 Linux 命令，增强 Windows 原生命令。

## 前置要求

- Go 1.22.0 或更高版本
- PowerShell 5.1 或更高版本

## 功能

### 已实现命令

| 命令 | 说明 | 状态 |
|------|------|------|
| `ls` | 目录列表（Windows/Linux 双模式） | ✅ |
| `grep` | 文本搜索 | ✅ |
| `cat` | 文件输出 | ✅ |
| `which` | 查找可执行文件 | ✅ |

### ls 命令特性

- **Windows 模式**（默认）：兼容 Windows dir 命令
- **Linux 模式**：支持 `-l`, `-a`, `-h`, `-t`, `-R` 等选项
- **自动切换**：使用 Linux 选项时自动切换到 Linux 模式

```bash
ls          # Windows dir 风格
ls -l       # Linux 长格式
ls -lah     # Linux 风格，显示所有文件、人类可读大小
ls /W       # Windows 宽列表
```

### grep 命令特性

- 支持正则表达式
- 支持 `-i`（忽略大小写）、`-n`（行号）、`-c`（计数）等选项
- 支持从文件或 stdin 读取

```bash
grep "pattern" file.txt
echo "hello" | grep -i "HELLO"
grep -n "error" *.log
```

> **注意**：`*.log` 等通配符由 shell 展开，而非 grep 命令本身。在 PowerShell 中，通配符行为可能与 bash 不同。

### cat 命令特性

- 支持 `-n`（行号）、`-E`（行尾标记）、`-T`（制表符显示）、`-s`（压缩空行）

```bash
cat file.txt
cat -n file.txt
cat -E file.txt
```

### which 命令特性

- 在 PATH 中查找可执行文件
- 支持 `-a` 显示所有匹配

```bash
which go
which -a python
```

## 安装

### 从源码构建

```powershell
# 编译所有命令
.\build.ps1

# 安装到系统
.\install.ps1 -AddToPath
```

> 编译输出位于 `.\dist\` 目录。

### 卸载

```powershell
# 卸载（保留备份）
.\uninstall.ps1

# 恢复原始命令
.\uninstall.ps1 -Restore
```

## 开发

### 项目结构

```
win-linux/
├── cmd/              # 可执行文件入口
├── pkg/              # 命令核心逻辑
├── shared/           # 共享工具库
├── test/             # 集成测试
├── build.ps1         # 编译脚本
├── install.ps1       # 安装脚本
├── uninstall.ps1     # 卸载脚本
├── go.mod            # Go 模块定义
└── go.work           # Go 工作空间配置
```

### 添加新命令

1. 在 `pkg/` 创建命令包
2. 在 `cmd/` 创建入口
3. 更新 `build.ps1` 的命令列表
4. 添加测试

### 运行测试

```bash
# 单元测试
go test ./...

# 集成测试（需要先编译）
.\build.ps1
go test ./test/ -v
```

## 许可证

MIT
