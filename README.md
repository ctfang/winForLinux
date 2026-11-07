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
| `echo` / `pwd` / `tee` / `clear` | 输出、路径、分流、清屏（`echo -n/-e`，`tee -a`） |
| `find` | 文件查找（`-name/-iname/-type/-maxdepth/-mindepth/-size/-mtime`） |
| `sed` | 基础流编辑（`s///[gi]`、`d`、`p`、`q`，`-n/-i/-e`） |
| `diff` | 文件差异（`-u` unified，`-q` brief） |
| `du` / `df` | 目录占用 / 磁盘空间（`-h`，`du -s/-d`） |
| `xargs` | 参数传递（`-n/-I/-0/-d/-r`；PowerShell 下 `{}` 请加引号） |
| `wfl` | Linux 风格命令组合器（`\|` `&&` `\|\|` `;` `&`） |

## 安装

```powershell
.\build.ps1
.\install.ps1 -AddToPath -Force
```

安装目录：`%LOCALAPPDATA%\win-tools`

安装脚本会：
1. 把所有 `.exe` 拷到安装目录并加入用户 PATH
2. 备份可能被覆盖的系统同名程序（如 Windows `find.exe`）
3. 在当前会话移除会遮挡的 PowerShell 别名（`ls`/`cat`/`cp`/`mv`/`rm`/`mkdir`/`sort`/`echo`/`pwd`/`clear`/`diff`/`tee` 等）
4. 更新 PowerShell profile：优先 win-tools，并启用 **UTF-8**（见下）

请**新开一个 PowerShell 窗口**后再用。

### PowerShell UTF-8 与管道

Windows PowerShell 5.1 默认把 `$OutputEncoding` 设为 ASCII。原生 `.exe` 之间用 `|` 时，非 ASCII 字符会被重编码成 `?`，例如：

```powershell
cat README.md | grep -i win   # 未设 UTF-8 时中文变成 ??
```

`install.ps1` 写入的 profile 会在交互式 ConsoleHost 中设置：

- `chcp 65001`（仅在控制台未重定向时）
- `[Console]::OutputEncoding` / `InputEncoding` → UTF-8（无 BOM）
- `$OutputEncoding` → UTF-8（无 BOM）

这样 `cat … | grep …` 能正确传递中文。`wfl -c "…"` 与 `cmd.exe` 管道本身不依赖该设置。新开窗口后生效；当前会话也可手动执行上述设置。

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
echo -e "a\tb"
pwd
find . -name "*.go" -type f
find docs -maxdepth 2 -iname "*.md"
sed "s/foo/bar/g" file.txt
sed -i "s/old/new/" file.txt
diff -u a.txt b.txt
du -sh .
df -h
echo one two three | xargs -n 1 echo
# PowerShell 下 {} 是脚本块，请加引号：
"a`nb" | xargs -I '{}' echo 'got-{}'
```

### 管道与命令组合

**在 PowerShell 里**，原生 `|` 已经可以连接这些 `.exe`（需 UTF-8 profile，见上）：

```powershell
cat README.md | grep -i win
echo hello | tee out.txt | grep hello
Get-Content README.md | wc -l
```

PowerShell 对 `&` / `&&` / `||` 的语义与 Linux 不完全相同。需要 Linux 风格组合时，用 **`wfl`**：

```powershell
wfl -c "echo hello | grep hello"
wfl -c "echo -e 'a\nb' | xargs -I {} echo got-{}"
wfl -c "find . -name *.txt | xargs cat"
wfl -c "false || echo fallback"
wfl -c "true && echo ok"
wfl -c "echo bg & echo fg"
wfl -c "ls -la ; echo done"
```

`wfl` 会优先调用与自己同目录的 win-tools 命令（含 `echo.exe`）。

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
.\build.ps1
.\install.ps1 -AddToPath -Force
```

当前版本：1.2.0