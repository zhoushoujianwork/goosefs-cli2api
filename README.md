<a name="itJfZ"></a>

# 项目背景

实现 api 接口管理 goosefs 集群，主要实现一下分布式缓存加载的数据预热更新的需求。<br />查看详细 API 文档：<br />[Apifox - 接口文档分享](https://apifox.com/apidoc/shared-78738557-618b-46ba-9a97-da1f90eeff26)
<a name="BeeEd"></a>

# 功能概述

- [x] docker-compose 容器部署启动；
- [x] 提交的任务支持多 Path, 支持任务带 `task_name`;
- [x] 完成的任务支持发送通知到钉钉群，异常也通知出来，有`task_name`的会以 task_name 为集合进行通报；注意钉钉机器人需要设置关键字：告警、通知；
- [x] 查询支持传入 `taskid`或者 `task_name`；
- [x] 在发起钉钉通知的时候，只告警有新数据加载的路径，没有的隐藏掉；
- [x] GooseFSForceLoad 该步骤执行的是先去 LoadMetadata，然后再去 DistributeLoad，这样彻底更新；
      <a name="ahcmA"></a>

# 如何部署

因为需要用到宿主机的 goosefs 环境，所以放弃了 docker 部署方案，采用 systemctl。<br />部署脚本：[https://github.com/zhoushoujianwork/goosefs-cli2api/blob/master/deploy.sh](https://github.com/zhoushoujianwork/goosefs-cli2api/blob/master/deploy.sh)<br />配置文件内容：[https://github.com/zhoushoujianwork/goosefs-cli2api/blob/master/config/config.yaml](https://github.com/zhoushoujianwork/goosefs-cli2api/blob/master/config/config.yaml)<br />配置文件支持当前运行目录，`/etc/goosefs-cli2api/`目录，日志文件`/opt/goosefs-cli2api/app.log`其他目录见 `config.yaml`，默认运行端口 8080。
<a name="pfpm4"></a>

# 举个例子

<a name="Edm4E"></a>

### 日常任务接口

支持三种任务：

1. GooseFSList，查询直接返回
2. GooseFSForceLoad，强制先执行 GooseFSLoadMetadata，然后再执行 GooseFSDistributeLoad；
3. GooseFSLoadMetadata，以任务的方式挂起，结果通过下面接口查询
4. GooseFSDistributeLoad，以任务的方式挂起，结果通过下面接口查询
   <a name="bdQGS"></a>

#### 查询缓存目录 GooseFSList

由于该请求不是挂起的后台执行，所以支持设置 timeout，默认 30 秒；

```bash
curl --location --request POST 'http://localhost:8080/api/v1/gfs' \
--header 'Content-Type: application/json' \
--data-raw '{
  "action": "GooseFSList",
  "path": [
    "/"
  ]
}'
```

<a name="fbMUq"></a>

#### GooseFSDistributeLoad

```bash
curl --location --request POST 'http://localhost:8080/api/v1/gfs' \
--header 'Content-Type: application/json' \
--data-raw '{
  "action": "GooseFSForceLoad",
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

<a name="ca6Jq"></a>

#### 发起分布式缓存目录请求 GooseFSDistributeLoad

```bash
curl --location --request POST 'http://localhost:8080/api/v1/gfs' \
--header 'Content-Type: application/json' \
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

<a name="yHWH4"></a>

#### 发起 loadMetadata 请求 GooseFSLoadMetadata

```bash
curl --location --request POST 'http://localhost:8080/api/v1/gfs' \
--header 'Content-Type: application/json' \
--data-raw '{
  "action": "GooseFSLoadMetadata",
  "path": [
    "/data-datalake-dataprod-bj-1251949819/deltalake/test"
  ]
}'

# 返回
[
    "204a65e4-4ad3-4f07-850c-b99077aa6d02",
    "160d7bf1-0da3-473d-b4c4-29f7c7d37e14"
]

```

<a name="BUBBw"></a>

### 任务查询接口

查询支持参数：<br />`task_id`：支持单个任务的查询<br />`test_task_name`：任务维度的查询，支持查询这一组任务的所有缓存任务的结果，返回整体任务状态<br />`action`：GooseFSForceLoad，GooseFSDistributeLoad，GooseFSLoadMetadata<br />`status`：success，failed，notallsuccess，running
<a name="FoopH"></a>

#### 查询任务输出

```bash
curl --location --request GET 'http://localhost:8080/api/v1/output?task_name=test_task_name'
```

<a name="ccgvf"></a>

#### 查询任务状态

```bash
curl --location --request GET 'http://localhost:8080/api/v1/status?task_name=test_task_name'
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
