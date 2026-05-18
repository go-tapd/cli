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

- [x] `tapd source commit add`
- [x] `tapd source commit list`
- [x] `tapd source commit objects`

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
