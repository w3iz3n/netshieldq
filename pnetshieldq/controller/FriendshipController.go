package controller

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net/http"
	"pnetshieldq/dao"
	"pnetshieldq/entity"
	"pnetshieldq/log"
	"pnetshieldq/utils"
	"strconv"
)

type FriendshipController struct {
	FriendshipDao *dao.FriendshipDao
	UserDao       *dao.UserDao
	mqttClient    mqtt.Client
}

func NewFriendshipController(friendshipDao *dao.FriendshipDao, userDao *dao.UserDao, mqttClient mqtt.Client) *FriendshipController {
	return &FriendshipController{
		FriendshipDao: friendshipDao,
		UserDao:       userDao,
		mqttClient:    mqttClient,
	}
}

func (fc *FriendshipController) AddFriend(w http.ResponseWriter, r *http.Request) {

	var requestData struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	frienduser, err := fc.UserDao.FindByUsername(requestData.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	userid := utils.GetUserIDFromJWT(r)
	log := log.GetLogger()
	log.Info(strconv.Itoa(userid) + " Add " + strconv.Itoa(frienduser.UserID) + "as friend")
	exists, err := fc.FriendshipDao.CheckFriendship(userid, frienduser.UserID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Friend request already exists or already friends", http.StatusBadRequest)
		return
	}

	friendship := &entity.Friendship{
		UserID1: userid,
		UserID2: frienduser.UserID,
		Status:  "pending",
	}
	if err := fc.FriendshipDao.AddFriend(friendship); err != nil {
		http.Error(w, "Failed to create friend request", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Friend request sent successfully"))
}

func (fc *FriendshipController) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID1 int `json:"user_id1"`
		UserID2 int `json:"user_id2"`
	}
	log := log.GetLogger()
	log.Info("RemoveFriend method czalled")
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//exists, err := fc.FriendshipDao.CheckFriendship(data.UserID1, data.UserID2)
	//if err != nil {
	//	http.Error(w, "Database error", http.StatusInternalServerError)
	//	return
	//}
	//if !exists {
	//	http.Error(w, "Friendship does not exist", http.StatusBadRequest)
	//	return
	//}

	// 删除好友关系
	if err := fc.FriendshipDao.RemoveFriend(data.UserID1, data.UserID2); err != nil {
		http.Error(w, "Failed to remove friend", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Friend removed successfully"))
}

func (controller *FriendshipController) GetFriendsList(w http.ResponseWriter, r *http.Request) {
	log := log.GetLogger()
	log.Info("GetFriendsList method called")
	userID := utils.GetUserIDFromJWT(r)
	if userID == -1 {
		http.Error(w, "Invalid token or no token provided", http.StatusUnauthorized)
		return
	}

	friendIDs, err := controller.FriendshipDao.GetFriends(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve friends", http.StatusInternalServerError)
		return
	}

	friends := make([]entity.FriendInfo, 0, len(friendIDs)) // 使用 FriendInfo 列表
	for _, id := range friendIDs {
		friend, err := controller.UserDao.FindByID(id)
		if err != nil {
			continue
		}
		friends = append(friends, entity.FriendInfo{ID: friend.UserID, Username: friend.Username})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(friends); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (fc *FriendshipController) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	log := log.GetLogger()
	var requestData struct {
		Username string `json:"username"`
		Status   string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	friendUser, err := fc.UserDao.FindByUsername(requestData.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userID := utils.GetUserIDFromJWT(r)
	if userID == -1 {
		http.Error(w, "Invalid token or no token provided", http.StatusUnauthorized)
		return
	}
	log.Info(strconv.Itoa(userID) + " AcceptFriendRequest" + strconv.Itoa(friendUser.UserID))
	err = fc.FriendshipDao.UpdateStatus(userID, friendUser.UserID, requestData.Status)
	if err != nil {
		http.Error(w, "Failed to accept friend request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Friend request accepted successfully"))
}

func (fc *FriendshipController) GetFriendRequests(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromJWT(r)
	if userID == -1 {
		http.Error(w, "Invalid token or no token provided", http.StatusUnauthorized)
		return
	}
	log := log.GetLogger()
	log.Info("GetFriendRequests method called")
	friendRequests, err := fc.FriendshipDao.GetFriendRequests(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve friend requests", http.StatusInternalServerError)
		return
	}

	friendInfos := make([]map[string]interface{}, 0, len(friendRequests))
	for _, request := range friendRequests {
		friend, err := fc.UserDao.FindByID(request.UserID1)
		if err != nil {
			continue
		}
		friendInfos = append(friendInfos, map[string]interface{}{
			"username": friend.Username,
			"status":   request.Status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(friendInfos); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
