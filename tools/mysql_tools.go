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
var dbOpenApi *sql.DB

// InitMySQLClient 初始化 MySQL 客户端
func InitMySQLClient() {
	var config, _ = GetAppConfig("config.yaml")
	var err error
	// 主库连接
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
		log.Fatal("main db ping error:", err)
	}

	// OpenAPI 日志库连接
	dsnOpenApi := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.OpenApiMySQL.Username, config.OpenApiMySQL.Password,
		config.OpenApiMySQL.Host, config.OpenApiMySQL.Port, config.OpenApiMySQL.Database)

	dbOpenApi, err = sql.Open("mysql", dsnOpenApi)
	if err != nil {
		log.Fatal("openapi db connect error:", err)
	}
	dbOpenApi.SetMaxIdleConns(5)
	dbOpenApi.SetMaxOpenConns(50)
	dbOpenApi.SetConnMaxLifetime(5 * time.Minute)

	if err = dbOpenApi.Ping(); err != nil {
		log.Fatal("openapi db ping error:", err)
	}
}

//////////////////////////////////////////////////////
// Novel (书籍) 操作
//////////////////////////////////////////////////////

// MySQLGetAllNovels (默认，只返回已审核通过的书籍)
func MySQLGetAllNovels() ([]models.Novel, bool) {
	query := `SELECT id, novel_id, name, author, category, status, description, cover_url, created_at, updated_at
              FROM novels
              WHERE audit_status = 1
              ORDER BY id DESC`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("查询所有小说失败: %v", err)
		return nil, false
	}
	defer rows.Close()

	var novels []models.Novel
	for rows.Next() {
		var n models.Novel
		err := rows.Scan(&n.ID, &n.Novel_Id, &n.Name, &n.Author, &n.Category, &n.Status, &n.Description, &n.CoverURL, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			log.Printf("扫描小说数据失败: %v", err)
			return nil, false
		}
		novels = append(novels, n)
	}
	return novels, true
}

// MySQLGetNovelByID (默认：只返回已通过审核的书籍)
func MySQLGetNovelByID(novel_id string) (models.Novel, bool) {
	var n models.Novel
	query := `SELECT id, novel_id, name, author, category, status, description, cover_url, created_at, updated_at
              FROM novels
              WHERE novel_id = ? AND audit_status = 1
              LIMIT 1`
	err := db.QueryRow(query, novel_id).Scan(&n.ID, &n.Novel_Id, &n.Name, &n.Author, &n.Category, &n.Status, &n.Description, &n.CoverURL, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return n, false
		}
		log.Printf("查询小说 ID=%v 出错: %v", novel_id, err)
		return n, false
	}
	return n, true
}

