[English](README.md) | [简体中文](README_CN.md)

# HotCache

## introduce
Although the cache server can reduce the access to the relational database, for hot data (large number of accesses in a short time), continuous access to the cache database is also a significant time cost. 
Therefore, HotCache came into being. It can realize local data caching through SDK and ensure data consistency to a certain extent.

## feature
- SDK local cache data
- Persistence on the server side
- Support cluster mode
- Master-slave replication ensures high availability