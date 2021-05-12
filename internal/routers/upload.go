package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/service"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/convert"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
	"github.com/go-programming-tour-book/blog-service/pkg/upload"
)

type Upload struct {}

func NewUpload() Upload {
	return Upload{}
}

func (u Upload) UploadFile(c *gin.Context)  {
	response := app.NewResponese(c)
	file, fileHeader, err := c.Request.FormFile("file") //读取入参file字段的上传文件信息
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	fileType := convert.StrTo(c.PostForm("type")).MustInt() //利用入参type字段确立文件上传类型（也可以通过解析上传文件后缀来确定文件类型）
	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	response.ToResponse(gin.H{
		"file_access_rul": fileInfo.AccessUrl,
	})
}
