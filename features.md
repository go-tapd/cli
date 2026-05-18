# TODO

This file tracks TAPD CLI feature coverage.

The SDK coverage ledger lives in the sibling SDK repository:
`../tapd/features.md`.

> [!NOTE]
> This file is command-oriented. A checked item means the CLI exposes a user-facing
> command, not merely that the SDK implements the API.
>
> When adding a command, confirm the matching typed API already exists in
> `github.com/go-tapd/tapd` before wiring it into the CLI.

```
开发说明：

1、CLI 命令应优先复用 go-tapd/tapd SDK 的 typed service、request 和 response，不在 CLI 内手写 HTTP 请求。
2、列表和详情类命令默认支持 table 输出，并支持 --format json。
3、需要项目上下文的命令统一使用 --workspace-id / -w。
4、只读命令优先于写操作命令；写操作命令应先明确 flag 契约、必填字段和输出形态。
5、SDK features.md 中标注为 “AI 实现，未人工验证” 的接口，接入 CLI 前应额外谨慎。
```

## 基础能力

- [x] `tapd login`
- [x] `tapd login --auth-method pat`
- [x] `tapd login --auth-method basic`
- [x] `tapd login --workspace-id <id>` 登录时验证项目访问权限
- [x] `tapd auth status`
- [x] `tapd auth logout`
- [x] `--config <path>` 指定配置文件
- [x] `--base-url <url>` 覆盖 TAPD API 地址
- [x] `--format table|json`
- [ ] Shell completion 文档与安装说明

## 研发协作 API

### 需求

- [ ] `tapd story create`
- [ ] `tapd story view`
- [x] `tapd story list`
- [x] `tapd story count`
- [ ] `tapd story update`
- [ ] `tapd story batch-update`
- [ ] `tapd story categories`
- [ ] `tapd story categories count`
- [ ] `tapd story changes`
- [ ] `tapd story changes count`
- [x] `tapd story fields`
- [ ] `tapd story field-labels`
- [ ] `tapd story templates`
- [ ] `tapd story template-fields`
- [ ] `tapd story removed`
- [ ] `tapd story related-bugs`
- [ ] `tapd story related-test-cases`
- [ ] `tapd story by-view`
- [ ] `tapd story convert-ids`

### 缺陷

- [ ] `tapd bug create`
- [ ] `tapd bug view`
- [x] `tapd bug list`
- [x] `tapd bug count`
- [ ] `tapd bug update`
- [ ] `tapd bug batch-update`
- [ ] `tapd bug changes`
- [ ] `tapd bug changes count`
- [x] `tapd bug fields`
- [ ] `tapd bug field-labels`
- [ ] `tapd bug templates`
- [ ] `tapd bug template-fields`
- [ ] `tapd bug removed`
- [ ] `tapd bug related-stories`
- [ ] `tapd bug link`
- [ ] `tapd bug unlink`
- [ ] `tapd bug by-view`
- [ ] `tapd bug convert-ids`

### 迭代

- [ ] `tapd iteration create`
- [ ] `tapd iteration view`
- [ ] `tapd iteration list`
- [ ] `tapd iteration count`
- [ ] `tapd iteration update`
- [ ] `tapd iteration changes`
- [ ] `tapd iteration fields`
- [ ] `tapd iteration workitem-types`
- [ ] `tapd iteration templates`
- [ ] `tapd iteration template-fields`
- [ ] `tapd iteration lock`
- [ ] `tapd iteration unlock`

### 任务

- [ ] `tapd task create`
- [ ] `tapd task view`
- [x] `tapd task list`
- [x] `tapd task count`
- [ ] `tapd task update`
- [ ] `tapd task batch-update`
- [ ] `tapd task changes`
- [ ] `tapd task changes count`
- [x] `tapd task fields`
- [ ] `tapd task removed`
- [ ] `tapd task by-view` SDK 暂未实现

### 测试

- [ ] `tapd test-case create`
- [ ] `tapd test-case list`
- [ ] `tapd test-case count`
- [ ] `tapd test-case update`
- [ ] `tapd test-case categories`
- [ ] `tapd test-case fields`
- [ ] `tapd test-case results`
- [ ] `tapd test-plan create`
- [ ] `tapd test-plan list`
- [ ] `tapd test-plan count`
- [ ] `tapd test-plan update`
- [ ] `tapd test-plan progress`
- [ ] `tapd test-plan result`
- [ ] `tapd test-plan related-bugs`
- [ ] `tapd test-plan related-stories`

