# 🚀 easypr — 高效管理多分支 PR 与 Cherry-Pick

[![Go Version](https://img.shields.io/github/go-mod/go-version/zhangqibupt/easypr)](https://go.dev/)

`easypr` 是一款 CLI 工具，用于自动化 **Git 多分支 cherry-pick 与 PR 创建流程**，帮助开发者高效同步代码到多个版本分支。

---

## ✨ 特性

- 一键创建主 PR 与多个 cherry-pick PR
- 自动分支创建、cherry-pick、推送与 PR 提交
- 增量同步最新提交至已存在的 PR
- 支持设置默认审查人与上游仓库
- 全交互式命令行，操作简单直观

---

## ⚙️ 安装

```bash
go install github.com/zhangqibupt/easypr@latest
```

请确保 `$GOPATH/bin` 已加入 `$PATH`，以便全局执行 `easypr`。

---

## 🧭 命令概览

| 命令 | 功能 | 示例 |
|------|------|------|
| `easypr create` | 一键创建目标 PR 与 cherry-pick PR | `easypr create` |
| `easypr sync` | 同步最新提交到所有 PR 分支 | `easypr sync` |
| `easypr config set-assignees` | 设置默认审查人 | `easypr config set-assignees alice bob` |
| `easypr config set-upstream` | 设置上游仓库（Fork 场景） | `easypr config set-upstream <url>` |

---

## 🧩 常见场景

### 1️⃣ 创建主 PR 并同步到多个版本分支

```bash
git checkout feature1
easypr create
```

交互式选择目标分支（如 `master`, `V_6_57`），工具将自动完成 cherry-pick、推送与 PR 创建。

---

### 2️⃣ 同步新提交到现有 Cherry-Pick PRs

```bash
git checkout feature1
easypr sync
```

选择要同步的目标分支，`easypr` 自动执行增量 cherry-pick 并更新远程 PR。

---

## ⚙️ 配置示例

设置上游仓库（适用于 Fork）：
```bash
easypr config set-upstream https://github.com/xx/xx.git
```

设置默认审查人：
```bash
easypr config set-assignees alice bob charlie
```

---

## ⚠️ 注意事项

- 发生冲突时需手动解决，工具会暂停并提示
- 当前支持 GitHub，后续计划支持 GitLab / Gitee

---

## 🤝 贡献

```bash
git checkout -b feature/AmazingFeature
git commit -m "Add AmazingFeature"
git push origin feature/AmazingFeature
```

提交 PR 即可参与贡献！

---

## 📄 许可证

MIT License，详见 [LICENSE](./LICENSE)。
