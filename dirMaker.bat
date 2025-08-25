@echo off

:: 创建根目录结构
mkdir server\cmd\server
mkdir server\internal\api
mkdir server\internal\service
mkdir server\internal\model
mkdir server\internal\middleware
mkdir server\internal\router
mkdir server\internal\initialize
mkdir server\internal\task
mkdir server\pkg\config
mkdir server\pkg\utils
mkdir server\pkg\global
mkdir server\pkg\core
mkdir server\configs
mkdir server\scripts
mkdir server\log
mkdir server\uploads

:: 创建 internal/model 子目录
mkdir server\internal\model\appTypes
mkdir server\internal\model\database
mkdir server\internal\model\elasticsearch
mkdir server\internal\model\other
mkdir server\internal\model\request
mkdir server\internal\model\response

:: 创建 internal/service 子目录
mkdir server\internal\service\ai

:: 创建 pkg/utils 子目录
mkdir server\pkg\utils\hotSearch
mkdir server\pkg\utils\upload

:: 创建 scripts 子目录
mkdir server\scripts\flag

echo 目录结构已创建！
echo.
echo 新的项目结构：
echo ├── cmd/server/          # 应用程序入口点
echo ├── internal/            # 私有应用程序和库代码
echo │   ├── api/            # API层
echo │   ├── service/        # 业务逻辑层
echo │   ├── model/          # 数据模型
echo │   ├── middleware/     # 中间件
echo │   ├── router/         # 路由
echo │   ├── initialize/     # 初始化
echo │   └── task/           # 任务
echo ├── pkg/                # 可以被外部应用程序使用的库代码
echo │   ├── config/         # 配置管理
echo │   ├── utils/          # 工具函数
echo │   ├── global/         # 全局变量
echo │   └── core/           # 核心功能
echo ├── configs/            # 配置文件
echo ├── scripts/            # 脚本文件
echo ├── log/                # 日志文件
echo └── uploads/            # 上传文件
echo.
pause


