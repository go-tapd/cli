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
- [x] Shell completion 文档与安装说明

## 研发协作 API

### 需求

- [x] `tapd story create`
- [x] `tapd story view`
- [x] `tapd story list`
- [x] `tapd story count`
- [x] `tapd story update`
- [x] `tapd story batch-update`
- [x] `tapd story categories`
- [x] `tapd story categories count`
- [x] `tapd story changes`
- [x] `tapd story changes count`
- [x] `tapd story fields`
- [x] `tapd story field-labels`
- [x] `tapd story templates`
- [x] `tapd story template-fields`
- [x] `tapd story removed`
- [x] `tapd story related-bugs`
- [x] `tapd story related-test-cases`
- [x] `tapd story by-view`
- [x] `tapd story convert-ids`

### 缺陷

- [x] `tapd bug create`
- [x] `tapd bug copy`
- [x] `tapd bug view`
- [x] `tapd bug list`
- [x] `tapd bug count`
- [x] `tapd bug update`
- [x] `tapd bug batch-update`
- [x] `tapd bug changes`
- [x] `tapd bug changes count`
- [x] `tapd bug fields`
- [x] `tapd bug custom-field-settings`
- [x] `tapd bug field-labels`
- [x] `tapd bug templates`
- [x] `tapd bug template-fields`
- [x] `tapd bug removed`
- [x] `tapd bug related-stories`
- [x] `tapd bug links`
- [x] `tapd bug link`
- [x] `tapd bug unlink`
- [x] `tapd bug update-system-options`
- [x] `tapd bug by-view`
- [x] `tapd bug convert-ids`
- [x] 缺陷说明：`docs/bug.md`

### 迭代

- [x] `tapd iteration create`
- [x] `tapd iteration view`
- [x] `tapd iteration list`
- [x] `tapd iteration count`
- [x] `tapd iteration update`
- [x] `tapd iteration changes`
- [x] `tapd iteration fields`
- [x] `tapd iteration workitem-types`
- [x] `tapd iteration templates`
- [x] `tapd iteration template-fields`
- [x] `tapd iteration lock`
- [x] `tapd iteration unlock`

### 任务

- [x] `tapd task create`
- [x] `tapd task view`
- [x] `tapd task list`
- [x] `tapd task count`
- [x] `tapd task update`
- [x] `tapd task batch-update`
- [x] `tapd task changes`
- [x] `tapd task changes count`
- [x] `tapd task fields`
- [x] `tapd task removed`
- [ ] `tapd task by-view` SDK 暂未实现

### 测试

- [x] `tapd test-case create`
- [x] `tapd test-case list`
- [x] `tapd test-case count`
- [x] `tapd test-case update`
- [x] `tapd test-case categories`
- [x] `tapd test-case fields`
- [x] `tapd test-case results`
- [x] `tapd test-plan create`
- [x] `tapd test-plan list`
- [x] `tapd test-plan count`
- [x] `tapd test-plan update`
- [x] `tapd test-plan progress`
- [x] `tapd test-plan result`
- [x] `tapd test-plan related-bugs`
- [x] `tapd test-plan related-stories`
- [x] 测试说明：`docs/test.md`

### 发布

- [x] `tapd release create`
- [x] `tapd release view`
- [x] `tapd release list`
- [x] `tapd release count`
- [x] `tapd release update`
- [x] `tapd launch-form create`
- [x] `tapd launch-form list`
- [x] `tapd launch-form count`
- [x] `tapd launch-form fields`
- [x] `tapd launch-form templates`
- [x] `tapd launch-form logs`

### 源码

- [x] `tapd source commit add`
- [x] `tapd source commit list`
- [x] `tapd source commit objects`

### Wiki

- [x] `tapd wiki create`
- [x] `tapd wiki list`
- [x] `tapd wiki count`
- [x] `tapd wiki update`
- [x] `tapd wiki drawio`
- [x] `tapd wiki followers`
- [x] `tapd wiki followers count`
- [x] `tapd wiki permissions`
- [x] `tapd wiki tags`
- [x] `tapd wiki tags count`
- [x] `tapd wiki attachments count`

### 看板

- [x] `tapd board card create`
- [x] `tapd board card list`
- [x] `tapd board card update`
- [x] `tapd board columns`

### 评论

- [x] `tapd comment create`
- [x] `tapd comment list`
- [x] `tapd comment count`
- [x] `tapd comment update`

### 报表

- [x] `tapd report list`

### 附件

- [x] `tapd attachment list`
- [x] `tapd attachment download-url`
- [x] `tapd attachment image-url`
- [x] `tapd attachment document-url`

### 度量

- [x] `tapd measure life-times`

### 工时

- [x] `tapd timesheet create`
- [x] `tapd timesheet list`
- [x] `tapd timesheet count`
- [x] `tapd timesheet update`
- [x] `tapd timesheet delete`

### 项目

- [x] `tapd workspace view`
- [x] `tapd workspace users`
- [x] `tapd workspace roles`
- [x] `tapd workspace sub-workspaces`
- [x] `tapd workspace company-workspaces`
- [x] `tapd workspace participant-workspaces`
- [x] `tapd workspace add-member`
- [x] `tapd workspace update`
- [x] `tapd workspace custom-fields`
- [x] `tapd workspace documents`
- [x] `tapd workspace short-id convert`
- [x] `tapd workspace member-activity-log`
- [x] `tapd workspace calendar set-custom`
- [x] `tapd workspace calendar enable`
- [x] `tapd workspace calendar view-custom`
- [x] `tapd workspace calendar settings`

### 工作流

- [x] `tapd workflow last-steps`

### 配置

- [x] `tapd setting workspace`

### 标签

- [x] `tapd label create`
- [x] `tapd label list`
- [x] `tapd label count`
- [x] `tapd label update`

### 用户

- [x] `tapd user roles`

## Webhook

- [x] `tapd webhook serve`
- [x] `tapd webhook validate`
- [x] `tapd webhook inspect`
- [x] Webhook 说明：`docs/webhook.md`

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
