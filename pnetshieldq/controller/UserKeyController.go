package controller

import (
	"encoding/json"
	"net/http"
	"pnetshieldq/log"
	"pnetshieldq/utils"
	"strconv"

	"pnetshieldq/dao"
	"pnetshieldq/entity"
)

type UserKeyController struct {
	UserkeyDao *dao.UserKeyDAO
}

func NewUserKeyController(dao *dao.UserKeyDAO) *UserKeyController {
	return &UserKeyController{UserkeyDao: dao}
}

func (c *UserKeyController) CreateKey(w http.ResponseWriter, r *http.Request) {
	var userKey entity.UserKey
	var publicKey struct {
		PK string `json:"publicKey"`
	}
	userID := utils.GetUserIDFromJWT(r)
	if userID == -1 {
		http.Error(w, "Invalid token or no token provided", http.StatusUnauthorized)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&publicKey); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log := log.GetLogger()
	log.Info("User_" + strconv.Itoa(userID) + " Generate PK")
	userKey.UserID = uint(userID)
	userKey.PK = publicKey.PK
	if err := c.UserkeyDao.UpsertKey(&userKey); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userKey)
}

func (c *UserKeyController) RetrieveKey(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userKey, err := c.UserkeyDao.RetrieveKeyByID(uint(id))
	if err != nil {
		http.Error(w, "UserKey not found", http.StatusNotFound)
		return
	}
	log := log.GetLogger()
	log.Info(idStr + " is applied PK")
	json.NewEncoder(w).Encode(userKey.PK)
}
func (c *UserKeyController) GetAllKeys(w http.ResponseWriter, r *http.Request) {
	userKeys, err := c.UserkeyDao.GetAllKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(userKeys)
}
