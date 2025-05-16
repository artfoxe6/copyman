package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// 管理剪切板历史数据
type Data struct {
	db     *sql.DB  // 数据库连接
	dbFile string   // 数据库文件路径
	recent []string // 最近使用的条目
}

const recentSize = 10
const tableName = "clipboard"

func NewData(path string) *Data {
	// 初始化数据库连接
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}

	// 创建表结构
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS clipboard (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT,
		content TEXT UNIQUE,
		collect BOOLEAN DEFAULT 0,
		use_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err = db.Exec(createTableSQL); err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	// 创建内容索引
	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_content ON clipboard(content);`
	if _, err = db.Exec(createIndexSQL); err != nil {
		log.Printf("创建索引失败: %v", err)
	}

	// 创建更新触发器
	createTriggerSQL := `
	CREATE TRIGGER IF NOT EXISTS update_clipboard_updated_at
	AFTER UPDATE ON clipboard
	FOR EACH ROW
	BEGIN
		UPDATE clipboard SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END;`
	if _, err = db.Exec(createTriggerSQL); err != nil {
		log.Printf("创建触发器失败: %v", err)
	}

	data := &Data{
		db:     db,
		dbFile: path,
		recent: make([]string, 0),
	}

	// 初始化最近条目
	data.GetRecent()
	return data
}

// 新增条目
func (d *Data) NewItem(item string) {
	// 更新内存中的最近列表
	d.recent = append([]string{item}, d.recent...)
	if len(d.recent) > recentSize {
		d.recent = d.recent[:recentSize]
	}

	// 写入数据库（带重复检查）
	upsertSQL := `
	INSERT INTO clipboard (content, type, use_count) 
	VALUES (?, ?, 1)
	ON CONFLICT(content) DO UPDATE SET 
		use_count = use_count + 1,
		updated_at = CURRENT_TIMESTAMP;`

	if _, err := d.db.Exec(upsertSQL, item, "text"); err != nil {
		log.Printf("插入数据失败: %v", err)
	}
}

// 从数据库加载最近条目
func (d *Data) GetRecent() {
	querySQL := `
	SELECT content 
	FROM clipboard 
	ORDER BY created_at DESC 
	LIMIT ?`

	rows, err := d.db.Query(querySQL, recentSize)
	if err != nil {
		log.Printf("查询最近条目失败: %v", err)
		return
	}
	defer rows.Close()

	d.recent = make([]string, 0)
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			log.Printf("扫描数据失败: %v", err)
			continue
		}
		d.recent = append(d.recent, content)
	}
}

// 收藏指定内容
func (d *Data) Collect(id int) {
	if _, err := d.db.Exec("UPDATE clipboard SET collect = 1 WHERE id = ?", id); err != nil {
		log.Printf("收藏失败: %v", err)
	}
}

// 取消收藏指定内容
func (d *Data) UnCollect(id int) {
	if _, err := d.db.Exec("UPDATE clipboard SET collect = 0 WHERE id = ?", id); err != nil {
		log.Printf("取消收藏失败: %v", err)
	}
}

// 增加使用次数
func (d *Data) IncreaseUseCount(id int) {
	if _, err := d.db.Exec("UPDATE clipboard SET use_count = use_count + 1 WHERE id = ?", id); err != nil {
		log.Printf("更新使用次数失败: %v", err)
	}
}

// 搜索内容（返回结果切片）
func (d *Data) Search(keyword string, size int) []string {
	querySQL := `
	SELECT content 
	FROM clipboard 
	WHERE content LIKE ? 
	ORDER BY use_count DESC 
	LIMIT ?`

	rows, err := d.db.Query(querySQL, "%"+keyword+"%", size)
	if err != nil {
		log.Printf("搜索失败: %v", err)
		return nil
	}
	defer rows.Close()

	results := make([]string, 0)
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			log.Printf("扫描结果失败: %v", err)
			continue
		}
		results = append(results, content)
	}
	return results
}
