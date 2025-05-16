package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	dbFile string
	recent []CopyItem
	db     *sql.DB
}

type CopyItem struct {
	id          int
	content     string
	contentType int
	contentHash string
	isCollect   bool
	useCount    int
	createdAt   int64
	updatedAt   int64
}

const recentSize = 10
const tableName = "clipboard"

func NewData(path string) *Data {
	// 检查数据库文件是否存在
	firstRun := false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		firstRun = true
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}

	if firstRun {
		sqlStmt := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT,
			content_type INTEGER,
			content_hash TEXT,
			is_collect BOOLEAN,
			use_count INTEGER,
			created_at INTEGER,
			updated_at INTEGER
		);
		CREATE INDEX IF NOT EXISTS idx_hash ON %s(content_hash);
		CREATE INDEX IF NOT EXISTS idx_collect ON %s(is_collect);`, tableName, tableName, tableName)

		if _, err := db.Exec(sqlStmt); err != nil {
			log.Fatalf("初始化数据库失败: %v", err)
		}
	}

	return &Data{
		dbFile: path,
		recent: nil,
		db:     db,
	}
}

func (d *Data) NewItem(item CopyItem) {
	// 去重
	var existingID int
	err := d.db.QueryRow(fmt.Sprintf("SELECT id FROM %s WHERE content_hash = ?", tableName), item.contentHash).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("查询失败:", err)
		return
	}
	if existingID != 0 {
		log.Println("内容已存在，跳过")
		return
	}

	now := time.Now().Unix()
	item.createdAt = now
	item.updatedAt = now

	_, err = d.db.Exec(fmt.Sprintf(`INSERT INTO %s (content, content_type, content_hash, is_collect, use_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, tableName),
		item.content, item.contentType, item.contentHash, item.isCollect, item.useCount, item.createdAt, item.updatedAt)
	if err != nil {
		log.Println("插入失败:", err)
	}

	d.recent = append([]CopyItem{item}, d.recent...)
	if len(d.recent) > recentSize {
		d.recent = d.recent[:recentSize]
	}
}

func (d *Data) GetRecent() {
	rows, err := d.db.Query(fmt.Sprintf(`SELECT id, content, content_type, content_hash, is_collect, use_count, created_at, updated_at 
		FROM %s ORDER BY created_at DESC LIMIT 100`, tableName))
	if err != nil {
		log.Println("查询失败:", err)
		return
	}
	defer rows.Close()

	var recent []CopyItem
	for rows.Next() {
		var item CopyItem
		err := rows.Scan(&item.id, &item.content, &item.contentType, &item.contentHash, &item.isCollect, &item.useCount, &item.createdAt, &item.updatedAt)
		if err != nil {
			log.Println("解析失败:", err)
			continue
		}
		recent = append(recent, item)
	}
	d.recent = recent
}

func (d *Data) SetCollect(id int, isCollect bool) {
	_, err := d.db.Exec(fmt.Sprintf(`UPDATE %s SET is_collect = ?, updated_at = ? WHERE id = ?`, tableName), isCollect, time.Now().Unix(), id)
	if err != nil {
		log.Println("更新收藏状态失败:", err)
	}
}

func (d *Data) AddCount(id int) {
	_, err := d.db.Exec(fmt.Sprintf(`UPDATE %s SET use_count = use_count + 1, updated_at = ? WHERE id = ?`, tableName), time.Now().Unix(), id)
	if err != nil {
		log.Println("增加使用次数失败:", err)
	}
}

func (d *Data) Search(keywords []string, limit int) {
	if len(keywords) == 0 {
		fmt.Println("无关键词")
		return
	}

	var conditions []string
	var args []interface{}
	for _, kw := range keywords {
		conditions = append(conditions, "content LIKE ?")
		args = append(args, "%"+kw+"%")
	}

	query := fmt.Sprintf(`SELECT id, content, content_type, content_hash, is_collect, use_count, created_at, updated_at 
		FROM %s WHERE %s ORDER BY created_at DESC LIMIT ?`, tableName, strings.Join(conditions, " AND "))
	args = append(args, limit)

	rows, err := d.db.Query(query, args...)
	if err != nil {
		log.Println("搜索失败:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item CopyItem
		err := rows.Scan(&item.id, &item.content, &item.contentType, &item.contentHash, &item.isCollect, &item.useCount, &item.createdAt, &item.updatedAt)
		if err != nil {
			log.Println("读取结果失败:", err)
			continue
		}
		fmt.Printf("ID:%d 内容:%s 使用次数:%d 收藏:%v\n", item.id, item.content, item.useCount, item.isCollect)
	}
}
