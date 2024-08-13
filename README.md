<a name="itJfZ"></a>

# 项目背景

实现 api 接口管理 goosefs 集群，主要实现一下分布式缓存加载的数据预热更新的需求。<br />查看详细 API 文档：<br />[Apifox - 接口文档分享](https://apifox.com/apidoc/shared-78738557-618b-46ba-9a97-da1f90eeff26)
<a name="BeeEd"></a>

# 功能概述

- [x] docker-compose 容器部署启动；
- [x] 提交的任务支持多 Path, 支持任务带 `task_name`;
- [x] 完成的任务支持发送通知到钉钉群，异常也通知出来，有`task_name`的会以 task_name 为集合进行通报；注意钉钉机器人需要设置关键字：告警、通知；
- [x] 查询支持传入 `taskid`或者 `task_name`;
      <a name="pfpm4"></a>

# 举个例子

<a name="Edm4E"></a>

### 日常任务接口

<a name="bdQGS"></a>

#### 查询缓存目录 GooseFSList

由于该请求不是挂起的后台执行，所以支持设置 timeout，默认 30 秒；

```bash
curl --location --request POST 'http://localhost:8080/api/v1/gfs' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Connection: keep-alive' \
--data-raw '{
  "action": "GooseFSList",
  "path": [
    "/"
  ]
}'
```

<a name="ca6Jq"></a>

#### 发起分布式缓存目录请求 GooseFSDistributeLoad

```bash
curl --location --request POST 'http://localhost:8080/api/v1/gfs' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Connection: keep-alive' \
--data-raw '{
  "action": "GooseFSDistributeLoad",
  "task_name": "test_task_name",
  "path": [
    "/data-datalake-dataprod-bj-1251949819/deltalake/npd_temp.db"
  ]
}'

# 返回任务 ID
[
    "cd3fe749-4e32-4415-85df-57c3ddbea1b4"
]
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

#### 发起 loadMetadata 请求 GooseFSLoadMetadata

> 参考 1.2 ，唯一不一样的是请求体 action 为 1；

```json
{
  "action": 1,
  "path": "/tmp/"
}
```

<a name="BUBBw"></a>

### 任务查询接口

查询支持参数：<br />`task_id`：支持单个任务的查询<br />`test_task_name`：任务维度的查询，支持查询这一组任务的所有缓存任务的结果，返回整体任务状态
<a name="FoopH"></a>

#### 查询任务输出

```bash
curl --location --request GET 'http://localhost:8080/api/v1/output?task_name=test_task_name' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Accept: */*' \
--header 'Connection: keep-alive'
```

<a name="ccgvf"></a>

#### 查询任务状态

```bash
curl --location --request GET 'http://localhost:8080/api/v1/status?task_name=test_task_name' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Accept: */*' \
--header 'Connection: keep-alive'
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
