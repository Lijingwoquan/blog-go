# Go 博客管理系统

这是一个使用Go语言开发的后端博客管理系统。它采用了以下技术:

- **Gin框架**作为Web框架
- **JWT(JSON Web Token)** 实现用户认证
- **SQLx** 用于与数据库交互
- **Redis** 用于缓存数据
- **HTTPS** 提供安全的加密通信

## 功能特性

- 创建、更新、删除和查看博客文章
- 对文章进行分类和标签管理
- 管理员可以管理博客内容
- 使用Redis缓存提高性能
- 使用HTTPS保护通信安全

## 技术栈

- **后端**: Go, Gin, SQLx, Redis
- **认证**: JSON Web Tokens (JWT)
- **数据库**: MySQL
- **缓存**: Redis
- **部署**: Docker

## 快速开始

1. 克隆仓库:
   git clone https://github.com/Lijingwoquan/blog-go.git
2. 配置环境变量:
    - 在 **./config/config.yaml**下修改自己的环境配置
    - 在**docker-compose** 下修改 **mysql** 密码(如需部署)
    - 在目录 **/config/redis.config**下修改**redis**信息(如需部署)
3. 运行应用程序:

```ssh
go build main.go
go run main.go
```

## 部署

- 该项目需要配置证书，需要在ssl文件夹下配置
> 若无证书可以在 `main.go` 中使用运行 `err := r.Run(port)`
- 该应用程序可以使用Docker进行部署。

```ssh
docker-compose up -d --build
```

## 贡献

如果你发现任何问题或有改进建议,欢迎提交issue或发起Pull Request。
