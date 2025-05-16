package main

// 管理剪切板历史数据
type Data struct {
	dbFile string
	recent []CopyItem
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

	//检查sqlite db file是否存在，不存说明是第一次启动程序，需要初始化table，
	//table的字段： id, content(内容)，content_type(内容类型,int 型),content_hash,is_collect(是否收藏)，use_count(使用次数), created_at(时间戳格式), updated_at(时间戳格式)
	//content_hash 需要加索引， collect需要加索引
	//todo

	return &Data{
		dbFile: path,
		recent: nil,
	}
}

// 新增条目
func (d *Data) NewItem(item CopyItem) {
	d.recent = append([]CopyItem{item}, d.recent...)

	if len(d.recent) > recentSize {
		d.recent = d.recent[:recentSize]
	}

	//写入table,注意通过 content_hash 检查重复
	//todo
}

// 初次启动程序冲数据库中读取最近的100条数据到recent中
func (d *Data) GetRecent() {
	//todo
}

// 取消/收藏 指定内容
func (d *Data) SetCollect(id int, isCollect bool) {
	//todo
}

// 增加一次使用次数
func (d *Data) AddCount(id int) {
	//todo
}

// 根据关键词搜索,关键词支持多个，多个关键字用 and 关系
// 按照时间倒序
func (d *Data) Search(keywords []string, limit int) {
	//todo
}
