# GEO (Generative Engine Optimization) 项目需求文档 (PRD)

## 1. 项目概览
GEO 项目旨在监控、分析和优化大模型（如 DeepSeek）中特定品牌的曝光情况。通过自动化脚本模拟用户提问，分析模型的回复及引用的链接，评估品牌在生成式引擎中的可见度和倾向性。同时，系统通过**生成针对性的品牌文章并发布到网络**，增加大模型训练或联网搜索时的品牌相关语料，从而主动优化品牌在 AI 模型中的曝光表现。

## 2. 系统架构
系统分为三个主要部分：
- **监控端 (Monitoring Client)**: 基于 Playwright 的自动化脚本，负责与大模型交互，采集品牌曝光数据。
- **服务端 (Server)**: 
    - **任务与数据管理**: 负责 Prompt 管理、任务调度、结果存储。
    - **分析引擎**: 对采集到的回复进行自动化评分与分析。
    - **内容生成模块 (AGO - Article Generation & Optimization)**: 调用大模型生成利于 SEO/GEO 的文章，并管理发布流程。
- **数据库**: MySQL 存储提示词、任务记录、分析结果及生成的优化文章。

## 3. 功能需求

### 3.1 监控端 (Monitoring Client)
- **任务轮询**: 定期请求服务端获取处于 `pending` 状态的任务。
- **自动化交互**: 使用 Playwright 驱动浏览器访问 DeepSeek 等模型，提交 Prompt 并截取完整响应。
- **链接提取**: 自动提取响应中包含的引用链接 (Citations/References)。
- **任务结果汇报**: 上报文本、链接及原始 HTML 输出。

### 3.2 服务端 (Server) - 监控与分析
- **Prompt 管理**: 支持 Prompt 的分类与 CRUD。
- **品牌曝光分析 (核心)**:
    - **自动化打分**: 对回复文本进行提及率、正面性、对比优势等维度的评分。
    - **引用分析**: 检查引用链接的权重，判断是否包含官方或正面渠道。

### 3.3 内容优化 (Content Optimization) - 新增
- **文章生成**:
    - 根据监控分析出的“品牌弱点”或“竞品优势”，自动调用 LLM 生成高质量的品牌软文。
    - 文章需包含目标关键词、品牌信息，并模仿权威媒体或真实用户的口吻。
- **发布管理**:
    - **多平台适配**: 支持配置不同社交平台或博客网站的发布接口（或模拟点击发布）。
    - **链接反馈**: 记录发布后的文章 URL，以便后续监控这些链接是否被 AI 引用。

## 4. 任务执行流程
1. **分析**: 服务端分析现有数据，发现品牌在某些提问下的曝光不足。
2. **生成**: 内容生成模块生成 3-5 篇针对性的优化文章。
3. **发布**: 系统自动或辅助发布文章到互联网。
4. **监控**: 监控端定期执行相关任务，验证新发布的文章是否被 AI 引用，以及品牌评分是否提升。

## 5. API 定义 (初步)

### 5.1 监控端接口
- `GET /api/client/task`: 获取任务。
- `POST /api/client/task/:id/report`: 上报结果（含 `raw_output`）。

### 5.2 优化与生成接口 (Admin)
- `POST /api/optimize/generate-article`: 触发文章生成任务。
- `GET /api/optimize/articles`: 获取已生成的文章列表。
- `POST /api/optimize/publish`: 记录/触发发布任务。

## 6. 数据模型 (初步)

### 6.1 Prompt (提示词)
- `id`, `content`, `category`

### 6.2 Task & Result
- `tasks`: `id`, `prompt_id`, `status`, `last_run`
- `results`: `id`, `task_id`, `response_text`, `raw_output`, `brand_score`, `analysis_report`
- `citations`: `id`, `task_id`, `url`, `title` (独立存放)

### 6.3 Optimization Article (优化文章) - 新增
- `id`: 唯一标识
- `brand_id`: 目标品牌
- `title`: 文章标题
- `content`: 文章正文
- `target_keywords`: 目标关键词
- `publish_status`: 发布状态（pending, published）
- `published_url`: 发布后的外链地址
- `created_at`: 生成时间

## 7. 技术栈建议
- **监控端**: Python + Playwright (同步 API)
- **服务端后端**: Go (trpc-go 框架)
- **服务端前端**: Vue 3 (Vite + Pinia)
- **数据库**: MySQL 8.0+

## 8. 后续规划
- 建立**品牌语料库**，确保持续输出高质量内容。
- 对接更多自媒体平台 API。
- 增加**GEO 效果追踪图表**，可视化展示优化前后的品牌权重变化。
