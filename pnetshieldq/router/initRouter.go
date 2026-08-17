package router

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net"
	"net/http"
	"pnetshieldq/controller"
	"pnetshieldq/utils"
)

func InitRouter(userController *controller.UserController, friendshipController *controller.FriendshipController, messageController *controller.MessageController, adminController *controller.AdminController, websocketutil *utils.MQTTWebSocketUtil, fileController *controller.FileController, userkeycontroller *controller.UserKeyController, logController *controller.LogController, messagerecordController *controller.MessageRecordController, configController *controller.ConfigController, friendcController *controller.FriendCController, useravatarcontroller *controller.UserAvatarController, serverStatusController *controller.ServerStatusController) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/update-avatar", useravatarcontroller.UpdateAvatarHandler).Methods("POST")
	router.HandleFunc("/get-avatar", useravatarcontroller.GetAvatarHandler).Methods("GET")
	router.HandleFunc("/get-friavatar/{userID}", useravatarcontroller.GetAvatarHandlerFriend).Methods("GET")
	router.HandleFunc("/server-status", serverStatusController.GetServerStatus).Methods("GET")
	router.HandleFunc("/logcount", logController.HandleLogCount).Methods("GET")
	router.HandleFunc("/usercount", userController.GetUserCount).Methods("GET")
	router.HandleFunc("/infor", messageController.GetMessageStatistics).Methods("GET")

	router.HandleFunc("/friend/add", friendshipController.AddFriend).Methods("POST")
	router.HandleFunc("/friend/remove", friendshipController.RemoveFriend).Methods("POST")
	router.HandleFunc("/message/send", messageController.SendMessage).Methods("POST")
	router.HandleFunc("/register", userController.Register).Methods("POST")
	router.HandleFunc("/login", userController.Login).Methods("POST")
	router.HandleFunc("/getfriends", friendshipController.GetFriendsList).Methods("GET")
	router.HandleFunc("/friend/accept", friendshipController.AcceptFriendRequest).Methods("POST")
	router.HandleFunc("/friend/requests", friendshipController.GetFriendRequests).Methods("GET")
	router.HandleFunc("/remove-friend", friendshipController.RemoveFriend).Methods("POST")

	router.HandleFunc("/admin/login", adminController.AdminLogin).Methods("POST")
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocketutil.HandleWebSocketConnection(w, r)
	})

	router.HandleFunc("/file/upload", fileController.UploadFile).Methods("POST")
	router.HandleFunc("/file/delete", fileController.DeleteFile).Methods("GET")
	router.HandleFunc("/message/history/{userID}/{friendID}", messageController.GetHistoricalMessages).Methods("GET")
	router.HandleFunc("/message/send", messageController.SendMessage).Methods("POST")
	router.HandleFunc("/file/upload", fileController.UploadFile).Methods("POST")
	router.HandleFunc("/file/delete", fileController.DeleteFile).Methods("DELETE")
	router.HandleFunc("/user/sendpk", userkeycontroller.CreateKey).Methods("POST")
	router.HandleFunc("/userkeys", userkeycontroller.GetAllKeys).Methods("GET")
	router.HandleFunc("/userkeys/retrieve", userkeycontroller.RetrieveKey).Methods("GET")
	router.HandleFunc("/friendc/add", friendcController.CreateFriendC).Methods("POST")
	router.HandleFunc("/friendc/get", friendcController.GetFriendCByFriendIDAndUserID).Methods("GET")

	router.HandleFunc("/user/get", userController.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", userController.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users", userController.GetAllUsers).Methods("GET")
	router.HandleFunc("/file/userfiles", fileController.GetUserFiles).Methods("GET")
	router.HandleFunc("/file/download", fileController.DownloadFile).Methods("GET")
	router.HandleFunc("/config", configController.GetConfig).Methods("GET")
	router.HandleFunc("/log", logController.HandleLog).Methods("GET")
	router.HandleFunc("/messagerecord", messagerecordController.FindAllMessages).Methods("GET")
	router.HandleFunc("/config", configController.UpdateConfig).Methods("POST")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/get-ip/", getIPHandler).Methods("GET")

	return handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router)
}
func getIPHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Fprintf(w, "Error getting IP: %s", err)
		return
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		fmt.Fprintf(w, "Unable to parse user IP")
		return
	}

	if userIP.IsLoopback() {
		fmt.Fprintf(w, "Loopback addresses are not processed")
		return
	}

	if userIP.IsPrivate() {
		fmt.Fprintf(w, "Private or local network addresses are not processed")
		return
	}

	fmt.Fprintf(w, "%s", userIP.String())
}
