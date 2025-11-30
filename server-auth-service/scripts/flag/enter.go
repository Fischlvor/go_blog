package flag

import (
	"fmt"
	"os"

	"auth-service/pkg/global"

	"github.com/urfave/cli"
	"go.uber.org/zap"
)

// 定义 CLI 标志
var (
	sqlFlag = &cli.BoolFlag{
		Name:  "sql",
		Usage: "Initializes the structure of the MySQL database table.",
	}
)

// Run 执行基于命令行标志的相应操作
func Run(c *cli.Context) {
	// 检查是否设置了多个标志
	if c.NumFlags() > 1 {
		err := cli.NewExitError("Only one command can be specified", 1)
		global.Log.Error("Invalid command usage:", zap.Error(err))
		os.Exit(1)
	}

	// 根据不同的标志选择执行的操作
	switch {
	case c.Bool(sqlFlag.Name):
		if err := SQL(); err != nil {
			global.Log.Error("Failed to create table structure:", zap.Error(err))
			return
		} else {
			global.Log.Info("Successfully created table structure")
		}
	default:
		err := cli.NewExitError("unknown command", 1)
		global.Log.Error(err.Error(), zap.Error(err))
	}
}

// NewApp 创建并配置一个新的 CLI 应用程序
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "SSO Auth Service"
	app.Flags = []cli.Flag{
		sqlFlag,
	}
	app.Action = Run
	return app
}

// InitFlag 初始化并运行 CLI 应用程序
func InitFlag() {
	if len(os.Args) > 1 {
		app := NewApp()
		err := app.Run(os.Args)
		if err != nil {
			global.Log.Error("Application execution encountered an error:", zap.Error(err))
			os.Exit(1)
		}
		if os.Args[1] == "-h" || os.Args[1] == "-help" {
			fmt.Println("Displaying help message...")
		}
		os.Exit(0)
	}
}
