package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"pnetshieldq/config" // replace with your actual config package path
)

type ConfigController struct {
}

func NewConfigController() *ConfigController {
	return &ConfigController{}
}

func (ctrl *ConfigController) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Read error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var newConfig map[string]interface{}
	if err := json.Unmarshal(body, &newConfig); err != nil {
		http.Error(w, "Unmarshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	currentConfig, err := config.LoadConfig("E:\\netshieldq\\pnetshieldq\\config\\config.json")
	if err != nil {
		http.Error(w, "Load config error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var currentConfigMap map[string]interface{}
	inrec, err := json.Marshal(currentConfig)
	if err != nil {
		http.Error(w, "Error marshaling current config: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(inrec, &currentConfigMap); err != nil {
		http.Error(w, "Error unmarshaling to map: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for key, value := range newConfig {
		currentConfigMap[key] = value
	}

	var mergedConfig config.Config
	inrec, err = json.Marshal(currentConfigMap)
	if err != nil {
		http.Error(w, "Error marshaling merged config: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(inrec, &mergedConfig); err != nil {
		http.Error(w, "Error unmarshaling to config struct: "+err.Error(), http.StatusInternalServerError)
		return
	}

	configFile, err := os.OpenFile("E:\\netshieldq\\pnetshieldq\\config\\config.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		http.Error(w, "File open error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer configFile.Close()

	if err := json.NewEncoder(configFile).Encode(mergedConfig); err != nil {
		http.Error(w, "Encode error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Config updated successfully"))
}

func (ctrl *ConfigController) GetConfig(w http.ResponseWriter, r *http.Request) {
	currentConfig, err := config.LoadConfig("E:\\netshieldq\\pnetshieldq\\config\\config.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	configJson, err := json.Marshal(currentConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(configJson)
}
