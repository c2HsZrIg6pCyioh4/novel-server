package tools

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"novel-server/web/models"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitMySQLClient 初始化 MySQL 客户端
func InitMySQLClient() {
	var config, _ = GetAppConfig("config.yaml")
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MySQL.Username, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port, config.MySQL.Database)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 配置连接池参数
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

//////////////////////////////////////////////////////
// Novel (书籍) 操作
//////////////////////////////////////////////////////

// MySQLGetAllNovels 获取所有小说
func MySQLGetAllNovels() ([]models.Novel, bool) {
	query := `SELECT id, name, author, category, status, description, cover_url, created_at, updated_at FROM novels ORDER BY id DESC`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("查询所有小说失败: %v", err)
		return nil, false
	}
	defer rows.Close()

	var novels []models.Novel
	for rows.Next() {
		var n models.Novel
		err := rows.Scan(&n.ID, &n.Name, &n.Author, &n.Category, &n.Status, &n.Description, &n.CoverURL, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			log.Printf("扫描小说数据失败: %v", err)
			return nil, false
		}
		novels = append(novels, n)
	}
	return novels, true
}

// MySQLGetNovelByID 根据ID获取小说
func MySQLGetNovelByID(id int64) (models.Novel, bool) {
	var n models.Novel
	query := `SELECT id, name, author, category, status, description, cover_url, created_at, updated_at FROM novels WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&n.ID, &n.Name, &n.Author, &n.Category, &n.Status, &n.Description, &n.CoverURL, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return n, false
		}
		log.Printf("查询小说 ID=%d 出错: %v", id, err)
		return n, false
	}
	return n, true
}

// MySQLCreateNovel 新增小说
func MySQLCreateNovel(n models.Novel) (int64, error) {
	query := `INSERT INTO novels (name, author, category, status, description, cover_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`
	result, err := db.Exec(query, n.Name, n.Author, n.Category, n.Status, n.Description, n.CoverURL)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// MySQLUpdateNovel 更新小说
func MySQLUpdateNovel(n models.Novel) (bool, error) {
	query := `UPDATE novels SET name=?, author=?, category=?, status=?, description=?, cover_url=?, updated_at=NOW() WHERE id=?`
	res, err := db.Exec(query, n.Name, n.Author, n.Category, n.Status, n.Description, n.CoverURL, n.ID)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

// MySQLDeleteNovel 删除小说
func MySQLDeleteNovel(id int64) (bool, error) {
	query := `DELETE FROM novels WHERE id=?`
	res, err := db.Exec(query, id)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

//////////////////////////////////////////////////////
// Chapter (章节) 操作
//////////////////////////////////////////////////////

// MySQLGetChaptersByNovelID 获取某本小说的所有章节
func MySQLGetChaptersByNovelID(novelID int64) ([]models.Chapter, bool) {
	query := `
		SELECT id, novel_id, title, content, word_count, chapter_index, created_at, updated_at
		FROM chapters 
		WHERE novel_id = ? 
		ORDER BY chapter_index ASC
	`

	rows, err := db.Query(query, novelID)
	if err != nil {
		log.Printf("查询小说 ID=%d 的章节失败: %v", novelID, err)
		return nil, false
	}
	defer rows.Close()

	var chapters []models.Chapter
	for rows.Next() {
		var c models.Chapter
		err := rows.Scan(&c.ID, &c.NovelID, &c.Title, &c.Content, &c.WordCount, &c.ChapterIndex, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			log.Printf("扫描章节数据失败: %v", err)
			return nil, false
		}
		chapters = append(chapters, c)
	}
	return chapters, true
}

// MySQLGetChapterByID 根据ID获取章节
func MySQLGetChapterByID(id int64) (models.Chapter, bool) {
	var c models.Chapter
	query := `SELECT id, novel_id, title, content, word_count, chapter_index, created_at, updated_at FROM chapters WHERE id=?`
	err := db.QueryRow(query, id).Scan(&c.ID, &c.NovelID, &c.Title, &c.Content, &c.WordCount, &c.ChapterIndex, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, false
		}
		log.Printf("查询章节 ID=%d 出错: %v", id, err)
		return c, false
	}
	return c, true
}

// MySQLCreateChapter 新增章节
func MySQLCreateChapter(c models.Chapter) (int64, error) {
	if c.WordCount == 0 && c.Content != "" {
		c.WordCount = len([]rune(c.Content))
	}
	query := `INSERT INTO chapters (novel_id, title, content, word_count, chapter_index, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())`
	result, err := db.Exec(query, c.NovelID, c.Title, c.Content, c.WordCount, c.ChapterIndex)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// MySQLUpdateChapter 更新章节
func MySQLUpdateChapter(c models.Chapter) (bool, error) {
	if c.WordCount == 0 && c.Content != "" {
		c.WordCount = len([]rune(c.Content))
	}
	query := `UPDATE chapters SET title=?, content=?, word_count=?, chapter_index=?, updated_at=NOW() WHERE id=?`
	res, err := db.Exec(query, c.Title, c.Content, c.WordCount, c.ChapterIndex, c.ID)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

// MySQLDeleteChapter 删除章节
func MySQLDeleteChapter(id int64) (bool, error) {
	query := `DELETE FROM chapters WHERE id=?`
	res, err := db.Exec(query, id)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}