### 发布

- [ ] `tapd release create`
- [ ] `tapd release view`
- [ ] `tapd release list`
- [ ] `tapd release count`
- [ ] `tapd release update`
- [ ] `tapd launch-form create`
- [ ] `tapd launch-form list`
- [ ] `tapd launch-form count`
- [ ] `tapd launch-form fields`
- [ ] `tapd launch-form templates`
- [ ] `tapd launch-form logs`

### 源码

- [ ] `tapd source commit add`
- [ ] `tapd source commit list`
- [ ] `tapd source commit objects`

### Wiki

- [ ] `tapd wiki create`
- [ ] `tapd wiki list`
- [ ] `tapd wiki count`
- [ ] `tapd wiki update`
- [ ] `tapd wiki drawio`
- [ ] `tapd wiki followers`
- [ ] `tapd wiki followers count`
- [ ] `tapd wiki permissions`
- [ ] `tapd wiki tags`
- [ ] `tapd wiki tags count`
- [ ] `tapd wiki attachments count`

### 看板

- [ ] `tapd board card create`
- [ ] `tapd board card list`
- [ ] `tapd board card update`
- [ ] `tapd board columns`

### 评论

- [ ] `tapd comment create`
- [ ] `tapd comment list`
- [ ] `tapd comment count`
- [ ] `tapd comment update`

### 报表

- [ ] `tapd report list`

### 附件

- [ ] `tapd attachment list`
- [ ] `tapd attachment download-url`
- [ ] `tapd attachment image-url`
- [ ] `tapd attachment document-url`

### 度量

- [ ] `tapd measure life-times`

### 工时

- [ ] `tapd timesheet create`
- [ ] `tapd timesheet list`
- [ ] `tapd timesheet count`
- [ ] `tapd timesheet update`
- [ ] `tapd timesheet delete`

### 项目

- [x] `tapd workspace view`
- [x] `tapd workspace users`
- [ ] `tapd workspace roles`
- [ ] `tapd workspace sub-workspaces`
- [ ] `tapd workspace company-workspaces`
- [ ] `tapd workspace participant-workspaces`
- [ ] `tapd workspace add-member`
- [ ] `tapd workspace update`
- [ ] `tapd workspace custom-fields`
- [ ] `tapd workspace documents`
- [ ] `tapd workspace short-id convert`
- [ ] `tapd workspace member-activity-log`
- [ ] `tapd workspace calendar set-custom`
- [ ] `tapd workspace calendar enable`
- [ ] `tapd workspace calendar view-custom`
- [ ] `tapd workspace calendar settings`

### 工作流

- [ ] `tapd workflow last-steps`

### 配置

- [ ] `tapd setting workspace`

### 标签

- [ ] `tapd label create`
- [ ] `tapd label list`
- [ ] `tapd label count`
- [ ] `tapd label update`

### 用户

- [ ] `tapd user roles`

## Webhook

- [ ] `tapd webhook serve`
- [ ] `tapd webhook validate`
- [ ] `tapd webhook inspect`

## 轻协作 API

SDK 当前对轻协作 API 覆盖较少，CLI 暂不优先接入。

### 工作项

- [ ] `tapd lite workitem create`
- [ ] `tapd lite workitem list`
- [ ] `tapd lite workitem count`
- [ ] `tapd lite workitem update`

### 空间

- [ ] `tapd lite space view`
- [ ] `tapd lite space create`
- [ ] `tapd lite space members`
- [ ] `tapd lite space add-member`

### 评论

- [ ] `tapd lite comment create`
- [ ] `tapd lite comment list`
- [ ] `tapd lite comment count`
- [ ] `tapd lite comment update`

### 附件

- [ ] `tapd lite attachment upload`
- [ ] `tapd lite attachment list`
- [ ] `tapd lite attachment download-url`

### 应用集成-工蜂

- [ ] `tapd lite git branch link`
- [ ] `tapd lite git branch unlink`
- [ ] `tapd lite git branch workitems`
- [ ] `tapd lite git commit add`
- [ ] `tapd lite git commit objects`
