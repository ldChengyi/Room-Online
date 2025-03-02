package core

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// HandleSignals 创建带有优雅关闭能力的上下文
// 返回值：带取消信号的上下文，用于协调服务关闭
func HandleSignals() context.Context {
    // 创建可取消的上下文
    ctx, cancel := context.WithCancel(context.Background())

    // 监听系统信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, 
        syscall.SIGINT,  // Ctrl+C
        syscall.SIGTERM, // Kubernetes等编排系统的终止信号
        syscall.SIGQUIT) // 优雅退出

    go func() {
        // 第一次接收信号 - 启动优雅关闭
        sig := <-sigCh
        log.Printf("Received signal: %v, initiating shutdown...", sig)
        
        // 执行取消操作（触发ctx.Done()）
        cancel()

        // 第二次接收信号 - 强制立即退出
        <-sigCh
        log.Printf("Force shutdown requested!")
        os.Exit(1)
    }()

    return ctx
}