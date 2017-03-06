# PCC API

> 服务器:
> user = pcc1:7000
> article = pcc2:7020
> action = pcc3:7010


###  API 列表

| 服务     | 方法     | 路由                  | 说明                     |
|:--------|:--------|:---------------------|:------------------------|
| user    | GET     | /api/v1/register     | 用户注册                  |
| article | GET     | /api/v1/articles     | 获取文章列表               |
| article | GET     | /api/v1/articles/:id | 获取某个文章的详情         |
| article | GET     | /api/v1/articles/:id/liked_count | 获取某个文章的点赞数         |
| article | POST    | /api/v1/articles     | 创建文章                 |
| article | PUT     | /api/v1/articles     | 更新文章                 |
| article | DELETE  | /api/v1/articles/:id | 删除文章                 |
| action  | GET     | /api/v1/articles/:id/liked_users | 获取点赞用户列表 |
| action  | GET     | /api/v1/articles/:id/is_liked    | 查看用户是否点赞 |
| action  | POST    | /api/v1/articles/:id/like        | 点赞           |
| action  | DELETE  | /api/v1/articles/:id/like        | 取消点赞        |
| action  | GET     | /api/v1/users/:id/followers      | 获取关注该用户的用户列表 |
| action  | GET     | /api/v1/users/:id/is_followed    | 该用户否被用户关注      |
| action  | POST    | /api/v1/users/:id/follow         | 关注该用户             |
| action  | DELETE  | /api/v1/users/:id/follow         | 取消关注该用户          |
