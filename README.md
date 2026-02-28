# xq-agent (桌面智能体) 操作手册

## 项目简介

本项目是一个基于 Golang 开发的强大桌面智能体，旨在通过自然语言交互，帮助用户自动执行各类桌面任务。

它不仅集成了 **OpenAI 兼容的大模型**，还具备 **多渠道接入**、**浏览器自动化**、**文件操作**、**定时任务调度** 以及 **OpenClaw 技能扩展** 等高级能力。

现在，它还拥有一个现代化的 **GUI 桌面窗口**，提供类似 ChatGPT 的流畅对话体验。

---

## 核心功能概览

1.  **多模态交互**
    *   **GUI 桌面窗口**: 独立的桌面应用窗口，支持 Markdown 渲染、代码高亮和流式打字机输出。
    *   **思考过程可视化**: 支持显示模型推理过程（Reasoning Content）和工具调用状态（Tool Execution），让智能体的思考“看得见”。
    *   **Console**: 传统的命令行交互模式。
    *   **WeCom (企业微信)**: 支持发送通知消息到企业微信应用。
2.  **内置能力 (Tools)**
    *   **浏览器**: 打开网页、读取内容、网页截图。
    *   **文件系统**: 读取、写入、列出文件。
    *   **Shell**: 执行系统命令。
    *   **定时任务**: 通过自然语言添加、查看、删除定时任务。
3.  **技能扩展 (Skills)**
    *   完全兼容 **OpenClaw** 生态的 `SKILL.md` 格式。
    *   将技能文件夹放入 `skills/` 目录即可自动加载并使用。
    *   **可执行脚本支持**: 在 `skills/` 目录下创建子文件夹，放入 `SKILL.md` 和你的 Python/Shell 脚本。Agent 会自动识别脚本路径并调用 Shell 工具执行。

---

## 扩展外部技能

Agent 支持通过 `skills` 目录加载外部技能。

### 1. 目录结构
在 `agent.exe` 同级目录下创建 `skills` 文件夹：

```
agent.exe
config.yaml
skills/
  ├── my-python-skill/
  │   ├── SKILL.md
  │   └── script.py
  └── my-shell-skill/
      ├── SKILL.md
      └── deploy.sh
```

### 2. 编写技能 (SKILL.md)
在技能文件夹中创建 `SKILL.md`，Agent 会自动识别并加载。

```yaml
---
name: my-python-skill
description: 一个使用 Python 脚本处理数据的技能示例
---
To use this skill, run the following command:

python {{path}}/script.py --input "some input"
```

> **注意**: `{{path}}` 占位符暂未内置支持，但 Agent 会自动获取技能的绝对路径并注入到 Prompt 中（显示为 `Path: ...`），模型会根据此路径构造正确的命令。

### 3. 重启 Agent
添加新技能后，重启 `agent.exe` 即可生效。

---

## 快速开始

### 1. 环境准备

*   **Go 语言环境**: 建议 Go 1.21+。
*   **GCC 编译器**: 编译 Webview 需要 CGO 支持（Windows 推荐 TDM-GCC 或 MinGW）。
*   **Chrome 浏览器**: 用于浏览器自动化功能（Agent 会自动查找本地 Chrome）。

### 2. 配置说明

在项目根目录下找到 `config.yaml`（如不存在请参考以下内容创建）：

```yaml
llm:
  api_key: "sk-your-key-here"  # 填入你的 OpenAI 或兼容 API Key
  base_url: "https://api.openai.com/v1" # 可选，填入你的 API Base URL
  model: "gpt-4o"              # 使用的模型名称

channels:
  wecom:
    enabled: true              # 是否开启企业微信
    corp_id: "ww..."           # 企业 ID
    agent_id: 1000001          # 应用 AgentId
    secret: "..."              # 应用 Secret
```

### 3. 编译与运行

**编译带 GUI 的版本（隐藏控制台窗口）**:

```bash
# 1. 下载依赖
go mod tidy

# 2. 编译（隐藏控制台，仅显示 GUI）
go build -ldflags "-H windowsgui" -o agent.exe ./cmd/agent

# 3. 运行
./agent.exe
```

> **提示**: 如果你需要查看详细的调试日志，可以去掉 `-ldflags "-H windowsgui"` 参数进行编译，这样运行时会同时保留控制台窗口。

