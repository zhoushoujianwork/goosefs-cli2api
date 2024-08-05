<a name="E1nLJ"></a>

## 项目背景

实现 api 接口管理 goosefs 集群，主要实现一下分布式缓存加载的数据预热更新的需求。<br />查看详细 API 文档：<br />[Apifox - 接口文档分享](https://apifox.com/apidoc/shared-78738557-618b-46ba-9a97-da1f90eeff26)

<a name="JnCZa"></a>

## 举个例子

<a name="Edm4E"></a>

### 日常任务接口

<a name="bdQGS"></a>

#### 查询缓存目录

由于该请求不是挂起的后台执行，所以支持设置 timeout，默认 30 秒；

```bash

$ curl --location --request POST 'http://localhost:8080/api/v1/gfs' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8080' \
--header 'Connection: keep-alive' \
--data-raw '{
    "action": 2,
    "path": "/xxx/deltalake",
    "timeout": 20
}'
```

<a name="ca6Jq"></a>

#### 发起分布式缓存目录请求

```bash
# 执行内置的 goosefs 命令，包括 0:distribute_load 1:load_metadata，返回 task_id，可以通过 task_id 获取执行状态或者输出

$ curl --request POST 'http://localhost:8080/api/v1/gfs' \
--header 'Content-Type: application/json' \
--header 'Host: localhost:8080' \
--data-raw '{
    "action": 0,
    "path": "/tmp/"
}'

# 返回任务 ID
683877b1-d293-4b89-b5c0-4f5742502dd5

```

```bash
# 因为任务是后台挂起的，可能长时间执行中
# 注意将 taskID 放到 Path中
$ curl --request GET 'http://localhost:8080/api/v1/status/683877b1-d293-4b89-b5c0-4f5742502dd5'

# 返回
{
    "id": "683877b1-d293-4b89-b5c0-4f5742502dd5",
    "status": "exit status 0"
}
```

```bash
# 同样的 url path加入 taskID
$ curl --request GET 'http://localhost:8080/api/v1/output/683877b1-d293-4b89-b5c0-4f5742502dd5'
```

<a name="yHWH4"></a>

#### 发起 loadMetadata 请求

> 参考 1.2 ，唯一不一样的是请求体 action 为 1；

```json
{
  "action": 1,
  "path": "/tmp/"
}
```

<a name="ucUTV"></a>

### 管理接口

<a name="kBTjA"></a>

#### 查看集群缓存状态

```bash
curl --request GET 'http://localhost:8080/api/v1/gfs/report'

# 返回结果
"GooseFS cluster summary: \n    Master Address: x.x.x.x:9200\n    Web Port: 9201\n    Rpc Port: 9200\n    Started: 07-30-2024 15:05:31:578\n    Uptime: 5 day(s), 23 hour(s), 25 minute(s), and 37 second(s)\n    Version: 1.4.5.3\n    Safe Mode: false\n    Zookeeper Enabled: false\n    Live Workers: 10\n    Lost Workers: 0\n    Total Capacity: 117.19TB\n        Tier: SSD  Size: 117.19TB\n    Used Capacity: 117.81GB\n        Tier: SSD  Size: 117.81GB\n    Free Capacity: 117.07TB\n"
```
