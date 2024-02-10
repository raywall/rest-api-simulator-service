package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"hawx.me/code/route"
)

var config settings

type settings struct {
	TemplateFormatVersion string              `yaml:"TemplateFormatVersion"`
	Description           string              `yaml:"Description"`
	Resources             map[string]resource `yaml:"Resources"`
}

type resource struct {
	ResourceType string     `yaml:"ResourceType"`
	Description  string     `yaml:"Description,omitempty"`
	Properties   properties `yaml:"Properties"`
}

type properties struct {
	Path     string   `yaml:"Path"`
	Method   string   `yaml:"Method"`
	Response response `yaml:"Response"`
}

type response struct {
	StatusCode int    `yaml:"StatusCode"`
	Content    string `yaml:"Content"`
}

func (res *resource) getResponse(vars map[string]string) (map[string]interface{}, error) {
	content, err := os.ReadFile(res.Properties.Response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed on get response body file: %v", err)
	}

	data := string(content)
	if len(vars) > 0 {
		for key, value := range vars {
			item := fmt.Sprintf("{%s}", key)
			if strings.Contains(data, item) {
				data = strings.ReplaceAll(data, item, value)
			}
		}
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(data), &jsonMap); err != nil {
		return nil, fmt.Errorf("failed marshal the json file: %v", err)
	}

	return jsonMap, nil
}

func (res resource) server(w http.ResponseWriter, r *http.Request) {
	logMessage := fmt.Sprintf("%s - Processing received request for %s", r.Method, r.RequestURI)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if res.Properties.Method == r.Method {
		data, err := res.getResponse(route.Vars(r))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.WriteHeader(res.Properties.Response.StatusCode)
		json.NewEncoder(w).Encode(data)

		log.Println(logMessage)
		return
	}

	logMessage += " (fail)"
	log.Println(logMessage)
	http.NotFound(w, r)
}

func init() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("falha ao carregar configurações: %v", err)
	}

	yaml.Unmarshal(data, &config)
	for _, res := range config.Resources {
		route.HandleFunc(res.Properties.Path, res.server)
	}
}

func main() {
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", route.Default); err != nil {
		log.Fatal(err)
	}
}
