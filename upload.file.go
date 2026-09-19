package blueprints

var UploadStorage = `package upload

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	OriginalName string ` + "`json:\"original_name\"`" + `
	StoredName string ` + "`json:\"stored_name\"`" + `
	Size int64 ` + "`json:\"size\"`" + `
	ContentType string ` + "`json:\"content_type\"`" + `
}

type Storage interface { Save(ctx context.Context, file *multipart.FileHeader) (File,error) }

type LocalStorage struct{ directory string }

func NewLocalStorage(directory string) *LocalStorage {
	if strings.TrimSpace(directory)=="" { directory="uploads" }
	return &LocalStorage{directory:directory}
}

func (s *LocalStorage) Save(ctx context.Context, header *multipart.FileHeader) (File,error) {
	if err:=ctx.Err();err!=nil{return File{},err}
	if err:=os.MkdirAll(s.directory,0o750);err!=nil{return File{},fmt.Errorf("create upload directory: %w",err)}
	source,err:=header.Open();if err!=nil{return File{},fmt.Errorf("open uploaded file: %w",err)}
	defer source.Close()
	random:=make([]byte,16);if _,err:=rand.Read(random);err!=nil{return File{},fmt.Errorf("create upload name: %w",err)}
	extension:=strings.ToLower(filepath.Ext(filepath.Base(header.Filename)))
	storedName:=hex.EncodeToString(random)+extension
	targetPath:=filepath.Join(s.directory,storedName)
	target,err:=os.OpenFile(targetPath,os.O_WRONLY|os.O_CREATE|os.O_EXCL,0o640)
	if err!=nil{return File{},fmt.Errorf("create upload target: %w",err)}
	prefix:=make([]byte,512);count,readErr:=source.Read(prefix);if readErr!=nil && readErr!=io.EOF { _=target.Close();_=os.Remove(targetPath);return File{},fmt.Errorf("inspect uploaded file: %w",readErr) }
	written,copyErr:=io.Copy(target,io.MultiReader(bytes.NewReader(prefix[:count]),source));closeErr:=target.Close()
	if copyErr!=nil { _=os.Remove(targetPath); return File{},fmt.Errorf("store uploaded file: %w",copyErr) }
	if closeErr!=nil { _=os.Remove(targetPath); return File{},fmt.Errorf("close uploaded file: %w",closeErr) }
	return File{OriginalName:filepath.Base(header.Filename),StoredName:storedName,Size:written,ContentType:http.DetectContentType(prefix[:count])},nil
}
`

var UploadService = `package uploads

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"{{ .ModuleName }}/pkg/upload"
)

var (
	ErrFileRequired=errors.New("file is required")
	ErrFileTooLarge=errors.New("file exceeds maximum size")
	ErrFileTypeNotAllowed=errors.New("file type is not allowed")
)

type Service struct { storage upload.Storage; maxSize int64; allowed map[string]struct{} }

func NewService(storage upload.Storage,maxSize int64,allowedExtensions []string) *Service {
	allowed:=make(map[string]struct{},len(allowedExtensions))
	for _,extension:=range allowedExtensions { extension=strings.ToLower(strings.TrimSpace(extension)); if extension==""{continue}; if !strings.HasPrefix(extension,"."){extension="."+extension}; allowed[extension]=struct{}{} }
	return &Service{storage:storage,maxSize:maxSize,allowed:allowed}
}

func (s *Service) Upload(ctx context.Context,file *multipart.FileHeader) (upload.File,error) {
	if file==nil{return upload.File{},ErrFileRequired}
	if s.maxSize>0 && file.Size>s.maxSize{return upload.File{},ErrFileTooLarge}
	extension:=strings.ToLower(filepath.Ext(filepath.Base(file.Filename)))
	if len(s.allowed)>0 { if _,ok:=s.allowed[extension];!ok{return upload.File{},ErrFileTypeNotAllowed} }
	return s.storage.Save(ctx,file)
}
`

var UploadHandler = `package uploads

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"{{ .ModuleName }}/pkg/response"
	service "{{ .ModuleName }}/services/uploads"
)

type Handler struct{ service *service.Service }
func NewHandler(service *service.Service) *Handler{return &Handler{service:service}}

func (h *Handler) Upload(c *fiber.Ctx) error {
	file,err:=c.FormFile("file")
	if err!=nil{return fiber.NewError(fiber.StatusBadRequest,"multipart field 'file' is required")}
	result,err:=h.service.Upload(c.UserContext(),file)
	if err!=nil {
		if errors.Is(err,service.ErrFileRequired)||errors.Is(err,service.ErrFileTooLarge)||errors.Is(err,service.ErrFileTypeNotAllowed){return fiber.NewError(fiber.StatusBadRequest,err.Error())}
		return err
	}
	return response.Send(c,fiber.StatusCreated,result)
}
`

var UploadRoute = `package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	handler "{{ .ModuleName }}/handlers/uploads"
	"{{ .ModuleName }}/pkg/upload"
	service "{{ .ModuleName }}/services/uploads"
)

func RegisterUploadRoutes(_ *gorm.DB,route fiber.Router) {
	storage:=upload.NewLocalStorage(viper.GetString("UPLOAD_DIR"))
	uploadService:=service.NewService(storage,viper.GetInt64("UPLOAD_MAX_SIZE"),strings.Split(viper.GetString("UPLOAD_ALLOWED_EXTENSIONS"),","))
	uploadHandler:=handler.NewHandler(uploadService)
	route.Post("/uploads",uploadHandler.Upload)
}
`
