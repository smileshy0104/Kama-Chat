package main

import (
	"Kama-Chat/core"
	"Kama-Chat/global"
	"Kama-Chat/initialize/dao"
	"Kama-Chat/initialize/zlog"
	"Kama-Chat/lib/kafka"
	"Kama-Chat/lib/redis"
	"Kama-Chat/router"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. 配置管理初始化
	global.VIPER = core.Viper() // 加载配置文件（如 config.yaml）
	conf := global.CONFIG
	host := conf.MainConfig.Host
	port := conf.MainConfig.Port
	kafkaConfig := conf.KafkaConfig
	dao.InitMysql()

	if kafkaConfig.MessageMode == "kafka" {
		kafka.KafkaService.KafkaInit()
	}

	if kafkaConfig.MessageMode == "channel" {
		//go chat.ChatServer.Start()
	} else {
		//go chat.KafkaChatServer.Start()
	}

	go func() {
		// Win10本地部署
		//if err := router.Router.RunTLS(fmt.Sprintf("%s:%d", host, port), "pkg/ssl/127.0.0.1+2.pem", "pkg/ssl/127.0.0.1+2-key.pem"); err != nil {
		//	zlog.Fatal("server running fault")
		//	return
		//}
		// Ubuntu22.04云服务器部署
		//if err := https_server.GE.RunTLS(fmt.Sprintf("%s:%d", host, port), "/etc/ssl/certs/server.crt", "/etc/ssl/private/server.key"); err != nil {
		//	zlog.Fatal("server running fault")
		//	return
		//}
	}()

	// 设置信号监听
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	//<-quit

	address := fmt.Sprintf("%s:%d", host, port)
	//监听http服务
	srv := &http.Server{
		Addr:    address,
		Handler: router.Router,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("web服务器 启动失败 : %v\n", err)
		}
	}()

	log.Println("web 服务器启动成功： ", address)

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var done = make(chan struct{}, 1)
	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Printf("web server shutdown error: %v", err)
		} else {
			fmt.Println("web server shutdown ok")
		}
		done <- struct{}{}
	}()

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		fmt.Println("web server shutdown timeout")
	case <-done:
	}

	fmt.Println("program exit ok")

	if kafkaConfig.MessageMode == "kafka" {
		kafka.KafkaService.KafkaClose()
	}

	//chat.ChatServer.Close()

	zlog.Info("关闭服务器...")

	// 删除所有Redis键
	if err := redis.DeleteAllRedisKeys(); err != nil {
		zlog.Error(err.Error())
	} else {
		zlog.Info("所有Redis键已删除")
	}

	zlog.Info("服务器已关闭")
}
