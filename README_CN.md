[English](https://github.com/nefukadia/hot-cache/blob/dev/README.md) |
[简体中文](https://github.com/nefukadia/hot-cache/blob/dev/README_CN.md)

# HotCache

## 介绍
虽然缓存服务器可以减轻对关系数据库的访问，但对于热点数据（短时间内访问次数巨大）来说， 不断访问缓存数据库也是一个不小的时间开销。
因此HotCache横空出世，它可以通过sdk来实现本地数据缓存，并在一定程度上保证数据的一致性。

## 特性
- sdk本地缓存数据
- 服务端可持久化
- 支持集群模式
- 主从复制保证高可用