package clouds

import (
	"github.com/joho/godotenv"
	"os"
)

// CloudService API 키를 관리하고자 작성한 구조체이지만, 로직이 구현되지는 않았습니다
type CloudService struct {
	Name   string
	Region []string
	APIKey string
}

func NewCloudService(name string, region []string, apiKey string) *CloudService {
	return &CloudService{
		Name:   name,
		Region: region,
		APIKey: apiKey,
	}
}

func LoadAPIKeys() map[string]string {
	err := godotenv.Load()
	if err != nil {
		panic(".env가 없습니다")
	}

	return map[string]string{
		"AWS":   os.Getenv("AWS_API_KEY"),
		"Azure": os.Getenv("AZURE_API_KEY"),
		"GCP":   os.Getenv("GCP_API_KEY"),
	}
}
