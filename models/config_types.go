package models

type Config struct {
	Api struct {
		Port           string `yaml:"port"`
		CollectionName string `yaml:"api_collection"`
	} `yaml:"base_api"`

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
		Endpoint       string `yaml:"endpoint"`
	} `yaml:"embedding"`

	Ollama struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"ollama"`

	Generator struct {
		ModelName string `yaml:"model_name"`
		Endpoint  string `yaml:"endpoint"`
	} `yaml:"generator"`

	Evaluation struct {
		RetrievalDataPath  string `yaml:"retrieval_data_path"`
		GenerationDataPath string `yaml:"generation_data_path"`
		SourceDataPath     string `yaml:"source_data_path"`
		CollectionName     string `yaml:"collection_name"`
	} `yaml:"evaluation"`
}
