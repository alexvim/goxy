package configreader

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type JsonConfig struct {
	Proxy    Proxy    `json:"proxy"`
	Host     Host     `json:"host"`
	Resolver Resolver `json:"resolver"`
}

type Proxy struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type Host struct {
	Address string `json:"address"`
}

type Resolver struct {
	DoH DoH `json:"doh"`
}

type DoH struct {
	URL string `json:"url"`
}

func ReadFromFile(fileName string) JsonConfig {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("cfgreader: failed to open %s err=%s\n", fileName, err)
		return JsonConfig{}
	}

	cfg, err := readFromJson(file)
	if err != nil {
		log.Printf("cfgreader: failed to read file err=%s\n", err)
		return JsonConfig{}
	}

	return cfg
}

func readFromJson(r io.Reader) (JsonConfig, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		log.Printf("cfgreader: failed reads json err=%s\n", err)
		return JsonConfig{}, ErrConfigFileParse
	}

	cfg := JsonConfig{}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return JsonConfig{}, ErrConfigFileParse
	}

	return cfg, nil
}
