# CinaSeek - 轻量级 Ubuntu 虚拟机引擎

> 基于 [Multipass](https://github.com/canonical/multipass) 深度定制的自主品牌虚拟机管理引擎

## 品牌说明
CinaSeek 是 CinaRoom 平台的核心虚拟机引擎，基于 Canonical Multipass 技术构建，
以自主品牌 Cinaseek 运营，遵循 GPLv3 开源协议。

## 架构关系
CinaGroup → CinaSeek（核心产品/LLM聚合）→ CinaRoom（虚拟机管理平台）→ CinaSeek（虚拟机引擎）

## 与上游 Multipass 的区别
- 品牌化：客户端/守护进程名称改为 cinaseek/cinaseekd
- 环境变量前缀：MULTIPASS_* → CINASEEK_*
- protobuf 命名空间：multipass → cinaseek
- 深度集成 CinaRoom 管理平台
- OpenClaw 专属优化

## 许可证
本项目基于 GPLv3 协议开源，原版权归属 Canonical, Ltd.
修改部分版权归 CinaGroup 所有。
