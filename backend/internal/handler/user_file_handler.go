package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
	"xorm.io/xorm"

	"billsoftware/backend/internal/config"
	"billsoftware/backend/internal/model"
	"billsoftware/backend/internal/response"
	"billsoftware/backend/internal/storage"
)

const ledgerImageMaxBytesDefault int64 = 5 * 1024 * 1024

type UserFileHandler struct {
	engine        *xorm.Engine
	storage       *storage.ObjectStorage
	thumbnailSize int
	maxUploadSize int64
}

func NewUserFileHandler(engine *xorm.Engine, objectStorage *storage.ObjectStorage, cfg config.MinIOConfig) *UserFileHandler {
	thumbnailSize := cfg.ThumbnailSize
	if thumbnailSize <= 0 {
		thumbnailSize = 160
	}

	maxUploadSize := cfg.MaxUploadMB * 1024 * 1024
	if maxUploadSize <= 0 {
		maxUploadSize = ledgerImageMaxBytesDefault
	}

	return &UserFileHandler{
		engine:        engine,
		storage:       objectStorage,
		thumbnailSize: thumbnailSize,
		maxUploadSize: maxUploadSize,
	}
}

func (h *UserFileHandler) Upload(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "image file is required")
		return
	}
	if fileHeader.Size > h.maxUploadSize {
		response.Fail(c, http.StatusBadRequest, "image size must not exceed 5MB")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "open image file failed")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, h.maxUploadSize+1))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "read image file failed")
		return
	}
	if int64(len(data)) > h.maxUploadSize {
		response.Fail(c, http.StatusBadRequest, "image size must not exceed 5MB")
		return
	}

	mimeType, extension, ok := detectAllowedImage(data)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "only jpeg, png and webp images are allowed")
		return
	}

	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "decode image failed")
		return
	}
	bounds := decoded.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	thumbnail, err := buildJPEGThumbnail(decoded, h.thumbnailSize)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "build thumbnail failed")
		return
	}

	now := time.Now()
	randomID, err := randomHex(16)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "generate object key failed")
		return
	}
	baseKey := fmt.Sprintf("ledger/%d/%s/%s", userID, now.Format("2006/01"), randomID)
	originalKey := baseKey + "/original" + extension
	thumbnailKey := baseKey + "/thumbnail.jpg"

	ctx := c.Request.Context()
	if err := h.storage.Put(ctx, originalKey, data, mimeType); err != nil {
		response.Fail(c, http.StatusInternalServerError, "upload original image failed")
		return
	}
	if err := h.storage.Put(ctx, thumbnailKey, thumbnail, "image/jpeg"); err != nil {
		_ = h.storage.Delete(ctx, originalKey)
		response.Fail(c, http.StatusInternalServerError, "upload thumbnail image failed")
		return
	}

	upload := &model.FileUpload{
		UserID:             userID,
		BizType:            model.FileBizTypeLedgerImage,
		OriginalName:       filepath.Base(fileHeader.Filename),
		MimeType:           mimeType,
		SizeBytes:          uint64(len(data)),
		OriginalObjectKey:  originalKey,
		ThumbnailObjectKey: thumbnailKey,
		Width:              width,
		Height:             height,
		Status:             model.FileStatusUploaded,
	}
	if _, err := h.engine.Insert(upload); err != nil {
		_ = h.storage.Delete(ctx, originalKey, thumbnailKey)
		response.Fail(c, http.StatusInternalServerError, "save image metadata failed")
		return
	}

	originalURL, err := h.storage.PresignedGetURL(ctx, originalKey)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "build original image url failed")
		return
	}
	thumbnailURL, err := h.storage.PresignedGetURL(ctx, thumbnailKey)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "build thumbnail image url failed")
		return
	}

	response.Success(c, gin.H{
		"fileId":       upload.ID,
		"originalName": upload.OriginalName,
		"mimeType":     upload.MimeType,
		"sizeBytes":    upload.SizeBytes,
		"width":        upload.Width,
		"height":       upload.Height,
		"originalUrl":  originalURL,
		"thumbnailUrl": thumbnailURL,
	})
}

func detectAllowedImage(data []byte) (string, string, bool) {
	if len(data) >= 3 && data[0] == 0xff && data[1] == 0xd8 && data[2] == 0xff {
		return "image/jpeg", ".jpg", true
	}
	if len(data) >= 8 && bytes.Equal(data[:8], []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}) {
		return "image/png", ".png", true
	}
	if len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return "image/webp", ".webp", true
	}
	return "", "", false
}

func buildJPEGThumbnail(src image.Image, maxSize int) ([]byte, error) {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid image bounds")
	}

	scale := float64(maxSize) / float64(width)
	if height > width {
		scale = float64(maxSize) / float64(height)
	}
	if scale > 1 {
		scale = 1
	}

	targetWidth := int(float64(width) * scale)
	targetHeight := int(float64(height) * scale)
	if targetWidth < 1 {
		targetWidth = 1
	}
	if targetHeight < 1 {
		targetHeight = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 82}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func randomHex(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return strings.ToLower(hex.EncodeToString(buf)), nil
}
