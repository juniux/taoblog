# taoblog

Forked from 大佬的博客
<https://github.com/movsb/taoblog>

## Fork的目的

拿来即用的博客系统，学习GO

## 关于项目

"除了我以外，可能没有其他人能够很好地使用这套系统。" 尝试理解后能够改写

## 文件说明

文件名|文件描述
------|--------
admin/      | 后台目录（弃用）
client/     | 博客客户端
config/     | 配置模块
docker/     | 容器镜像
gateway/    | 网关接口层
modules/    | 公共模块
protocols/  | 协议定义
server/     | 博客后台
service/    | 服务实现
setup/      | 安装管理
run/        | 临时目录
themes/     | 主题目录
main.go     | 入口程序

## RoadMap

- 抄几个主题：

  - <http://www.templex.xyz/> 同王垠的博客
  - <https://www.tbfeng.com/> 天边风
  - <https://www.v2ex.com/t/561257> 你见过最简约美观的技术博客

## 联系原作者

- QQ: 191035066
- EM: chkesp@gmail.com

## Quick start

需要 Docker

```bash
$ make try
```
or
```bash
sudo make try
```

然后打开：<http://localhost:2564>。
