# student

基于 Go + chi + sqlc 的**模块化单体（Modular Monolith）**示例工程：业务按模块纵向切片，横切能力通过共享内核复用，持久化由 sqlc 生成类型安全查询。

## 目录结构

```
cmd/                      程序入口：加载配置、创建基础设施、装配模块
internal/
  student/                业务模块（限界上下文）：entity/service/repository/handler/dto/routes
  handler/                HTTP 处理器的共享基座（Base）
  service/                应用服务的共享基座（Base）
  repository/             仓储的共享基座（Base）+ 数据库连接
  platform/database/      连接池句柄与事务管理器（TxManager）
  platform/optional/      PATCH 三态字段类型（缺失 / null / 有值）
  server/                 HTTP 传输层：中间件 + 版本前缀 + 挂载各模块路由
  sqlc/                   sqlc 生成代码
pkg/                      基础设施包（app 生命周期、config、logger、middleware 等）
sqlc/                     schema.sql + queries/ + sqlc.yaml
```

依赖方向：

```
cmd ──► internal/server ──► 业务模块
cmd ──► 业务模块 ──► 共享内核（handler/service/repository）──► platform/*
```

## 共享内核四条规则

1. **共享内核不得 import 任何业务模块。**
   `internal/handler`、`internal/service`、`internal/repository` 只允许依赖标准库、`internal/platform/*` 与 `internal/sqlc`。
   一旦共享内核引入某个业务模块，所有模块都会被连带依赖，模块化即失效。

2. **业务模块之间禁止直接 import。**
   跨模块协作由**消费方**定义窄接口（例如 teacher 需要读取学生时，在 teacher 包内声明 `StudentReader`），在 `cmd` 装配处注入实现，避免循环依赖与整包耦合。

3. **模块只通过 Service 接口对外。**
   模块的 handler、repository、DTO 一律不导出给其它模块使用；模块内部实现细节可以自由演进。

4. **Base 只放横切能力。**
   `Base` 中不得出现任何业务概念；字段由构造函数统一注入，模块内只读、不得覆写。需要新能力时优先以字段组合进 Base，而不是为某个类型内嵌多个 Base。

## 事务

应用服务通过 `service.Base.Tx` 开启事务，仓储通过 `repository.Base.Q(ctx)` 自动加入当前事务，方法签名不需要暴露 `*sql.Tx`：

```go
err := s.Tx.Do(ctx, func(ctx context.Context) error {
	if err := s.studentRepo.Update(ctx, params); err != nil {
		return err
	}
	return s.teacherRepo.Create(ctx, params) // 同一个事务
})
```

`Q(ctx)` 在无事务时返回默认查询句柄，在有事务时返回 `Queries.WithTx(tx)`。

## 新增一个业务模块

1. 新建 `internal/<module>/`，包含 entity、service、repository、dto、handler、routes；
2. 在模块内实现 `Routes(chi.Router)`，自行注册 `/<module>` 下的路由；
3. 在 `cmd` 中装配并传给 `server.NewHTTPServer(addr, studentHandler, teacherHandler, ...)`；
4. `internal/server` 无需修改。

## 运行

```bash
go run ./cmd
```

配置文件：`config/config.json`；数据库默认位于 `data/student.db`。

## 代码生成

```bash
cd sqlc && sqlc generate
```
