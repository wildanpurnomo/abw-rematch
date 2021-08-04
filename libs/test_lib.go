package libs

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockObject struct {
	mock.Mock
}

func (s *MockObject) UploadFile(fileHeader *multipart.FileHeader, bucketName string) error {
	return nil
}

func InitGinForTesting() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}
