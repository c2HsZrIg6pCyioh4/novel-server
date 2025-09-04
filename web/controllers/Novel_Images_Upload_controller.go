package controllers

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"novel-server/tools"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/chai2010/webp"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type UploadController struct {
	Ctx iris.Context
}

type UploadResp struct {
	OK      bool              `json:"ok"`
	Message string            `json:"message,omitempty"`
	Links   map[string]string `json:"links,omitempty"`
}

func (c *UploadController) Post() mvc.Result {
	// 获取上传文件
	file, fileHeader, err := c.Ctx.FormFile("file")
	if err != nil {
		return mvc.Response{Code: iris.StatusBadRequest, Object: UploadResp{OK: false, Message: "file required"}}
	}
	defer file.Close()

	// 检查文件大小 (5MB限制)
	if fileHeader.Size > 5<<20 {
		return mvc.Response{Code: iris.StatusBadRequest, Object: UploadResp{OK: false, Message: "file too large (max 5MB)"}}
	}

	// 使用有限缓冲读取
	var buf bytes.Buffer
	tee := io.TeeReader(file, &buf)
	maxBytes := int64(5 << 20) // 5MB
	limitedReader := &io.LimitedReader{R: tee, N: maxBytes}

	// 读取前512字节用于MIME检测
	mimeBuf := make([]byte, 512)
	if _, err := limitedReader.Read(mimeBuf); err != nil {
		return mvc.Response{Code: iris.StatusInternalServerError, Object: UploadResp{OK: false, Message: "read failed"}}
	}

	// Magic number 检查
	signatures := map[string][]byte{
		"image/jpeg": {0xFF, 0xD8, 0xFF},
		"image/png":  {0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
		"image/gif":  {0x47, 0x49, 0x46, 0x38},
		"image/webp": {0x52, 0x49, 0x46, 0x46},
	}

	validMagic := false
	for _, sig := range signatures {
		if len(mimeBuf) >= len(sig) {
			match := true
			for i, b := range sig {
				if mimeBuf[i] != b {
					match = false
					break
				}
			}
			if match {
				validMagic = true
				break
			}
		}
	}

	if !validMagic {
		return mvc.Response{Code: iris.StatusBadRequest, Object: UploadResp{OK: false, Message: "invalid image signature"}}
	}

	// MIME 检查
	mime := http.DetectContentType(mimeBuf)
	allowedMimes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	// Block potentially dangerous formats
	blockedMimes := map[string]bool{
		"image/svg+xml":            true,
		"image/x-icon":             true,
		"image/tiff":               true,
		"image/vnd.microsoft.icon": true,
	}
	if blockedMimes[mime] {
		return mvc.Response{Code: iris.StatusBadRequest, Object: UploadResp{OK: false, Message: "format not allowed"}}
	}
	if !allowedMimes[mime] {
		return mvc.Response{Code: iris.StatusBadRequest, Object: UploadResp{OK: false, Message: "invalid mime type"}}
	}

	// 解码验证，防伪造
	img, _, err := image.Decode(io.MultiReader(bytes.NewReader(mimeBuf), limitedReader))
	if err != nil {
		return mvc.Response{Code: iris.StatusBadRequest, Object: UploadResp{OK: false, Message: "invalid image format"}}
	}

	// 去除元数据
	cleanImg := image.NewRGBA(img.Bounds())
	draw.Draw(cleanImg, img.Bounds(), img, image.Point{}, draw.Src)
	img = cleanImg

	// 检查图片尺寸
	const maxDimension = 5000
	bounds := img.Bounds()
	if bounds.Dx() > maxDimension || bounds.Dy() > maxDimension {
		return mvc.Response{Code: iris.StatusBadRequest, Object: UploadResp{OK: false, Message: "image dimensions too large"}}
	}

	// 计算文件MD5
	hash := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
	outDir := "./uploads"
	if err := os.MkdirAll(outDir, 0750); err != nil {
		c.Ctx.Application().Logger().Errorf("Failed to create upload directory: %v", err)
		return mvc.Response{Code: iris.StatusInternalServerError, Object: UploadResp{OK: false, Message: "server configuration error"}}
	}

	// 获取当前日期
	now := time.Now()
	year := fmt.Sprintf("%d", now.Year())
	month := fmt.Sprintf("%02d", now.Month())
	day := fmt.Sprintf("%02d", now.Day())

	// 文件名使用UUID
	uuid := uuid.New().String()
	c.Ctx.Application().Logger().Infof("Uploaded file MD5: %s, UUID: %s", hash, uuid)

	links := make(map[string]string)

	// 并发保存不同格式
	var wg sync.WaitGroup
	var mu sync.Mutex
	formats := []struct {
		ext     string
		encoder func(io.Writer, image.Image) error
	}{
		{".png", func(w io.Writer, img image.Image) error { return png.Encode(w, img) }},
		{".webp", func(w io.Writer, img image.Image) error {
			return webp.Encode(w, img, &webp.Options{Lossless: false, Quality: 80})
		}},
		{".jpg", func(w io.Writer, img image.Image) error { return jpeg.Encode(w, img, &jpeg.Options{Quality: 85}) }},
	}

	for _, format := range formats {
		wg.Add(1)
		go func(f struct {
			ext     string
			encoder func(io.Writer, image.Image) error
		}) {
			defer wg.Done()

			formatDir := filepath.Join(outDir, f.ext[1:], year, month, day)

			// 检查目录安全性
			safeDir, err := tools.EnsureSafePath(outDir, formatDir)
			if err != nil {
				c.Ctx.Application().Logger().Errorf("Security violation: %v", err)
				return
			}
			if err := os.MkdirAll(safeDir, 0750); err != nil {
				c.Ctx.Application().Logger().Errorf("Failed to create directory %s: %v", safeDir, err)
				return
			}

			// 生成文件路径
			uuidFilename := uuid + f.ext
			path := filepath.Join(safeDir, uuidFilename)

			// 再次校验最终路径
			safePath, err := tools.EnsureSafePath(outDir, path)
			if err != nil {
				c.Ctx.Application().Logger().Errorf("Security violation: %v", err)
				return
			}

			file, err := os.Create(safePath)
			if err != nil {
				c.Ctx.Application().Logger().Errorf("Failed to create %s: %v", safePath, err)
				return
			}
			defer file.Close()

			if err := f.encoder(file, img); err != nil {
				c.Ctx.Application().Logger().Errorf("Failed to encode %s: %v", safePath, err)
				return
			}

			// 生成访问路径
			fullPath := filepath.ToSlash(filepath.Join(f.ext[1:], year, month, day, uuidFilename))

			mu.Lock()
			links[f.ext[1:]] = "/" + fullPath
			mu.Unlock()
		}(format)
	}
	wg.Wait()

	if len(links) < len(formats) {
		c.Ctx.Application().Logger().Warnf("Not all formats were generated successfully. Expected %d, got %d", len(formats), len(links))
	}

	return mvc.Response{
		Object: UploadResp{OK: true, Links: links},
	}
}
