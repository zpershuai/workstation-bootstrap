# Dwell CLI 设计文档

## 1. 项目概述

### 1.1 背景
原项目使用 shell 脚本管理开发环境初始化，每次修改需要编辑多个脚本文件，流程繁琐且容易出错。需要抽象出一个命令行工具，支持模块化管理和独立更新。

### 1.2 目标
- 单一二进制工具管理整个开发环境
- 模块化架构，支持独立更新每个组件
- 向后兼容现有 repos.lock 配置
- 可扩展支持更多模块类型（brew, npm, dotfiles）

### 1.3 范围
**In Scope:**
- Git 仓库管理（clone, pull, symlink）
- 配置文件解析（YAML + repos.lock）
- 命令行接口（sync, status, doctor, init）
- 健康检查系统

**Out of Scope (Future):**
- Homebrew 包管理
- npm 全局包管理  
- Dotfiles 同步
- 系统默认值设置

---

## 2. 架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                        CLI Layer                             │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│  │  sync    │ │  status  │ │  doctor  │ │   init   │       │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘       │
└───────┼────────────┼────────────┼────────────┼─────────────┘
        │            │            │            │
        └────────────┴────────────┴────────────┘
                         │
┌────────────────────────┼────────────────────────────────────┐
│                   Core Engine                                │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Module Registry                         │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐    │  │
│  │  │  Git    │ │  Brew   │ │   NPM   │ │Dotfiles │    │  │
│  │  │ Module  │ │ Module  │ │ Module  │ │ Module  │    │  │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘    │  │
│  └──────────────────────────────────────────────────────┘  │
│                         │                                    │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Config Loader                           │  │
│  │         (YAML + repos.lock support)                  │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 模块系统

#### 2.2.1 Module 接口

所有模块类型必须实现此接口：

```go
type Module interface {
    Name() string                    // 模块唯一标识
    Type() string                    // 模块类型（git, brew, npm, etc.）
    Description() string             // 人类可读描述
    Status(ctx context.Context) (*State, error)   // 获取状态
    Sync(ctx context.Context) error               // 同步到期望状态
    Check(ctx context.Context) []CheckResult      // 健康检查
}
```

#### 2.2.2 状态类型

```go
type Status string

const (
    StatusOK       Status = "ok"      // 已是最新
    StatusPending  Status = "pending" // 等待同步
    StatusSyncing  Status = "syncing" // 同步中
    StatusError    Status = "error"   // 发生错误
    StatusDirty    Status = "dirty"   // 本地有修改
    StatusBehind   Status = "behind"  // 落后远程
    StatusAhead    Status = "ahead"   // 领先远程
    StatusMissing  Status = "missing" // 未安装
    StatusUnknown  Status = "unknown" // 未知状态
)
```

---

## 3. 模块详细设计

### 3.1 Git 模块

#### 3.1.1 配置结构

```yaml
git:
  - name: nvim
    url: git@github.com:zpershuai/nvim.git
    path: ~/.dotfiles.d/repos/nvim
    ref: main
    links:
      - from: ~/.dotfiles.d/repos/nvim
        to: ~/.config/nvim
    post_sync: ~/.tmux/install.sh
```

#### 3.1.2 同步流程

```
Sync()
├── Ensure parent directory exists
├── Check if already cloned
│   ├── No  → git clone
│   └── Yes → skip
├── git fetch --all --tags
├── Checkout ref (if specified)
│   ├── Is branch? → git pull --ff-only
│   └── Is tag/commit? → just checkout
├── Create symlinks
│   ├── Remove existing file/link
│   ├── Ensure parent directory
│   └── os.Symlink(from, to)
└── Run post_sync script (if exists)
```

#### 3.1.3 状态检测

```
Status()
├── Check directory exists
│   └── No → StatusMissing
├── Check .git directory
│   └── No → StatusError
├── Check local changes
│   └── Yes → StatusDirty
├── Fetch remote (if ref specified)
├── Check behind upstream
│   └── Yes → StatusBehind + commit count
├── Check ahead upstream
│   └── Yes → StatusAhead + commit count
└── StatusOK
```

#### 3.1.4 健康检查

- **git-binary**: Git 是否安装且在 PATH 中
- **repo-access**: 是否能访问远程仓库
- **symlink-{name}**: 每个 symlink 是否正确配置

### 3.2 配置加载器

#### 3.2.1 加载优先级

1. `dwell.yaml`（新格式，首选）
2. `repos/repos.lock`（向后兼容）

#### 3.2.2 repos.lock 解析

```
name url dest [ref]
```

映射到 Git 模块配置：
- `name` → Module.Name
- `url` → Module.URL
- `dest` → Module.Path
- `ref` → Module.Ref

内置 link 配置映射（硬编码）：

| 模块名 | From | To |
|--------|------|-----|
| nvim | ~/.dotfiles.d/repos/nvim | ~/.config/nvim |
| tmux | ~/.dotfiles.d/repos/tmux | ~/.tmux |
| claudecode_dotfiles | ~/.dotfiles.d/repos/claudecode_dotfiles | ~/.claude |
| tpm | ~/.dotfiles.d/repos/tpm | ~/.tmux/plugins/tpm |
| zsh-syntax-highlighting | ~/.dotfiles.d/repos/zsh-syntax-highlighting | ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting |
| zsh-navigation-tools | ~/.dotfiles.d/repos/zsh-navigation-tools | ~/.oh-my-zsh/custom/plugins/zsh-navigation-tools |