---

## 企业微信接入指南

要让 Agent 通过企业微信发送消息，请按照以下步骤操作：

1.  **注册企业微信**: 访问 [企业微信官网](https://work.weixin.qq.com/) 注册企业。
2.  **创建应用**:
    *   进入 [管理后台](https://work.weixin.qq.com/wework_admin/frame) -> **应用管理**。
    *   点击 **创建应用**，填写名称（如 "桌面智能体"）和 Logo。
    *   创建成功后，记录 **AgentId** and **Secret**。
3.  **获取 CorpID**:
    *   进入 **我的企业** -> **企业信息**，复制底部的 **企业ID (CorpID)**。
4.  **修改配置**:
    *   打开 `config.yaml`。
    *   将 `channels.wecom.enabled` 设置为 `true`。
    *   填入 `corp_id`, `agent_id`, `secret`。
5.  **重启 Agent**:
    *   Agent 启动时会自动获取 AccessToken，如果成功会显示 `[WeCom] Successfully connected`。

**注意**: 目前仅支持 **发送消息** (Agent -> 微信)。接收微信消息需要配置公网回调服务器，暂未内置支持。

---

## 功能使用指南

### 1. 基础对话
启动后，在桌面窗口的输入框直接输入你的需求。Agent 支持流式输出，你会看到回复像打字机一样逐字出现。

> **用户**: 你好，介绍一下你自己。
> **Agent**: 你好！我是你的桌面智能助手...

### 2. 浏览器自动化
Agent 内置了基于 `chromedp` 的浏览器控制能力。

*   **访问网页**:
    > "帮我看看现在 Hacker News 的头条是什么？"
    > (Agent 会调用 `browser_open` 读取网页内容并总结)

*   **网页截图**:
    > "打开百度首页并截图保存为 baidu.png"
    > (Agent 会调用 `browser_screenshot`，截图将保存在当前目录)

### 3. 文件操作
Agent 可以帮你管理本地文件。

*   **读取文件**: "读取一下当前目录下的 config.yaml 文件内容。"
*   **写入文件**: "帮我创建一个 hello.txt，内容是 'Hello World'。"
*   **列出文件**: "看看 skills 目录下有哪些文件？"

### 4. 定时任务 (Cron)
Agent 内置了 Cron 调度器，你可以用自然语言管理任务。

*   **添加任务**:
    > "每天早上 9 点提醒我开早会。"
    > "每隔 10 秒钟告诉我一次现在的时间。"
*   **查看任务**: "列出当前所有的定时任务。"
*   **删除任务**: "删除 ID 为 2 的那个任务。"

### 5. 扩展技能 (Skills)
本项目支持加载外部技能，兼容 OpenClaw 规范。

*   **安装技能**: 将包含 `SKILL.md` 的技能文件夹放入 `skills/` 目录。
    *   例如：`skills/demo/SKILL.md`
*   **使用技能**: Agent 启动时会自动读取 `SKILL.md` 中的描述。
    *   如果安装了 `demo` 技能，你可以说："运行 demo skill 的问候命令。"
    *   Agent 会根据文档自动执行对应的 CLI 命令（如 `echo ...`）。

---

## 目录结构说明

*   `cmd/agent/`: 程序入口。
*   `internal/core/`: Agent 核心逻辑（LLM 交互、工具分发）。
*   `internal/channels/`: 渠道层（Webview GUI, Console, Telegram, WeCom）。
*   `internal/tools/`: 内置工具实现（Browser, File, Shell）。
*   `internal/cron/`: 定时任务管理。
*   `internal/skills/`: OpenClaw 技能管理器。
*   `skills/`: **用户技能目录**，存放外部技能。

## 常见问题

1.  **浏览器截图失败？**
    *   请确保本地安装了 Chrome 浏览器。
    *   如果是服务器环境，请确保安装了相关依赖库。
    *   日志中会显示具体错误，如超时（Timeout），请检查网络连接。

2.  **LLM 响应慢？**
    *   检查网络连接。
    *   在 `config.yaml` 中配置国内可访问的 `base_url`。

3.  **流式输出乱码？**
    *   Windows PowerShell/CMD 默认编码可能导致问题，建议使用 Windows Terminal 或设置 `chcp 65001`。
    *   **推荐使用 GUI 模式**，可完美解决乱码问题。
