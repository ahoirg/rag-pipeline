package models

type Config struct {
	Chunk struct {
		Size    int `yaml:"size"`
		Overlap int `yaml:"overlap"`
	} `yaml:"chunk"`

	Qdrant struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"qdrant"`

	Embedding struct {
		ModelDimension int    `yaml:"model_dimension"`
		ModelName      string `yaml:"model_name"`
	} `yaml:"embedding"`

	Ollama struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"ollama"`
}
