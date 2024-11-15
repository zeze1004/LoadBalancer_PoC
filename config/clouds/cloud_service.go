package clouds

import (
	"github.com/joho/godotenv"
	"os"
)

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
		panic(".env이 없습니다")
	}

	return map[string]string{
		"AWS":   os.Getenv("AWS_API_KEY"),
		"Azure": os.Getenv("AZURE_API_KEY"),
		"GCP":   os.Getenv("GCP_API_KEY"),
	}
}
