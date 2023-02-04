[English](https://github.com/nefukadia/hot-cache/blob/dev/README.md) | 
[简体中文](https://github.com/nefukadia/hot-cache/blob/dev/README_CN.md)

# HotCache

## Introduce
Although the cache server can reduce the access to the relational database, for hot data (large number of accesses in a short time), continuous access to the cache database is also a significant time cost. 
Therefore, HotCache came into being. It can realize local data caching through SDK and ensure data consistency to a certain extent.

## Feature
- SDK local cache data
- Persistence on the server side
- Support cluster mode
- Master-slave replication ensures high availability