// MySQLCreateNovel 新增小说
func MySQLCreateNovel(n models.Novel) (int64, error) {
	query := `INSERT INTO novels (novel_id,name, author, category, status, description, cover_url, created_at, updated_at) VALUES (?,?, ?, ?, ?, ?, ?, NOW(), NOW())`
	result, err := db.Exec(query, n.Novel_Id, n.Name, n.Author, n.Category, n.Status, n.Description, n.CoverURL)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// MySQLUpdateNovel 更新小说
func MySQLUpdateNovel(n models.Novel) (bool, error) {
	query := `UPDATE novels SET name=?, author=?, category=?, status=?, description=?, cover_url=?, updated_at=NOW() WHERE novel_id =?`
	res, err := db.Exec(query, n.Name, n.Author, n.Category, n.Status, n.Description, n.CoverURL, n.Novel_Id)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

// MySQLDeleteNovel 删除小说
func MySQLDeleteNovel(novel_id string) (bool, error) {
	query := `DELETE FROM novels WHERE novel_id=?`
	res, err := db.Exec(query, novel_id)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

//////////////////////////////////////////////////////
// Chapter (章节) 操作
//////////////////////////////////////////////////////

// MySQLGetChaptersByNovelID (默认：只有当小说已审核通过时，返回该小说的所有章节)
func MySQLGetChaptersByNovelID(novel_id string) ([]models.Chapter, bool) {
	// 先检查书籍审核状态
	var auditStatus int
	err := db.QueryRow("SELECT audit_status FROM novels WHERE novel_id = ?", novel_id).Scan(&auditStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("小说不存在：%s", novel_id)
			return nil, false
		}
		log.Printf("查询小说审核状态失败: %v", err)
		return nil, false
	}
	if auditStatus != 1 {
		// 未审核通过，不允许查看章节
		return nil, false
	}

	query := `
        SELECT id, novel_id, title, content, word_count, chapter_index, created_at, updated_at
        FROM chapters 
        WHERE novel_id = ? 
        ORDER BY chapter_index ASC
    `
	rows, err := db.Query(query, novel_id)
	if err != nil {
		log.Printf("查询小说 ID=%v 的章节失败: %v", novel_id, err)
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

// MySQLGetChaptersDetailByNovelID (默认：只有当小说已审核通过时，返回章节目录)
func MySQLGetChaptersDetailByNovelID(novel_id string) ([]models.ChapterDetail, bool) {
	// 检查小说是否通过审核
	var auditStatus int
	err := db.QueryRow("SELECT audit_status FROM novels WHERE novel_id = ?", novel_id).Scan(&auditStatus)
	if err != nil {
		log.Printf("查询小说审核状态失败: %v", err)
		return nil, false
	}
	if auditStatus != 1 {
		return nil, false
	}

	query := `
        SELECT chapter_index, title
        FROM chapters 
        WHERE novel_id = ? 
        ORDER BY chapter_index ASC
    `
	rows, err := db.Query(query, novel_id)
	if err != nil {
		log.Printf("查询小说 ID=%v 的章节失败: %v", novel_id, err)
		return nil, false
	}
	defer rows.Close()

	var chaptersdetail []models.ChapterDetail
	for rows.Next() {
		var c models.ChapterDetail
		err := rows.Scan(&c.ChapterIndex, &c.Title)
		if err != nil {
			log.Printf("扫描章节数据失败: %v", err)
			return nil, false
		}
		chaptersdetail = append(chaptersdetail, c)
	}
	return chaptersdetail, true
}

// MySQLGetChaptersDetailByNovelIDForAudit (管理后台/审核专用：返回所有章节，不限制小说审核状态)
func MySQLGetChaptersDetailByNovelIDForAudit(novel_id string) ([]models.ChapterDetail, bool) {
	query := `
        SELECT chapter_index, title
        FROM chapters 
        WHERE novel_id = ? 
        ORDER BY chapter_index ASC
    `
	rows, err := db.Query(query, novel_id)
	if err != nil {
		log.Printf("查询小说 ID=%v 的章节失败: %v", novel_id, err)
		return nil, false
	}
	defer rows.Close()

	var chaptersdetail []models.ChapterDetail
	for rows.Next() {
		var c models.ChapterDetail
		err := rows.Scan(&c.ChapterIndex, &c.Title)
		if err != nil {
			log.Printf("扫描章节数据失败: %v", err)
			return nil, false
		}
		chaptersdetail = append(chaptersdetail, c)
	}
	return chaptersdetail, true
}

// MySQLGetChapterByID 根据ID获取章节
func MySQLGetChapterByID(novel_id string) (models.Chapter, bool) {
	var c models.Chapter
	query := `SELECT id, novel_id, title, content, word_count, chapter_index, created_at, updated_at FROM chapters WHERE novel_id=?`
	err := db.QueryRow(query, novel_id).Scan(&c.ID, &c.NovelID, &c.Title, &c.Content, &c.WordCount, &c.ChapterIndex, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, false
		}
		log.Printf("查询章节 ID=%v 出错: %v", novel_id, err)
		return c, false
	}
	return c, true
}

// 获取小说某一章（通过 novel_id + chapter_index）
func MySQLGetChapterByNovelIDAndIndex(novel_id string, chapterIndex int) (models.Chapter, bool) {
	println("执行到了这里")
	query := `
		SELECT c.id, c.novel_id, c.title, c.content, c.word_count, 
		       c.chapter_index, c.created_at, c.updated_at
		FROM chapters c
		INNER JOIN novels n ON c.novel_id = n.novel_id
		WHERE c.novel_id = ? 
		  AND c.chapter_index = ? 
		  AND n.audit_status = 1   -- 只允许已审核的小说
		LIMIT 1
	`
	row := db.QueryRow(query, novel_id, chapterIndex)

	var c models.Chapter
	err := row.Scan(
		&c.ID,
		&c.NovelID,
		&c.Title,
		&c.Content,
		&c.WordCount,
		&c.ChapterIndex,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		log.Printf("查询小说 ID=%v 第 %d 章失败: %v", novel_id, chapterIndex, err)
		return models.Chapter{}, false
	}
	return c, true
}

func MySQLGetLatestChapter(novel_id string, chapter *models.Chapter) error {
	query := `
		SELECT id, novel_id, title, content, word_count, chapter_index, created_at, updated_at
		FROM chapters
		WHERE novel_id = ?
		ORDER BY chapter_index DESC
		LIMIT 1
	`
	return db.QueryRow(query, novel_id).Scan(
		&chapter.ID,
		&chapter.NovelID,
		&chapter.Title,
		&chapter.Content,
		&chapter.WordCount,
		&chapter.ChapterIndex,
		&chapter.CreatedAt,
		&chapter.UpdatedAt,
	)
}

// 获取某本小说的最新章节序号
func MySQLGetLatestChapterIndex(novel_id string) (int, error) {
	var chapterIndex int
	query := `SELECT COALESCE(MAX(chapter_index), 0) FROM chapters WHERE novel_id = ?`
	err := db.QueryRow(query, novel_id).Scan(&chapterIndex)
	if err != nil {
		return 0, err
	}
	return chapterIndex, nil
}

// MySQLCreateChapterDetail 新增文章目录
func MySQLCreateChapterDetail(novel_id string, detail models.ChapterDetail) error {
	query := `INSERT INTO chapter_details (novel_id, title, chapter_index) VALUES (?, ?, ?)`
	_, err := db.Exec(query, novel_id, detail.Title, detail.ChapterIndex)
	return err
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
	query := `UPDATE chapters SET title=?, content=?, word_count=?, chapter_index=?, updated_at=NOW() WHERE novel_id=?`
	res, err := db.Exec(query, c.Title, c.Content, c.WordCount, c.ChapterIndex, c.NovelID)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

// MySQLDeleteChapter 删除章节
func MySQLDeleteChapter(novel_id string) (bool, error) {
	query := `DELETE FROM chapters WHERE novel_id=?`
	res, err := db.Exec(query, novel_id)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

// 更新章节内容
func MySQLUpdateChapterByNovelIDAndIndex(novel_id string, chapterIndex int, updated models.Chapter) (bool, error) {
	query := `
		UPDATE chapters
		SET title = ?, content = ?, word_count = ?, updated_at = ?
		WHERE novel_id = ? AND chapter_index = ?
	`
	result, err := db.Exec(query,
		updated.Title,
		updated.Content,
		updated.WordCount,
		time.Now(),
		novel_id,
		chapterIndex,
	)
	if err != nil {
		log.Printf("更新小说 ID=%v 第 %d 章失败: %v", novel_id, chapterIndex, err)
		return false, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}

// 删除章节
func MySQLDeleteChapterByNovelIDAndIndex(novel_id string, chapterIndex int) (bool, error) {
	query := `
		DELETE FROM chapters
		WHERE novel_id = ? AND chapter_index = ?
	`
	result, err := db.Exec(query, novel_id, chapterIndex)
	if err != nil {
		log.Printf("删除小说 ID=%v 第 %d 章失败: %v", novel_id, chapterIndex, err)
		return false, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}

// MySQLGetOpenapiUserbyApplesub 根据ID获取用户信息
func MySQLGetOpenapiUserbyApplesub(apple_sub string) (models.User, bool) {
	var user models.User

	query := `SELECT id, username, sub, apple_sub 
	          FROM openapi_user 
	          WHERE apple_sub = ? 
	          LIMIT 1`

	err := dbOpenApi.QueryRow(query, apple_sub).Scan(
		&user.ID,
		&user.Username,
		&user.Sub,
		&user.AppleSub,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, false
		}
		log.Println("query error:", err)
		return user, false
	}

	return user, true
}

// MySQLCreateOpenapiUser 新增用户
func MySQLCreateOpenapiUser(user models.User) (int64, error) {
	query := `INSERT INTO openapi_user (sub,username, apple_sub, email,is_active) VALUES (?,?,?,?,0)`
	result, err := dbOpenApi.Exec(query,
		user.Sub,
		user.Nickname,
		user.AppleSub,
		user.Email,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// MySQLAuditNovel (更新小说审核状态：0/1/2)
func MySQLAuditNovel(novelID string, status int, auditBy string) (bool, error) {
	query := `UPDATE novels SET audit_status = ?, audit_by = ?, audit_at = NOW(), updated_at = NOW() WHERE novel_id = ?`
	res, err := db.Exec(query, status, auditBy, novelID)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

// MySQLGetNovelWithChaptersForAudit (审核专用：获取一本书并带上其章节)
func MySQLGetNovelWithChaptersForAudit(novel_id string) (models.Novel, []models.Chapter, bool) {
	novel, ok := MySQLGetNovelByIDForAudit(novel_id)
	if !ok {
		return models.Novel{}, nil, false
	}
	chapters, ok := MySQLGetChaptersByNovelIDForAudit(novel_id)
	if !ok {
		// 即便书存在，可能没有章节也返回空 slice
		return novel, nil, true
	}
	return novel, chapters, true
}

// MySQLGetChaptersByNovelIDForAudit (审核专用：返回某书的所有章节，不受书籍审核状态限制)
func MySQLGetChaptersByNovelIDForAudit(novel_id string) ([]models.Chapter, bool) {
	query := `
        SELECT id, novel_id, title, content, word_count, chapter_index, created_at, updated_at
        FROM chapters 
        WHERE novel_id = ? 
        ORDER BY chapter_index ASC
    `
	rows, err := db.Query(query, novel_id)
	if err != nil {
		log.Printf("查询小说 ID=%v 的章节失败: %v", novel_id, err)
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

// MySQLGetNovelByIDForAudit (审核专用：返回书籍并包含审核字段)
func MySQLGetNovelByIDForAudit(novel_id string) (models.Novel, bool) {
	var n models.Novel
	query := `SELECT id, novel_id, name, author, category, status, description, cover_url, audit_status, audit_by, audit_at, created_at, updated_at
              FROM novels
              WHERE novel_id = ?
              LIMIT 1`
	err := db.QueryRow(query, novel_id).Scan(&n.ID, &n.Novel_Id, &n.Name, &n.Author, &n.Category, &n.Status, &n.Description, &n.CoverURL, &n.AuditStatus, &n.AuditBy, &n.AuditAt, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return n, false
		}
		log.Printf("查询小说(审核用) ID=%v 出错: %v", novel_id, err)
		return n, false
	}
	return n, true
}

// MySQLGetAllNovelsForAudit (管理后台/审核专用：返回所有书籍，带审核字段)
func MySQLGetAllNovelsForAudit() ([]models.Novel, bool) {
	query := `SELECT id, novel_id, name, author, category, status, description, cover_url, audit_status, audit_by, audit_at, created_at, updated_at
              FROM novels
              ORDER BY id DESC`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("查询所有小说(审核专用)失败: %v", err)
		return nil, false
	}
	defer rows.Close()

	var novels []models.Novel
	for rows.Next() {
		var n models.Novel
		err := rows.Scan(&n.ID, &n.Novel_Id, &n.Name, &n.Author, &n.Category, &n.Status, &n.Description, &n.CoverURL, &n.AuditStatus, &n.AuditBy, &n.AuditAt, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			log.Printf("扫描小说数据失败: %v", err)
			return nil, false
		}
		novels = append(novels, n)
	}
	return novels, true
}

// MySQLUpdateNovelAudit 更新小说的审核状态
func MySQLUpdateNovelAudit(novelId string, status int, auditBy string, auditAt time.Time) (bool, error) {
	query := `UPDATE novels 
	          SET audit_status = ?, audit_by = ?, audit_at = ?, updated_at = NOW()
	          WHERE novel_id = ?`

	result, err := db.Exec(query, status, auditBy, auditAt, novelId)
	if err != nil {
		log.Printf("更新小说审核状态失败 novel_id=%v err=%v", novelId, err)
		return false, err
	}

	rows, _ := result.RowsAffected()
	return rows > 0, nil
}
