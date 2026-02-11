package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// 从环境变量或配置文件读取数据库配置
	host := getEnvOrDefault("DB_HOST", "192.168.18.6")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "postgres")
	password := getEnvOrDefault("DB_PASSWORD", "123456")
	dbname := getEnvOrDefault("DB_NAME", "nookverse")

	fmt.Printf("开始数据库迁移: %s@%s:%s/%s\n", user, host, port, dbname)

	// 第一步：创建数据库（如果不存在）
	if err := createDatabaseIfNotExists(host, port, user, password, dbname); err != nil {
		log.Fatal("创建数据库失败:", err)
	}

	// 第二步：连接到目标数据库并执行初始化脚本
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("无法连接数据库:", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	fmt.Println("数据库连接成功!")

	// 读取SQL文件
	sqlFile := "db/init.sql"
	content, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Fatal("读取SQL文件失败:", err)
	}

	// 执行SQL
	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatal("执行SQL失败:", err)
	}

	fmt.Println("数据库迁移完成!")
	
	// 验证结果
	validateDatabase(db)
}

func createDatabaseIfNotExists(host, port, user, password, dbname string) error {
	// 先连接到默认数据库（postgres）
	defaultConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password)

	fmt.Printf("连接到PostgreSQL服务器: %s@%s:%s\n", user, host, port)

	db, err := sql.Open("postgres", defaultConnStr)
	if err != nil {
		return fmt.Errorf("无法连接到PostgreSQL服务器: %w", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("PostgreSQL服务器连接失败: %w", err)
	}

	fmt.Println("PostgreSQL服务器连接成功!")

	// 检查数据库是否已存在
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbname).Scan(&exists)
	if err != nil {
		return fmt.Errorf("检查数据库存在性失败: %w", err)
	}

	if exists {
		fmt.Printf("数据库 '%s' 已存在\n", dbname)
	} else {
		// 创建数据库
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
		if err != nil {
			return fmt.Errorf("创建数据库失败: %w", err)
		}
		fmt.Printf("数据库 '%s' 创建成功!\n", dbname)
	}

	return nil
}

func validateDatabase(db *sql.DB) {
	fmt.Println("\n=== 数据库验证 ===")
	
	// 检查关键表
	tables := []string{"houses", "rooms", "categories", "items", "users"}
	for _, table := range tables {
		var count int
		err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			fmt.Printf("%s 表: 查询失败 (%v)\n", table, err)
		} else {
			fmt.Printf("%s 表记录数: %d\n", table, count)
		}
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}