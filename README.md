# file_server
> a file_server with golang for practicing 

#### 简单的介绍
  这个项目是在linux环境下，通过golang实现的云盘存储功能的一个小项目。
  旨在练习golang语言，对数据库的连接及操作，以及后序会尝试加入redis的中间缓存机制，实现高并发下的云盘服务。
  
#### 项目的环境 Linux + Docker
  因为正好学习了docker容器，所以数据库是使用的docker中的mysql，并实现了主从复制功能，以保证数据持久性和安全性
