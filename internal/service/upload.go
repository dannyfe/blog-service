package service

import (
	"errors"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/pkg/upload"
	"mime/multipart"
	"os"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

func (svc *Service) UploadFile(fileType upload.FileType, file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	fileName := upload.GetFileName(fileHeader.Filename)
	//检查文件格式
	if !upload.CheckContainExt(fileType, fileName) {
		return nil, errors.New("file suffix is not supported. ")
	}

	//检查文件大小
	if upload.CheckMaxSize(fileType, file) {
		return nil, errors.New("exceeded maximum file limit. ")
	}

	//检查文件存放路径是否存在，不存在就创建路径
	uploadSavePath := upload.GetSavePath()
	if upload.CheckSavePath(uploadSavePath) {
		if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("failed to create directory. ")
		}
	}

	//检查文件存放路径是否满足权限
	if upload.CheckPermission(uploadSavePath) {
		return nil, errors.New("insufficient file permissions. ")
	}

	//检查完成开始上传文件
	dst := uploadSavePath + "/" + fileName
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}

	accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}