post_sync 脚本（仅 tmux）：
- `~/.tmux/install.sh`

---

## 4. 命令设计

### 4.1 sync 命令

**用途**: 同步模块到期望状态

```bash
dwell sync [module] [flags]
```

**参数**:
- `module`: 可选，指定模块名。不提供则同步所有。

**Flags**:
- `--dry-run`: 预览变更，不实际执行

**输出示例**:
```
[nvim] ✓ synced (1.2s)
[tmux] ✓ synced (0.8s)
[claudecode_dotfiles] ✗ failed to pull: error message

Sync completed (2 succeeded, 1 failed)
```

### 4.2 status 命令

**用途**: 查看所有模块状态

```bash
dwell status
```

**输出示例**:
```
MODULE                  TYPE  STATUS  REF     MESSAGE
------                  ----  ------  ---     -------
nvim                    git   behind  main    2 commits behind remote
tmux                    git   ok      main    Up to date
claudecode_dotfiles     git   dirty   main    Local changes detected

Total: 3 modules
```

### 4.3 doctor 命令

**用途**: 运行健康检查

```bash
dwell doctor
```

**输出示例**:
```
[nvim] git
  ✓ git-binary
  ✓ repo-access
  ✓ symlink-nvim

[tmux] git
  ✓ git-binary
  ✓ repo-access
  ⚠ symlink-tmux: Symlink missing: ~/.tmux

✓ All checks passed (5/6)
```

### 4.4 init 命令

**用途**: 从 repos.lock 生成 dwell.yaml

```bash
dwell init [--force]
```

**Flags**:
- `--force`: 覆盖已存在的 dwell.yaml

---

## 5. 项目结构

```
dotfiles/
├── bin/
│   └── dwell                   # 编译后的二进制
├── internal/
│   ├── cmd/
│   │   ├── main.go            # CLI 入口
│   │   ├── sync.go            # sync 命令
│   │   ├── status.go          # status 命令
│   │   ├── doctor.go          # doctor 命令
│   │   └── init.go            # init 命令
│   └── pkg/
│       ├── modules/
│       │   └── module.go      # Module 接口
│       ├── git/
│       │   └── module.go      # Git 模块实现
│       └── config/
│           └── loader.go      # 配置加载器
├── dwell.yaml                 # 新配置格式
├── repos/repos.lock           # 旧配置格式（向后兼容）
├── go.mod                     # Go 模块
├── Makefile                   # 构建脚本
└── DWELL_README.md            # 使用文档
```

---

## 6. 依赖项

### 6.1 运行时依赖
- Go 1.21+
- Git
- SSH（用于 Git 认证）

### 6.2 编译依赖
```
github.com/fatih/color v1.16.0          # 终端颜色
github.com/urfave/cli/v2 v2.27.1        # CLI 框架
gopkg.in/yaml.v3 v3.0.1                 # YAML 解析
```

---

## 7. 向后兼容策略

### 7.1 配置加载
1. 首先尝试加载 `dwell.yaml`
2. 如果不存在，尝试加载 `repos/repos.lock`
3. 都不存在则报错

### 7.2 repos.lock 解析
- 忽略空行和 `#` 开头的注释行
- 每行格式：`name url dest [ref]`
- 少于 3 个字段的行被跳过
- 使用硬编码的 link 配置映射

### 7.3 迁移路径
```bash
# 用户可以选择保持使用 repos.lock（完全兼容）
# 或者迁移到新格式
dwell init  # 生成 dwell.yaml
# 之后可以删除 repos.lock
```

---

## 8. 测试策略

### 8.1 单元测试
- Module 接口实现测试
- Config 加载器测试（YAML + repos.lock）
- Git 操作封装测试

### 8.2 集成测试
- 完整 sync 流程测试
- repos.lock 向后兼容测试
- 健康检查测试

### 8.3 TDD 流程
1. 写测试 → 看失败
2. 写最小代码 → 看通过
3. 重构 → 保持通过
4. 重复

---

## 9. 未来扩展

### 9.1 计划中的模块
- **Brew Module**: 管理 Homebrew 包
- **NPM Module**: 管理全局 npm 包
- **Dotfiles Module**: 管理配置文件 symlink

### 9.2 增强功能
- 状态持久化（记录上次同步时间）
- 并行同步
- 配置验证
- 交互式初始化向导

---

## 10. 决策记录

### 10.1 使用 Go 而非 Rust
**决策**: 使用 Go
**理由**: 
- 更快的编译速度
- 更简单的构建流程
- 丰富的标准库
- 团队熟悉度

### 10.2 保留 repos.lock
**决策**: 向后兼容 repos.lock
**理由**:
- 平滑迁移路径
- 用户可以选择保持旧格式
- 降低切换成本

### 10.3 硬编码 link 映射
**决策**: repos.lock 使用硬编码 link 配置
**理由**:
- 保持 repos.lock 格式简单
- 新格式 dwell.yaml 支持自定义 links
- 满足向后兼容需求

---

## 11. 附录

### 11.1 术语表
- **Module**: 可管理的组件（git repo, brew bundle, etc.）
- **Sync**: 将模块同步到期望状态
- **Symlink**: 符号链接
- **Ref**: Git 引用（branch, tag, commit）

### 11.2 参考资料
- [urfave/cli documentation](https://cli.urfave.org/)
- [Git Go bindings](https://github.com/go-git/go-git)
- Original repos.lock format
