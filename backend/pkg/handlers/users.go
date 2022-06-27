package handlers

import (
	"net/http"
	"social-network/pkg/models"
	"social-network/pkg/utils"
	"strings"
)

// Find all users and they relation with current user
func (handler *Handler) AllUsers(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request all users exccept current + relations
	users, errUsers := handler.repos.UserRepo.GetAllAndFollowing(userId)
	if errUsers != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, users, 200)
}

// Returns user nickname, id and path to avatar
func (handler *Handler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	user, err := handler.repos.UserRepo.GetDataMin(userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, []models.User{user}, 200)
}

// Find all followers
func (handler *Handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request all  following users
	followers, errUsers := handler.repos.UserRepo.GetFollowers(userId)
	if errUsers != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, followers, 200)
}

// Find all who clinet is following
func (handler *Handler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request all  following users
	followers, errUsers := handler.repos.UserRepo.GetFollowing(userId)
	if errUsers != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, followers, 200)
}

// Returns user data based on public / private profile and user_id from request
// waits for GET request with query "userId" ->user client is looking for
//  can be used both on own profile and other users
func (handler *Handler) UserData(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	currentUserId := r.Context().Value(utils.UserKey).(string)
	// get user id from request
	query := r.URL.Query()
	userId := query.Get("userId")
	// get if profile public or private
	status, err := handler.repos.UserRepo.ProfileStatus(userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// check if client looking for own profile
	currentUser := (currentUserId == userId)
	var following bool
	if !currentUser {
		// check if current user following user he is looking for
		following, err = handler.repos.UserRepo.IsFollowing(userId, currentUserId)
		if err != nil {
			utils.RespondWithError(w, "Error on getting data", 200)
			return
		}
	}
	// request user info based on following/ and profile status
	// if public or current user or if following  get large data set
	// if private and not following => get small data set
	var user models.User
	if currentUser || following || status == "PUBLIC" { // get full data set
		user, err = handler.repos.UserRepo.GetProfileMax(userId)
	} else {
		user, err = handler.repos.UserRepo.GetProfileMin(userId)
	}
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// tie together stats to user object
	user.Following = following
	user.CurrentUser = currentUser
	user.Status = status

	utils.RespondWithUsers(w, []models.User{user}, 200)
}

// changes user status in db return status
// in case of turning to PUBLIC -> also accept follow requests
func (handler *Handler) UserStatus(w http.ResponseWriter, r *http.Request) {
	statusList := []string{"PUBLIC", "PRIVATE"} //possible status
	var client models.User

	w = utils.ConfigHeader(w)
	// access user id
	client.ID = r.Context().Value(utils.UserKey).(string)
	// get status from request
	query := r.URL.Query()
	reqStatus := strings.ToUpper(query.Get("status"))

	// check if valid value and asign to user
	if reqStatus == statusList[0] {
		client.Status = statusList[0]
	} else if reqStatus == statusList[1] {
		client.Status = statusList[1]
	} else {
		utils.RespondWithError(w, "Requested status not valid", 200)
		return
	}
	// request current status from db
	currentStatus, err := handler.repos.UserRepo.GetStatus(client.ID)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// check if requested status is not the same as current
	if currentStatus == client.Status {
		utils.RespondWithError(w, "Status change not valid", 200)
		return
	}
	// Set new status
	err = handler.repos.UserRepo.SetStatus(client)
	if err != nil {
		utils.RespondWithError(w, "Error on saving status", 200)
		return
	}
	// if new status is public -> also accept pending follow requests
	// responds with success and newly created status
	utils.RespondWithSuccess(w, client.Status, 200)
}
