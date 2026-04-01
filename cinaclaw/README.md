# CinaClaw - 轻量级 Ubuntu 虚拟机引擎

> 基于 [Multipass](https://github.com/canonical/multipass) 深度定制的自主品牌虚拟机管理引擎

## 品牌说明
CinaClaw 是 CinaSeek 平台的核心虚拟机引擎，基于 Canonical Multipass 技术构建，
遵循 GPLv3 开源协议。

## 架构关系
```
CinaGroup 技术生态
├── CinaToken         # LLM 聚合平台 + OAuth 认证
└── CinaSeek          # 虚拟机远程管理平台
    └── CinaClaw      # VM 引擎（基于 Multipass fork）
```

## 与上游 Multipass 的区别
- 品牌化：客户端/守护进程名称改为 cinaclaw/cinaclawd
- 环境变量前缀：MULTIPASS_* → CINACLAWS_*
- protobuf 命名空间：multipass → cinaclaw
- 深度集成 CinaSeek 管理平台
- OpenClaw 专属优化

## 许可证
本项目基于 GPLv3 协议开源，原版权归属 Canonical, Ltd.
修改部分版权归 CinaGroup 所有。
