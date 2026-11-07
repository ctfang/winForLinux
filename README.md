# winForLinux (win-tools)

在 Windows 上提供常用 Linux 命令的独立可执行文件（Go 实现），可直接在 PowerShell / cmd 中使用，无需 WSL。

## 前置要求

- Go 1.22+（仅从源码编译时需要）
- PowerShell 5.1+

## 已实现命令

| 命令 | 说明 |
|------|------|
| `ls` | 目录列表（默认 Windows dir 风格；遇 `-l/-a/-h/-R` 等切到 Linux 风格） |
| `grep` | 文本搜索（支持合并短选项如 `-in`，以及 `-r` 递归） |
| `cat` | 文件输出 |
| `which` | 在 PATH 中查找可执行文件 |
| `head` / `tail` | 文件头 / 尾（`tail -f` 跟踪） |
| `wc` | 行/词/字节统计 |
| `touch` | 创建文件或更新时间戳 |
| `sort` / `uniq` | 排序 / 去重 |
| `basename` / `dirname` | 路径拆分 |
| `cp` / `mv` / `rm` / `mkdir` | 常用文件操作（`cp -r`、`rm -rf`、`mkdir -p`） |
| `wfl` | Linux 风格命令组合器（`\|` `&&` `\|\|` `;` `&`） |

## 安装

```powershell
.\build.ps1
.\install.ps1 -AddToPath -Force
```

安装目录：`%LOCALAPPDATA%\win-tools`

安装脚本会：
1. 把所有 `.exe` 拷到安装目录并加入用户 PATH
2. 备份可能被覆盖的系统同名程序
3. 在当前会话移除会遮挡的 PowerShell 别名（`ls`/`cat`/`cp`/`mv`/`rm`/`mkdir`/`sort`）

请**新开一个 PowerShell 窗口**后再用。配置文件也会在启动时自动去掉这些别名，并优先使用 win-tools。

卸载 / 恢复：

```powershell
.\uninstall.ps1
# 或
.\install.ps1 -Restore
```

## 用法示例

```powershell
ls -lah
ls -R docs
grep -in "error" *.log
grep -r "TODO" .
cat -n README.md
head -n 5 README.md
tail -n 3 README.md
wc -l README.md
touch newfile.txt
mkdir -p a\b\c
cp -r docs docs_backup
sort file.txt | uniq
```

### 管道与命令组合

**在 PowerShell 里**，原生 `|` 已经可以连接这些 `.exe`：

```powershell
cat README.md | grep -i win
Get-Content README.md | wc -l
```

PowerShell 对 `&` / `&&` / `||` 的语义与 Linux 不完全相同。需要 Linux 风格组合时，用 **`wfl`**：

```powershell
wfl -c "echo hello | grep hello"
wfl -c "false || echo fallback"
wfl -c "true && echo ok"
wfl -c "echo bg & echo fg"
wfl -c "ls -la ; echo done"
```

`wfl` 会优先调用与自己同目录的 win-tools 命令。

## 开发

```
winForLinux/
├── cmd/          # 各命令入口
├── pkg/          # 命令核心逻辑
├── shared/       # argparse / color / osutil
├── test/         # 集成测试
├── build.ps1
├── install.ps1
└── uninstall.ps1
```

添加新命令：在 `pkg/` 与 `cmd/` 各加一处，更新 `build.ps1` 的命令列表。

```powershell
go test ./pkg/... ./shared/...
.\build.ps1
go test ./test/ -v
```

## 设计说明

- 每个命令编译为独立 `.exe`，放入 PATH 即可用（方案与设计文档一致）
- 管道在宿主 shell（PowerShell）中天然可用；`&&` / `||` / `&` / `;` 由 `wfl` 提供
- 不追求完整 GNU 兼容，覆盖常用基础场景

## 许可证

MIT