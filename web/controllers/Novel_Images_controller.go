package controllers

import (
	"github.com/kataras/iris/v12"
	"novel-server/tools"
	"os"
	"path/filepath"
)

type ImageController struct {
	Ctx iris.Context
}

func (c *ImageController) Get(format, year, month, day, filename string) {
	uploadDir := "uploads"
	safeFilename := filepath.Base(filename)
	filePath := filepath.Join(uploadDir, format, year, month, day, safeFilename)

	// 路径安全检查
	absFilePath, err := tools.EnsureSafePath(uploadDir, filePath)
	if err != nil {
		c.Ctx.Application().Logger().Warnf("Blocked path traversal attempt: %s", filePath)
		c.Ctx.StatusCode(iris.StatusForbidden)
		_ = c.Ctx.JSON(iris.Map{"error": "access denied"})
		return
	}

	// 文件是否存在
	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		c.Ctx.StatusCode(iris.StatusNotFound)
		_ = c.Ctx.JSON(iris.Map{"error": "file not found"})
		return
	}

	c.Ctx.Application().Logger().Infof("Serving file: %s", absFilePath)
	c.Ctx.ServeFile(absFilePath)
}
