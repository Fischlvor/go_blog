@echo off

:: 创建根目录结构
mkdir server\api
mkdir server\assets
mkdir server\config
mkdir server\core
mkdir server\flag
mkdir server\global
mkdir server\initialize
mkdir server\log
mkdir server\middleware
mkdir server\model
mkdir server\router
mkdir server\service
mkdir server\task
mkdir server\uploads
mkdir server\utils

:: 创建 model 子目录
mkdir server\model\appTypes
mkdir server\model\database
mkdir server\model\elasticsearch
mkdir server\model\other
mkdir server\model\request
mkdir server\model\response

:: 创建 utils 子目录
mkdir server\utils\hotSearch
mkdir server\utils\upload

echo 目录结构已创建！
pause


