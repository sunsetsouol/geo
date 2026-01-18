# GEO (Generative Engine Optimization) 项目需求文档 (PRD)

## 1. 项目概览
GEO 项目旨在监控、分析和优化大模型（如 DeepSeek）中特定品牌的曝光情况。通过自动化脚本模拟用户提问，分析模型的回复及引用的链接，评估品牌在生成式引擎中的可见度和倾向性，从而指导优化策略。

## 2. 系统架构
系统分为两个主要部分：
- **监控端 (Monitoring Client)**: 基于 Playwright 的自动化脚本，负责与大模型交互。
- **服务端 (Server)**: 负责任务管理、提示词 (Prompt) 管理、数据接收、品牌曝光分析及 API 接口提供。

## 3. 功能需求

### 3.1 监控端 (Monitoring Client)
- **任务轮询**: 定期请求服务端获取处于 `pending` 状态的任务。
- **任务锁定**: 成功获取任务后，监控端进入任务执行状态，服务端应标记该任务为 `processing` 以防重复分配。
- **自动化交互**:
    - 使用 Playwright 驱动浏览器访问 DeepSeek 等模型。
    - 输入从服务端获取的 Prompt 并提交。
    - 实时监控页面，捕获大模型的完整响应文本。
    - 自动提取响应中包含的引用链接 (Citations/References)。
- **任务结果汇报**: 
    - 任务成功完成后，调用服务端 API 上报结果（文本 + 链接），并将任务状态更新为 `completed`。
    - 若执行失败（如网络、脚本错误），需上报失败原因，并将任务状态置为 `failed` 或回退到 `pending`。
- **异常处理**: 处理验证码、模型限流、页面结构变化等。

### 3.2 服务端 (Server)
- **Prompt 管理**:
    - 提供管理员 API 进行 Prompt 的 CRUD 操作。
    - 支持按标签或业务场景对 Prompt 进行分组。
- **任务生命周期管理**:
    - **分配逻辑**: 响应监控端的请求，分配优先级最高的 `pending` 任务。
    - **超时监控**: 监控长时间处于 `processing` 状态的任务，若超时则重置。
    - **结果处理**: 接收并存储监控端回传的模型原始数据。
- **品牌曝光分析 (核心)**:
    - **自动化打分**: 对 `Result` 中的文本进行多维度分析（品牌提及、推荐度、竞品对比）。
    - **链接分析**: 检查引用链接是否包含目标品牌官网或正面报道。
- **API 服务**:
    - 为监控端提供任务分发和结果回收接口。
    - 为管理后台提供 Prompt 管理和数据分析展示接口。

## 4. 任务执行流程

1. **监控端** -> 发起 `GET /api/client/task`: "有新任务吗？"
2. **服务端** -> 返回 `Task ID + Prompt`: "去问这个问题。" (并将任务置为 `processing`)
3. **监控端** -> 执行 Playwright 脚本: 打开 DeepSeek，提问，抓取结果。
4. **监控端** -> 发起 `POST /api/client/task/:id/report`: "任务做完了，这是结果。"
5. **服务端** -> 接收结果: 更新任务状态为 `completed`，触发异步品牌分析逻辑。

## 5. API 定义 (初步)

### 5.1 监控端接口
- `GET /api/client/task`: 获取当前待执行任务。
- `POST /api/client/task/:id/report`: 上报任务执行结果（`response_text`, `references`）。
- `POST /api/client/task/:id/error`: 上报任务失败信息。

### 5.2 Prompt 管理接口 (Admin)
- `GET /api/prompts`: 获取列表。
- `POST /api/prompts`: 创建。
- `PUT /api/prompts/:id`: 更新。
- `DELETE /api/prompts/:id`: 删除。

### 5.3 数据分析与展示接口
- `GET /api/analytics/overview`: 整体评分趋势。
- `GET /api/analytics/results`: 详细任务回复与得分详情。

## 6. 数据模型 (初步)

### 6.1 Prompt (提示词)
- `id`: 唯一标识
- `content`: 提示词内容
- `category`: 分类
- `created_at`: 创建时间

### 6.2 Task (任务)
- `id`: 唯一标识
- `prompt_id`: 关联的 Prompt
- `status`: 状态（pending, processing, completed, failed）
- `last_run`: 最近运行时间

### 6.3 Result (结果)
- `id`: 唯一标识
- `task_id`: 关联任务
- `response_text`: 模型回复全文
- `references`: 引用链接列表 (JSON)
- `brand_score`: 品牌曝光得分
- `analysis_report`: 详细分析报告 (JSON)

## 7. 技术栈建议
- **监控端**: Node.js + Playwright
- **服务端**: Node.js (NestJS/Express) 或 Python (FastAPI)
- **数据库**: PostgreSQL / MongoDB (用于结果存储)
- **分析引擎**: 可集成简单 NLP 库或调用大模型 API 进行二次打分分析。

## 8. 后续规划
- 增加更多大模型支持（ChatGPT, Claude, Gemini 等）。
- 增加代理池支持以应对反爬。
- 可视化看板 (Dashboard) 开发。
