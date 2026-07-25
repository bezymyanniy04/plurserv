package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	// "sync/atomic"
	"time"

	"github.com/bezymyanniy04/plurserv/internal/auth"
	"github.com/bezymyanniy04/plurserv/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

//consts

const default_for_newbies string = `Hello! You might be feeling a little confused. Don't worry, that's normal. You're a new headmate in a plural system. To put it simply, you share a body with other people just like you. You're living a life together, so try to find common ground with them.
Well, first off, let's help you out. On the previous page you saw the "Headmates" button. You can go there and fill out your information.
What kind of information is this? Let me help you.
Your: appearance; name; pronouns; your role in the system; age. There will also be a free field where you can write anything you like about yourself. Below this there will be a button to go to your personal diary.
It's okay if you don't know or remember everything about yourself right away, take your time.
Good luck!`

//structs

type apiConfig struct {
	db        *database.Queries
	platform  string
	JWTSecret string
	dbc       *sql.DB
}

type UserPass struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Avatar     string    `json:"avatar"`
	SystemName string    `json:"system_name"`
	Theme      int32     `json:"theme"`
	Font       string    `json:"font"`
}

type User struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	Avatar     string    `json:"avatar"`
	SystemName string    `json:"system_name"`
	Theme      int32     `json:"theme"`
	Font       string    `json:"font"`
	Token      string    `json:"token"`
	Refresh    string    `json:"refresh_token"`
}

type Alter struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Avatar        string    `json:"avatar"`
	Name          string    `json:"name"`
	Pronouns      string    `json:"pronouns"`
	Age           string    `json:"age"`
	Role          string    `json:"role"`
	Description   string    `json:"description"`
	Colour        string    `json:"colour"`
	Fronting      bool      `json:"fronting"`
	FrontingSince time.Time `json:"fronting_since"`
	User_id       uuid.UUID `json:"userId"`
}

type Front struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	EndedAt   time.Time `json:"ended_at"`
	AlterId   uuid.UUID `json:"alter_id"`
}

type FriendRequest struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
	SenderId   uuid.UUID `json:"sender_id"`
	RecieverId uuid.UUID `json:"reciever_id"`
	Answer     int32     `json:"answer"`
}

type Diary struct {
	ID          uuid.UUID `json:"id"`
	AlterId     uuid.UUID `json:"alter_id"`
	BgColour    string    `json:"bg_colour"`
	BgColour2   string    `json:"bg_colour2"`
	TextColour  string    `json:"text_colour"`
	BlockColour string    `json:"block_colour"`
	Font        string    `json:"font"`
	UserId      uuid.UUID `json:"user_id"`
	Name        string    `json:"alter_name"`
	AlterAvatar string    `json:"alter_avatar"`
}

type DiaryEntry struct {
	ID      uuid.UUID `json:"id"`
	DiaryId uuid.UUID `json:"diary_id"`
	Name    string    `json:"name"`
	Date    time.Time `json:"date"`
	Text    string    `json:"text"`
	UserId  uuid.UUID `json:"user_id"`
}

type DiaryEntryFile struct {
	ID        uuid.UUID `json:"id"`
	EntryId   uuid.UUID `json:"entry_id"`
	File      string    `json:"file"`
	UserId    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ForNewbies struct {
	ID     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"user_id"`
	Text   string    `json:"text"`
}
type FriendReq struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	Avatar     string    `json:"avatar"`
	SystemName string    `json:"system_name"`
	RequestID  uuid.UUID `json:"request_id"`
	Answer     int32     `json:"answer"`
}
type AvatarPath struct {
	AvatarPath string `json:"path"`
}

// type FrontingAlter struct {
// 	CreatedAt time.Time `json:"created_at"`
// 	EndedAt   time.Time `json:"ended_at"`
// 	AlterId   uuid.UUID `json:"alter_id"`
// 	Name      string    `json:"name"`
// 	User_id   uuid.UUID `json:"userId"`
// }

//middleware

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

//helpers

func marsh[str any](bod str, code int, w http.ResponseWriter) {
	dat, err := json.Marshal(bod)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error marshalling JSON: %s", err)))
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func hex_to_rgb(hex string, w http.ResponseWriter) string {
	colour := ""
	for i := range 3 {
		c, err := strconv.ParseInt(hex[2*i+1:2*i+3], 16, 0)

		if err != nil {
			err_mes("bad colour", 400, w)
			return ""
		}
		cc := strconv.Itoa(int(c))
		colour += cc + " "
	}
	return colour
}

//json

func err_mes(err string, code int, w http.ResponseWriter) {
	type e struct {
		Error string `json:"error"`
	}
	bod := e{Error: err}
	marsh(bod, code, w)
}

//requests

func Readiness(w http.ResponseWriter, r *http.Request) {
	// w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// w.WriteHeader(200)
	// w.Write([]byte("OK\n"))
	type try struct {
		Text  string `json:"text"`
		ID    string `json:"id"`
		Image string `json:"image"`
	}
	t := try{
		Text:  "aboba",
		ID:    "some id",
		Image: "https://avatars.mds.yandex.net/i?id=e66ce3e9b4ad13f1e6d3fd1c539de1ea7e157ea1-5130204-images-thumbs&n=13",
	}
	marsh(t, 200, w)
}

//users

func (cfg *apiConfig) post_user(w http.ResponseWriter, r *http.Request) {
	param := UserPass{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}
	if param.Email == "" {
		err_mes("Empty email", 400, w)
		return
	}
	if param.Password == "" {
		err_mes("Empty password", 400, w)
		return
	}
	if param.Avatar == "" {
		param.Avatar = "assets/default_avater.jpg"
	}
	if param.SystemName == "" {
		param.SystemName = "System"
	}
	if param.Font == "" {
		param.Font = "Arial"
	}
	hashed_password, err := auth.HashPassword(param.Password)
	if err != nil {
		err_mes("Something went wrong with hash", 400, w)
		return
	}

	userdb, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: param.Email, HashedPassword: hashed_password, Avatar: param.Avatar, SystemName: param.SystemName, Theme: param.Theme, Font: param.Font})
	if err != nil {
		err_mes(fmt.Sprintf("%v", err), 400, w)
		return
	}
	_, err = cfg.db.NewForNewbies(r.Context(), database.NewForNewbiesParams{UserID: userdb.ID, Text: default_for_newbies})
	if err != nil {
		err_mes("Something went wrong with db2", 400, w)
		return
	}

	user := User{
		ID:         userdb.ID,
		CreatedAt:  userdb.CreatedAt,
		UpdatedAt:  userdb.UpdatedAt,
		Email:      userdb.Email,
		Avatar:     userdb.Avatar,
		SystemName: userdb.SystemName,
		Theme:      userdb.Theme,
		Font:       userdb.Font,
	}
	marsh(user, 201, w)

}

func (cfg *apiConfig) get_user(w http.ResponseWriter, r *http.Request) {

	userIdString := r.PathValue("userId")
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid_token, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	if userId != userid_token {
		_, err = cfg.db.GetIfFriends(r.Context(), database.GetIfFriendsParams{UserID: userid_token, FriendID: userId})
		if err != nil {
			err_mes("You're not friends", 401, w)
			return
		}
	}

	userdb, err := cfg.db.GetUser(r.Context(), userId)
	if err != nil {
		err_mes("failed to get alter", 404, w)
		return
	}
	user := User{
		ID:         userdb.ID,
		CreatedAt:  userdb.CreatedAt,
		UpdatedAt:  userdb.UpdatedAt,
		Email:      userdb.Email,
		Avatar:     userdb.Avatar,
		SystemName: userdb.SystemName,
		Theme:      userdb.Theme,
		Font:       userdb.Font,
	}
	marsh(user, 200, w)

}

func (cfg *apiConfig) edit_userinfo(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	param := UserPass{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}
	if param.Email == "" {
		err_mes("Empty email", 400, w)
		return
	}
	if param.Password == "" {
		err_mes("Empty password", 400, w)
		return
	}
	if param.Avatar == "" {
		param.Avatar = "assets/default_avater.jpg"
	}
	if param.SystemName == "" {
		param.SystemName = "System"
	}
	if param.Font == "" {
		param.Font = "Arial"
	}
	hashed_password, err := auth.HashPassword(param.Password)
	if err != nil {
		err_mes("Something went wrong with hash", 400, w)
		return
	}
	userinfodb, err := cfg.db.EditUserInfo(r.Context(), database.EditUserInfoParams{Email: param.Email, HashedPassword: hashed_password, SystemName: param.SystemName, Theme: param.Theme, Font: param.Font, ID: userid})
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}
	userinfo := User{
		ID:         userinfodb.ID,
		CreatedAt:  userinfodb.CreatedAt,
		UpdatedAt:  userinfodb.UpdatedAt,
		Email:      userinfodb.Email,
		Avatar:     userinfodb.Avatar,
		SystemName: userinfodb.SystemName,
		Theme:      userinfodb.Theme,
		Font:       userinfodb.Font,
	}
	marsh(userinfo, 201, w)
}

func (cfg *apiConfig) edit_userinfo_settings(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	param := UserPass{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}

	if param.Avatar == "" {
		param.Avatar = "assets/default_avater.jpg"
	}
	if param.SystemName == "" {
		param.SystemName = "System"
	}
	if param.Font == "" {
		param.Font = "Arial"
	}

	userinfodb, err := cfg.db.EditUserSettings(r.Context(), database.EditUserSettingsParams{SystemName: param.SystemName, Theme: param.Theme, Font: param.Font, ID: userid})
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}
	userinfo := User{
		ID:         userinfodb.ID,
		CreatedAt:  userinfodb.CreatedAt,
		UpdatedAt:  userinfodb.UpdatedAt,
		Email:      userinfodb.Email,
		Avatar:     userinfodb.Avatar,
		SystemName: userinfodb.SystemName,
		Theme:      userinfodb.Theme,
		Font:       userinfodb.Font,
	}
	marsh(userinfo, 201, w)
}

//alters

func (cfg *apiConfig) post_alter(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	param := Alter{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong with decoding", 400, w)
		return
	}
	if param.Name == "" {
		err_mes("Empty name", 400, w)
		return
	}
	colour := ""
	if param.Colour == "" {
		colour = "102 102 102"
	} else {

		for i := range 3 {
			parsed, err := strconv.ParseInt(param.Colour[2*i+1:2*i+3], 16, 0)
			if err != nil {
				err_mes("Bad colour", 400, w)
				return
			}

			colour += strconv.FormatInt(parsed, 10) + " "
		}
	}

	if param.Avatar == "" {
		param.Avatar = "assets/default_avater.jpg"
	}
	alterdb, err := cfg.db.CreateAlter(r.Context(), database.CreateAlterParams{Avatar: param.Avatar, Name: param.Name, Pronouns: param.Pronouns, Age: param.Age, AlterRole: param.Role, Description: param.Description, Colour: colour, UserID: userid})
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}
	alter := Alter{
		ID:          alterdb.ID,
		CreatedAt:   alterdb.CreatedAt,
		UpdatedAt:   alterdb.UpdatedAt,
		Avatar:      alterdb.Avatar,
		Name:        alterdb.Name,
		Pronouns:    alterdb.Pronouns,
		Age:         alterdb.Age,
		Role:        alterdb.AlterRole,
		Description: alterdb.Description,
		Colour:      alterdb.Colour,
		Fronting:    alterdb.Fronting,
		User_id:     alterdb.UserID,
	}
	marsh(alter, 201, w)
}

func (cfg *apiConfig) get_alters(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("user_id")
	// altersdb, err := cfg.db.GetAlters(r.Context())
	if s == "" {
		err_mes("user uuid incorrect", 400, w)
		return
	}

	query := r.URL.Query().Get("query")
	query = "%" + strings.ToLower(query) + "%"
	id, err := uuid.Parse(s)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid_token, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	if id != userid_token {
		_, err = cfg.db.GetIfFriends(r.Context(), database.GetIfFriendsParams{UserID: userid_token, FriendID: id})
		if err != nil {
			err_mes("You're not friends", 401, w)
			return
		}
	}

	altersdb, err := cfg.db.GetAltersByAuthor(r.Context(), database.GetAltersByAuthorParams{UserID: id, Name: query})
	if err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}

	sort.Slice(altersdb, func(i, j int) bool {
		return strings.ToLower(altersdb[i].Name) < strings.ToLower(altersdb[j].Name)
	})

	alters := []Alter{}
	for _, alterdb := range altersdb {
		alter := Alter{
			ID:          alterdb.ID,
			CreatedAt:   alterdb.CreatedAt,
			UpdatedAt:   alterdb.UpdatedAt,
			Name:        alterdb.Name,
			Avatar:      alterdb.Avatar,
			Pronouns:    alterdb.Pronouns,
			Age:         alterdb.Age,
			Role:        alterdb.AlterRole,
			Description: alterdb.Description,
			Colour:      alterdb.Colour,
			Fronting:    alterdb.Fronting,
			User_id:     alterdb.UserID,
		}
		alters = append(alters, alter)
	}
	marsh(alters, 200, w)
}

func (cfg *apiConfig) get_alters_without_diary(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	altersdb, err := cfg.db.GetAltersWithoutDiaries(r.Context(), userid)
	if err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}

	sort.Slice(altersdb, func(i, j int) bool {
		return strings.ToLower(altersdb[i].Name) < strings.ToLower(altersdb[j].Name)
	})

	alters := []Alter{}
	for _, alterdb := range altersdb {
		alter := Alter{
			ID:          alterdb.ID,
			CreatedAt:   alterdb.CreatedAt,
			UpdatedAt:   alterdb.UpdatedAt,
			Name:        alterdb.Name,
			Avatar:      alterdb.Avatar,
			Pronouns:    alterdb.Pronouns,
			Age:         alterdb.Age,
			Role:        alterdb.AlterRole,
			Description: alterdb.Description,
			Colour:      alterdb.Colour,
			Fronting:    alterdb.Fronting,
			User_id:     alterdb.UserID,
		}
		alters = append(alters, alter)
	}
	marsh(alters, 200, w)
}

func (cfg *apiConfig) get_alter(w http.ResponseWriter, r *http.Request) {
	alterIdString := r.PathValue("alterId")
	alterId, err := uuid.Parse(alterIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	alterdb, err := cfg.db.GetAlter(r.Context(), alterId)
	if err != nil {
		err_mes("failed to get alter", 404, w)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid_token, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	if alterdb.UserID != userid_token {
		_, err = cfg.db.GetIfFriends(r.Context(), database.GetIfFriendsParams{UserID: userid_token, FriendID: alterdb.UserID})
		if err != nil {
			err_mes("You're not friends", 401, w)
			return
		}
	}

	alter := Alter{
		ID:          alterdb.ID,
		CreatedAt:   alterdb.CreatedAt,
		UpdatedAt:   alterdb.UpdatedAt,
		Name:        alterdb.Name,
		Avatar:      alterdb.Avatar,
		Pronouns:    alterdb.Pronouns,
		Age:         alterdb.Age,
		Role:        alterdb.AlterRole,
		Description: alterdb.Description,
		Colour:      alterdb.Colour,
		Fronting:    alterdb.Fronting,
		User_id:     alterdb.UserID,
	}
	marsh(alter, 200, w)
}

func (cfg *apiConfig) edit_alter(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	alterIdString := r.PathValue("alterId")
	alterId, err := uuid.Parse(alterIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	param := Alter{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong with decoding", 400, w)
		return
	}
	if param.Name == "" {
		err_mes("Empty name", 400, w)
		return
	}
	colour := ""
	if param.Colour == "" {
		colour = "102 102 102"
	} else {

		for i := range 3 {
			parsed, err := strconv.ParseInt(param.Colour[2*i+1:2*i+3], 16, 0)
			if err != nil {
				err_mes("Bad colour", 400, w)
				return
			}

			colour += strconv.FormatInt(parsed, 10) + " "
		}
	}

	if param.Avatar == "" {
		param.Avatar = "assets/default_avater.jpg"
	}
	alterdb, err := cfg.db.EditAlter(r.Context(), database.EditAlterParams{ID: alterId, UserID: userid, Name: param.Name, Pronouns: param.Pronouns, Age: param.Age, AlterRole: param.Role, Description: param.Description, Colour: colour})
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}
	alter := Alter{
		ID:          alterdb.ID,
		CreatedAt:   alterdb.CreatedAt,
		UpdatedAt:   alterdb.UpdatedAt,
		Name:        alterdb.Name,
		Avatar:      alterdb.Avatar,
		Pronouns:    alterdb.Pronouns,
		Age:         alterdb.Age,
		Role:        alterdb.AlterRole,
		Description: alterdb.Description,
		Colour:      alterdb.Colour,
		Fronting:    alterdb.Fronting,
		User_id:     alterdb.UserID,
	}
	marsh(alter, 201, w)
}

func (cfg *apiConfig) delete_alter(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	alterIdString := r.PathValue("alterId")
	alterId, err := uuid.Parse(alterIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	alterdb, err := cfg.db.GetAlter(r.Context(), alterId)
	if err != nil {
		err_mes("failed to get alter", 404, w)
		return
	}
	if alterdb.UserID != userid {
		err_mes("You're not the author", 403, w)
		return
	}
	cfg.db.DeleteAlter(r.Context(), alterId)
	w.WriteHeader(204)
}

//fronts

func (cfg *apiConfig) change_front(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	_, err = auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	alterIdString := r.PathValue("alterId")
	alterId, err := uuid.Parse(alterIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	alterdb := database.Alter{}
	frontdb, err := cfg.db.GetNowFrontByAlter(r.Context(), alterId)
	if err != nil {

		erri := cfg.db.NewFront(r.Context(), alterId)
		if erri != nil {
			err_mes("couldn't create new front", 400, w)
			return
		}
		alterdb, erri = cfg.db.NewFrontAlter(r.Context(), alterId)

	} else {

		erri := cfg.db.EndFrontFronts(r.Context(), frontdb.ID)
		if erri != nil {
			err_mes("couldn't end front", 400, w)
			return
		}
		alterdb, erri = cfg.db.EndFrontAlter(r.Context(), alterId)
	}

	alter := Alter{
		ID:          alterdb.ID,
		CreatedAt:   alterdb.CreatedAt,
		UpdatedAt:   alterdb.UpdatedAt,
		Name:        alterdb.Name,
		Avatar:      alterdb.Avatar,
		Pronouns:    alterdb.Pronouns,
		Age:         alterdb.Age,
		Role:        alterdb.AlterRole,
		Description: alterdb.Description,
		Colour:      alterdb.Colour,
		Fronting:    alterdb.Fronting,
		User_id:     alterdb.UserID,
	}
	marsh(alter, 201, w)
}

func (cfg *apiConfig) get_fronting_alters(w http.ResponseWriter, r *http.Request) {

	userIdString := r.PathValue("userId")
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid_token, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	query := r.URL.Query().Get("query")
	query = "%" + strings.ToLower(query) + "%"

	if userId != userid_token {
		_, err = cfg.db.GetIfFriends(r.Context(), database.GetIfFriendsParams{UserID: userid_token, FriendID: userId})
		if err != nil {
			err_mes("You're not friends", 401, w)
			return
		}
	}

	altersdb, err := cfg.db.GetAllNowFronts(r.Context(), database.GetAllNowFrontsParams{UserID: userId, Name: query})
	if err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}

	sort.Slice(altersdb, func(i, j int) bool {
		return strings.ToLower(altersdb[i].Name) < strings.ToLower(altersdb[j].Name)
	})
	alters := []Alter{}
	for _, alterdb := range altersdb {
		alter := Alter{
			ID:            alterdb.ID,
			CreatedAt:     alterdb.CreatedAt,
			UpdatedAt:     alterdb.UpdatedAt,
			Name:          alterdb.Name,
			Pronouns:      alterdb.Pronouns,
			Avatar:        alterdb.Avatar,
			Age:           alterdb.Age,
			Role:          alterdb.AlterRole,
			Description:   alterdb.Description,
			Colour:        alterdb.Colour,
			Fronting:      alterdb.Fronting,
			FrontingSince: alterdb.FrontingSince.Time,
			User_id:       alterdb.UserID,
		}
		alters = append(alters, alter)
	}
	marsh(alters, 200, w)

}

func (cfg *apiConfig) get_fronts_by_time(w http.ResponseWriter, r *http.Request) {

	type Body_front struct {
		UserId uuid.UUID `json:"user_id"`
		From   string    `json:"from"`
		To     string    `json:"to"`
	}

	param := Body_front{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong with decoding", 400, w)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid_token, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	if param.UserId != userid_token {
		// _, err = cfg.db.GetIfFriends(r.Context(), database.GetIfFriendsParams{UserID: userid_token, FriendID: param.UserId})
		// if err != nil {
		err_mes("You're not friends", 401, w)
		return
		// }
	}

	timefrom, err := time.Parse("2006-01-02", param.From)
	timeto, err := time.Parse("2006-01-02", param.To)
	altersdb, err := cfg.db.GetFrontsByTime(r.Context(), database.GetFrontsByTimeParams{StartedAt: timefrom, StartedAt_2: timeto, UserID: param.UserId})
	if err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}

	// order := r.URL.Query().Get("sort")
	// if order == "asc" {
	// 	sort.Slice(altersdb, func(i, j int) bool {
	// 		return strings.ToLower(altersdb[i].Name) < strings.ToLower(altersdb[j].Name)
	// 	})
	// }

	type AlterFront struct {
		ID          uuid.UUID `json:"id"`
		StartedAt   time.Time `json:"started_at"`
		EndedAt     time.Time `json:"ended_at"`
		Avatar      string    `json:"avatar"`
		Name        string    `json:"name"`
		Pronouns    string    `json:"pronouns"`
		Age         string    `json:"age"`
		Role        string    `json:"role"`
		Description string    `json:"description"`
		Colour      string    `json:"colour"`
		Fronting    bool      `json:"fronting"`
		User_id     uuid.UUID `json:"userId"`
	}

	alters := []AlterFront{}
	for _, alterdb := range altersdb {
		alter := AlterFront{
			ID:          alterdb.ID,
			StartedAt:   alterdb.StartedAt,
			EndedAt:     alterdb.EndedAt.Time,
			Name:        alterdb.Name,
			Pronouns:    alterdb.Pronouns,
			Avatar:      alterdb.Avatar,
			Age:         alterdb.Age,
			Role:        alterdb.AlterRole,
			Description: alterdb.Description,
			Colour:      alterdb.Colour,
			Fronting:    alterdb.Fronting,
			User_id:     alterdb.UserID,
		}
		alters = append(alters, alter)
	}
	marsh(alters, 200, w)
}

//friends

//requests

func (cfg *apiConfig) post_friend_request(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	receiverIdString := r.PathValue("receiverId")
	receiverId, err := uuid.Parse(receiverIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	_, err = cfg.db.GetFriendRequest(r.Context(), database.GetFriendRequestParams{SenderID: userid, ReceiverID: receiverId})
	if err == nil {
		err_mes("You've already sent friend request", 400, w)
		return
	}

	requestdb, err := cfg.db.NewFriendRequest(r.Context(), database.NewFriendRequestParams{SenderID: userid, ReceiverID: receiverId})
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}

	request := FriendRequest{
		ID:         requestdb.ID,
		CreatedAt:  requestdb.CreatedAt,
		ExpiredAt:  requestdb.ExpiresAt.Time,
		SenderId:   requestdb.SenderID,
		RecieverId: requestdb.ReceiverID,
		Answer:     requestdb.Answer,
	}
	marsh(request, 201, w)
}

func (cfg *apiConfig) get_friend_requests_as_sender(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	requestsdb, err := cfg.db.GetNewFriendRequestsSender(r.Context(), userid)
	if err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}

	order := r.URL.Query().Get("sort")
	if order == "asc" {
		sort.Slice(requestsdb, func(i, j int) bool { return requestsdb[i].CreatedAt.After(requestsdb[j].CreatedAt) })
	}

	requests := []FriendReq{}
	for _, requestdb := range requestsdb {
		request := FriendReq{
			ID:         requestdb.ID,
			CreatedAt:  requestdb.CreatedAt,
			UpdatedAt:  requestdb.UpdatedAt,
			Email:      requestdb.Email,
			Avatar:     requestdb.Avatar,
			SystemName: requestdb.SystemName,
			Answer:     requestdb.Answer,
			RequestID:  requestdb.ID_2,
		}
		requests = append(requests, request)
	}
	marsh(requests, 200, w)
}

func (cfg *apiConfig) get_friend_requests_as_receiver(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	requestsdb, err := cfg.db.GetNewFriendRequestsReciever(r.Context(), userid)
	if err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}

	order := r.URL.Query().Get("sort")
	if order == "asc" {
		sort.Slice(requestsdb, func(i, j int) bool { return requestsdb[i].CreatedAt.After(requestsdb[j].CreatedAt) })
	}

	requests := []FriendReq{}
	for _, requestdb := range requestsdb {
		request := FriendReq{
			ID:         requestdb.ID,
			CreatedAt:  requestdb.CreatedAt,
			UpdatedAt:  requestdb.UpdatedAt,
			Email:      requestdb.Email,
			Avatar:     requestdb.Avatar,
			SystemName: requestdb.SystemName,
			Answer:     requestdb.Answer,
			RequestID:  requestdb.ID_2,
		}
		requests = append(requests, request)
	}
	marsh(requests, 200, w)
}

func (cfg *apiConfig) put_answer_friend_request(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	param := FriendRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong with decoding", 400, w)
		return
	}

	requestdb, err := cfg.db.AnswerFriendRequest(r.Context(), database.AnswerFriendRequestParams{Answer: param.Answer, ID: param.ID, ReceiverID: userid})
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}
	if param.Answer == 1 {
		err = cfg.db.NewFriends(r.Context(), database.NewFriendsParams{RequestID: requestdb.ID, UserID: requestdb.SenderID, FriendID: requestdb.ReceiverID})
		if err != nil {
			err_mes("Something went wrong with friends db 1", 400, w)
			return
		}
		err = cfg.db.NewFriends(r.Context(), database.NewFriendsParams{RequestID: requestdb.ID, UserID: requestdb.ReceiverID, FriendID: requestdb.SenderID})
		if err != nil {
			err_mes("Something went wrong with friends db 2", 400, w)
			return
		}
	} else {
		err = cfg.db.DeleteFriendRequest(r.Context(), requestdb.ID)
		if err != nil {
			err_mes("Something went wrong with db3", 400, w)
			return
		}
	}

	request := FriendRequest{
		ID:         requestdb.ID,
		CreatedAt:  requestdb.CreatedAt,
		ExpiredAt:  requestdb.ExpiresAt.Time,
		SenderId:   requestdb.SenderID,
		RecieverId: requestdb.ReceiverID,
		Answer:     requestdb.Answer,
	}
	marsh(request, 201, w)
}

//friends

func (cfg *apiConfig) get_friends(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	friendsdb, err := cfg.db.GetAllFriends(r.Context(), userid)
	if err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}

	sort.Slice(friendsdb, func(i, j int) bool {
		return strings.ToLower(friendsdb[i].SystemName) < strings.ToLower(friendsdb[j].SystemName)
	})

	friends := []User{}
	for _, frienddb := range friendsdb {

		friend := User{
			ID:         frienddb.ID,
			CreatedAt:  frienddb.CreatedAt,
			UpdatedAt:  frienddb.UpdatedAt,
			Email:      frienddb.Email,
			Avatar:     frienddb.Avatar,
			SystemName: frienddb.SystemName,
			Theme:      frienddb.Theme,
			Font:       frienddb.Font,
		}
		friends = append(friends, friend)
	}
	marsh(friends, 200, w)
}

func (cfg *apiConfig) delete_friend_request(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	friendIdString := r.PathValue("friendId")
	friendId, err := uuid.Parse(friendIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	friendsdb, err := cfg.db.GetIfFriends(r.Context(), database.GetIfFriendsParams{UserID: userid, FriendID: friendId})
	if err != nil {
		err_mes("Not friends", 400, w)
		return
	}

	err = cfg.db.DeleteFriendRequest(r.Context(), friendsdb.RequestID)
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}
	w.WriteHeader(204)
}

// //diaries

//diaries

func (cfg *apiConfig) post_diary(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	alterIdString := r.PathValue("alterId")
	alterId, err := uuid.Parse(alterIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}

	alterdb, err := cfg.db.GetAlter(r.Context(), alterId)
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}
	if userid != alterdb.UserID {
		err_mes("You're not the chosen one", 401, w)
		return
	}
	userdb, err := cfg.db.GetUser(r.Context(), alterdb.UserID)
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}

	block_colour_str := strings.Split(alterdb.Colour, " ")

	var block_colour []int
	for i := range 3 {
		c, err := strconv.Atoi(block_colour_str[i])
		if err != nil {
			err_mes("bad colour", 400, w)
			return
		}
		block_colour = append(block_colour, c)
	}
	less := ""
	more := ""
	for i := range 3 {
		moreplus := int(math.Round(float64(block_colour[i]) * (3)))

		if moreplus > 210 {
			more += "241 "
		} else if moreplus < 45 {
			more += "178 "
		} else {
			more += strconv.Itoa(int(float64(moreplus-45)*0.3+178)) + " "
		}

		less += strconv.Itoa(block_colour[i]/2) + " "

	}

	var bg_colour string
	var text_colour string
	if userdb.Theme%2 == 0 {
		bg_colour = more
		text_colour = less
	} else {
		bg_colour = less
		text_colour = more
	}
	diarydb, err := cfg.db.NewDiary(r.Context(), database.NewDiaryParams{AlterID: alterId, BgColour: bg_colour, BlockColour: alterdb.Colour, TextColour: text_colour, UserID: alterdb.UserID, Name: alterdb.Name})
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}
	diary := Diary{
		ID:          diarydb.ID,
		AlterId:     diarydb.AlterID,
		BgColour:    diarydb.BgColour,
		BgColour2:   diarydb.BgColour2,
		BlockColour: diarydb.BlockColour,
		TextColour:  diarydb.TextColour,
		Font:        diarydb.Font,
		UserId:      diarydb.UserID,
		Name:        diarydb.Name,
	}
	marsh(diary, 201, w)
}

func (cfg *apiConfig) get_diaries(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	diariesdb, err := cfg.db.GetDiaries(r.Context(), userid)
	if err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}

	sort.Slice(diariesdb, func(i, j int) bool {
		return strings.ToLower(diariesdb[i].Name) < strings.ToLower(diariesdb[j].Name)
	})

	diaries := []Diary{}
	for _, diarydb := range diariesdb {
		diary := Diary{
			ID:          diarydb.ID,
			AlterId:     diarydb.AlterID,
			BgColour:    diarydb.BgColour,
			BgColour2:   diarydb.BgColour2,
			BlockColour: diarydb.BlockColour,
			TextColour:  diarydb.TextColour,
			Font:        diarydb.Font,
			UserId:      diarydb.UserID,
			Name:        diarydb.Name,
			AlterAvatar: diarydb.Avatar,
		}
		diaries = append(diaries, diary)
	}
	marsh(diaries, 200, w)
}

func (cfg *apiConfig) get_diary(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	diaryIdString := r.PathValue("diaryId")
	diaryId, err := uuid.Parse(diaryIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	diarydb, err := cfg.db.GetDiary(r.Context(), database.GetDiaryParams{ID: diaryId, UserID: userid})
	if err != nil {
		err_mes("failed to get alter", 404, w)
		return
	}

	diary := Diary{
		ID:          diarydb.ID,
		AlterId:     diarydb.AlterID,
		BgColour:    diarydb.BgColour,
		BgColour2:   diarydb.BgColour2,
		BlockColour: diarydb.BlockColour,
		TextColour:  diarydb.TextColour,
		Font:        diarydb.Font,
		UserId:      diarydb.UserID,
		Name:        diarydb.Name,
	}
	marsh(diary, 200, w)
}

func (cfg *apiConfig) get_diary_by_alter(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	alterIdString := r.PathValue("alterId")
	alterId, err := uuid.Parse(alterIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	diarydb, err := cfg.db.GetDiaryByALter(r.Context(), database.GetDiaryByALterParams{AlterID: alterId, UserID: userid})
	if err != nil {
		err_mes("failed to get diary", 404, w)
		return
	}

	diary := Diary{
		ID:          diarydb.ID,
		AlterId:     diarydb.AlterID,
		BgColour:    diarydb.BgColour,
		BgColour2:   diarydb.BgColour2,
		BlockColour: diarydb.BlockColour,
		TextColour:  diarydb.TextColour,
		Font:        diarydb.Font,
		UserId:      diarydb.UserID,
		Name:        diarydb.Name,
	}
	marsh(diary, 200, w)
}

func (cfg *apiConfig) edit_diary(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	diaryIdString := r.PathValue("diaryId")
	diaryId, err := uuid.Parse(diaryIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	param := Diary{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong with decoding", 400, w)
		return
	}
	bgcolour := hex_to_rgb(param.BgColour, w)
	blockcolour := hex_to_rgb(param.BlockColour, w)
	textcolour := hex_to_rgb(param.TextColour, w)
	bgcolour2 := hex_to_rgb(param.BgColour2, w)

	diarydb, err := cfg.db.EditDiary(r.Context(), database.EditDiaryParams{BgColour: bgcolour, BgColour2: bgcolour2, BlockColour: blockcolour, TextColour: textcolour, Font: param.Font, ID: diaryId, UserID: userid})
	if err != nil {
		err_mes("Something went wrong with editing", 400, w)
		return
	}
	diary := Diary{
		ID:          diarydb.ID,
		AlterId:     diarydb.AlterID,
		BgColour:    diarydb.BgColour,
		BgColour2:   diarydb.BgColour2,
		BlockColour: diarydb.BlockColour,
		TextColour:  diarydb.TextColour,
		Font:        diarydb.Font,
		UserId:      diarydb.UserID,
		Name:        diarydb.Name,
	}
	marsh(diary, 201, w)
}

func (cfg *apiConfig) delete_diary(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	diaryIdString := r.PathValue("diaryId")
	diaryId, err := uuid.Parse(diaryIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}

	cfg.db.DeleteDiary(r.Context(), database.DeleteDiaryParams{ID: diaryId, UserID: userid})
	w.WriteHeader(204)
}

//diary_entries

func (cfg *apiConfig) post_diary_entry(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	diaryIdString := r.PathValue("diaryId")
	diaryId, err := uuid.Parse(diaryIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}

	param := DiaryEntry{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong with decoding", 400, w)
		return
	}

	diarydb, err := cfg.db.GetDiary(r.Context(), database.GetDiaryParams{ID: diaryId, UserID: userid})
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}
	if userid != diarydb.UserID {
		err_mes("You're not the chosen one", 401, w)
		return
	}
	if param.Name == "" {
		err_mes("Empty name", 409, w)
		return
	}
	if param.Text == "" {
		err_mes("Empty message", 400, w)
		return
	}

	entrydb, err := cfg.db.NewDiaryEntry(r.Context(), database.NewDiaryEntryParams{DiaryID: diaryId, Name: param.Name, Text: param.Text, UserID: userid})
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}
	entry := DiaryEntry{
		ID:      entrydb.ID,
		UserId:  entrydb.UserID,
		Name:    entrydb.Name,
		DiaryId: entrydb.DiaryID,
		Date:    entrydb.Date,
		Text:    entrydb.Text,
	}
	marsh(entry, 201, w)
}

func (cfg *apiConfig) get_diary_entries(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	diaryIdString := r.PathValue("diaryId")
	diaryId, err := uuid.Parse(diaryIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}

	entriesdb, err := cfg.db.GetDiaryEntriesByDiary(r.Context(), database.GetDiaryEntriesByDiaryParams{DiaryID: diaryId, UserID: userid})
	if err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}

	sort.Slice(entriesdb, func(i, j int) bool {
		return entriesdb[i].Date.After(entriesdb[j].Date)
	})

	entries := []DiaryEntry{}
	for _, entrydb := range entriesdb {
		if entrydb.Name == "StandardInvisibleEntryofkw" && entrydb.Text == "StandardInvisibleEntryofkw" {
			cfg.db.DeleteDiaryEntry(r.Context(), database.DeleteDiaryEntryParams{ID: entrydb.ID, UserID: userid})

			continue
		}
		entry := DiaryEntry{
			ID:      entrydb.ID,
			UserId:  entrydb.UserID,
			Name:    entrydb.Name,
			DiaryId: entrydb.DiaryID,
			Date:    entrydb.Date,
			Text:    entrydb.Text,
		}

		entries = append(entries, entry)
	}
	marsh(entries, 200, w)
}

func (cfg *apiConfig) get_diary_entry(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	entryIdString := r.PathValue("entryId")
	entryId, err := uuid.Parse(entryIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	entrydb, err := cfg.db.GetDiaryEntry(r.Context(), database.GetDiaryEntryParams{ID: entryId, UserID: userid})
	if err != nil {
		err_mes("failed to get alter", 404, w)
		return
	}

	entry := DiaryEntry{
		ID:      entrydb.ID,
		UserId:  entrydb.UserID,
		Name:    entrydb.Name,
		DiaryId: entrydb.DiaryID,
		Date:    entrydb.Date,
		Text:    entrydb.Text,
	}
	marsh(entry, 200, w)
}

func (cfg *apiConfig) edit_diary_entry(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	entryIdString := r.PathValue("entryId")
	entryId, err := uuid.Parse(entryIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	param := DiaryEntry{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong with decoding", 400, w)
		return
	}

	entrydb, err := cfg.db.EditDiaryEntry(r.Context(), database.EditDiaryEntryParams{Name: param.Name, Text: param.Text, ID: entryId, UserID: userid})
	if err != nil {
		err_mes("Something went wrong with editing", 400, w)
		return
	}
	entry := DiaryEntry{
		ID:      entrydb.ID,
		UserId:  entrydb.UserID,
		Name:    entrydb.Name,
		DiaryId: entrydb.DiaryID,
		Date:    entrydb.Date,
		Text:    entrydb.Text,
	}
	marsh(entry, 201, w)
}

func (cfg *apiConfig) delete_diary_entry(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	entryIdString := r.PathValue("entryId")
	entryId, err := uuid.Parse(entryIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}

	cfg.db.DeleteDiaryEntry(r.Context(), database.DeleteDiaryEntryParams{ID: entryId, UserID: userid})
	w.WriteHeader(204)
}

//diary files

func (cfg *apiConfig) post_diary_file(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	entryIdString := r.PathValue("entryId")
	entryId, err := uuid.Parse(entryIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	fileparam, header, err := r.FormFile("photo")

	if err != nil {
		err_mes("Something went wrong with file", 400, w)
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filepath := fmt.Sprintf("assets/files/%v", filename)
	outFile, err := os.Create("web/" + filepath)

	if err != nil {
		err_mes("Something went wrong with creating a file", 400, w)
		return
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, fileparam)
	if err != nil {
		err_mes("Something went wrong with copying a file", 400, w)
		return
	}

	_, err = cfg.db.GetDiaryEntry(r.Context(), database.GetDiaryEntryParams{ID: entryId, UserID: userid})
	if err != nil {
		// fmt.Println(err)
		err_mes("Something went wrong with db1", 400, w)
		return
	}

	filedb, err := cfg.db.NewDiaryEntryFile(r.Context(), database.NewDiaryEntryFileParams{EntryID: entryId, File: filepath, UserID: userid})
	if err != nil {
		// fmt.Println(err)
		err_mes("Something went wrong with db2", 400, w)
		return
	}
	file := DiaryEntryFile{
		ID:        filedb.ID,
		UserId:    filedb.UserID,
		EntryId:   filedb.EntryID,
		File:      filedb.File,
		CreatedAt: filedb.CreatedAt,
	}
	marsh(file, 201, w)
}

func (cfg *apiConfig) get_diary_files_by_entry(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	entryIdString := r.PathValue("entryId")
	entryId, err := uuid.Parse(entryIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}

	filesdb, err := cfg.db.GetDiaryEntryFilesByEntry(r.Context(), database.GetDiaryEntryFilesByEntryParams{EntryID: entryId, UserID: userid})
	if err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}

	sort.Slice(filesdb, func(i, j int) bool {
		return filesdb[i].CreatedAt.Before(filesdb[j].CreatedAt)
	})

	files := []DiaryEntryFile{}
	for _, filedb := range filesdb {
		file := DiaryEntryFile{
			ID:        filedb.ID,
			UserId:    filedb.UserID,
			EntryId:   filedb.EntryID,
			File:      filedb.File,
			CreatedAt: filedb.CreatedAt,
		}

		files = append(files, file)
	}
	marsh(files, 200, w)
}

func (cfg *apiConfig) get_diary_entry_file(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	fileIdString := r.PathValue("fileId")
	fileId, err := uuid.Parse(fileIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	filedb, err := cfg.db.GetDiaryEntryFile(r.Context(), database.GetDiaryEntryFileParams{ID: fileId, UserID: userid})
	if err != nil {
		err_mes("failed to get file", 404, w)
		return
	}
	file := DiaryEntryFile{
		ID:        filedb.ID,
		UserId:    filedb.UserID,
		EntryId:   filedb.EntryID,
		File:      filedb.File,
		CreatedAt: filedb.CreatedAt,
	}
	marsh(file, 200, w)
}

func (cfg *apiConfig) delete_diary_entry_file(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	fileIdString := r.PathValue("fileId")
	fileId, err := uuid.Parse(fileIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}
	filedb, err := cfg.db.GetDiaryEntryFile(r.Context(), database.GetDiaryEntryFileParams{ID: fileId, UserID: userid})
	if err != nil {
		err_mes("failed to get file", 404, w)
		return
	}
	if strings.Contains(filedb.File, "files") {
		err = os.Remove(filedb.File)
		if err != nil {
			err_mes("failed to delete the file", 404, w)
			return
		}
	}

	cfg.db.DeleteDiaryEntryFile(r.Context(), database.DeleteDiaryEntryFileParams{ID: fileId, UserID: userid})
	w.WriteHeader(204)
}

// avatars
//user

func (cfg *apiConfig) edit_user_avatar(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	usdb, err := cfg.db.GetUser(r.Context(), userId)
	if strings.Contains(usdb.Avatar, "avatars") {
		err = os.Remove("web/" + usdb.Avatar)
		if err != nil {
			fmt.Println(err)
			err_mes("failed to delete the file", 404, w)
			return
		}
	}

	fileparam, header, err := r.FormFile("photo")

	if err != nil {
		err_mes("Something went wrong with file", 400, w)
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filepath := fmt.Sprintf("assets/avatars/%v", filename)
	outFile, err := os.Create("web/" + filepath)

	if err != nil {
		err_mes("Something went wrong with creating a file", 400, w)
		return
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, fileparam)
	if err != nil {
		err_mes("Something went wrong with copying a file", 400, w)
		return
	}

	userdb, err := cfg.db.EditUserAvatar(r.Context(), database.EditUserAvatarParams{Avatar: filepath, ID: userId})
	user := User{
		ID:         userdb.ID,
		CreatedAt:  userdb.CreatedAt,
		UpdatedAt:  userdb.UpdatedAt,
		Email:      userdb.Email,
		Avatar:     userdb.Avatar,
		SystemName: userdb.SystemName,
		Theme:      userdb.Theme,
		Font:       userdb.Font,
	}
	marsh(user, 201, w)
}

func (cfg *apiConfig) delete_user_avatar(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	usdb, err := cfg.db.GetUser(r.Context(), userid)
	if strings.Contains(usdb.Avatar, "avatars") {
		err = os.Remove("web/" + usdb.Avatar)
		if err != nil {
			fmt.Println(err)
			err_mes("failed to delete the file", 404, w)
			return
		}
	}

	userdb, err := cfg.db.EditUserAvatar(r.Context(), database.EditUserAvatarParams{Avatar: "assets/default_avater.jpg", ID: userid})
	user := User{
		ID:         userdb.ID,
		CreatedAt:  userdb.CreatedAt,
		UpdatedAt:  userdb.UpdatedAt,
		Email:      userdb.Email,
		Avatar:     userdb.Avatar,
		SystemName: userdb.SystemName,
		Theme:      userdb.Theme,
		Font:       userdb.Font,
	}
	marsh(user, 201, w)
}

//alter

func (cfg *apiConfig) edit_alter_avatar(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	alterIdString := r.PathValue("alterId")
	alterId, err := uuid.Parse(alterIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}

	aldb, err := cfg.db.GetAlter(r.Context(), alterId)
	if strings.Contains(aldb.Avatar, "avatars") {
		err = os.Remove("web/" + aldb.Avatar)
		if err != nil {
			fmt.Println(err)
			err_mes("failed to delete the file", 404, w)
			return
		}
	}

	fileparam, header, err := r.FormFile("photo")

	if err != nil {
		err_mes("Something went wrong with file", 400, w)
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filepath := fmt.Sprintf("assets/avatars/%v", filename)
	outFile, err := os.Create("web/" + filepath)

	if err != nil {
		err_mes("Something went wrong with creating a file", 400, w)
		return
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, fileparam)
	if err != nil {
		err_mes("Something went wrong with copying a file", 400, w)
		return
	}

	alterdb, err := cfg.db.EditAlterAvatar(r.Context(), database.EditAlterAvatarParams{Avatar: filepath, ID: alterId, UserID: userId})
	alter := Alter{
		ID:          alterdb.ID,
		CreatedAt:   alterdb.CreatedAt,
		UpdatedAt:   alterdb.UpdatedAt,
		Name:        alterdb.Name,
		Avatar:      alterdb.Avatar,
		Pronouns:    alterdb.Pronouns,
		Age:         alterdb.Age,
		Role:        alterdb.AlterRole,
		Description: alterdb.Description,
		Colour:      alterdb.Colour,
		Fronting:    alterdb.Fronting,
		User_id:     alterdb.UserID,
	}
	marsh(alter, 201, w)

}

func (cfg *apiConfig) delete_alter_avatar(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}
	alterIdString := r.PathValue("alterId")
	alterId, err := uuid.Parse(alterIdString)
	if err != nil {
		err_mes("bad uuid", 400, w)
		return
	}

	aldb, err := cfg.db.GetAlter(r.Context(), alterId)
	if strings.Contains(aldb.Avatar, "avatars") {
		err = os.Remove("web/" + aldb.Avatar)
		if err != nil {
			fmt.Println(err)
			err_mes("failed to delete the file", 404, w)
			return
		}
	}

	alterdb, err := cfg.db.EditAlterAvatar(r.Context(), database.EditAlterAvatarParams{Avatar: "assets/default_avater.jpg", ID: alterId, UserID: userId})
	alter := Alter{
		ID:          alterdb.ID,
		CreatedAt:   alterdb.CreatedAt,
		UpdatedAt:   alterdb.UpdatedAt,
		Name:        alterdb.Name,
		Avatar:      alterdb.Avatar,
		Pronouns:    alterdb.Pronouns,
		Age:         alterdb.Age,
		Role:        alterdb.AlterRole,
		Description: alterdb.Description,
		Colour:      alterdb.Colour,
		Fronting:    alterdb.Fronting,
		User_id:     alterdb.UserID,
	}
	marsh(alter, 201, w)

}

//for newbies

func (cfg *apiConfig) post_for_newbies(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	newdb, err := cfg.db.NewForNewbies(r.Context(), database.NewForNewbiesParams{UserID: userid, Text: default_for_newbies})
	if err != nil {
		err_mes("Something went wrong with db", 400, w)
		return
	}
	new := ForNewbies{
		ID:     newdb.ID,
		UserId: newdb.UserID,
		Text:   newdb.Text,
	}
	marsh(new, 201, w)
}

func (cfg *apiConfig) get_for_newbies(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	newdb, err := cfg.db.GetForNewbies(r.Context(), userid)
	if err != nil {
		err_mes("failed to get alter", 404, w)
		return
	}

	new := ForNewbies{
		ID:     newdb.ID,
		UserId: newdb.UserID,
		Text:   newdb.Text,
	}
	marsh(new, 200, w)
}

func (cfg *apiConfig) edit_for_newbies(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	userid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		err_mes("Invalid jwt", 401, w)
		return
	}

	param := ForNewbies{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong with decoding", 400, w)
		return
	}

	newdb, err := cfg.db.EditForNewbies(r.Context(), database.EditForNewbiesParams{Text: param.Text, UserID: userid})
	if err != nil {
		err_mes("Something went wrong with editing", 400, w)
		return
	}
	new := ForNewbies{
		ID:     newdb.ID,
		UserId: newdb.UserID,
		Text:   newdb.Text,
	}
	marsh(new, 201, w)
}

//security

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	param := UserPass{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		err_mes("Something went wrong", 400, w)
		return
	}
	if param.Email == "" {
		err_mes("Empty email", 400, w)
		return
	}
	if param.Password == "" {
		err_mes("Empty password", 400, w)
		return
	}
	userdb, err := cfg.db.GetUserByEmail(r.Context(), param.Email)

	if err != nil {
		err_mes(fmt.Sprintf("%v", err), 400, w)
		return
	}

	valid_password, err := auth.CheckPasswordHash(param.Password, userdb.HashedPassword)

	if err != nil {
		err_mes("error email or password", 400, w)
		return
	}
	if !valid_password {
		err_mes("Incorrect email or password", 401, w)
		return
	}

	tok, err := auth.MakeJWT(userdb.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		err_mes("Something went wrong with the jwt token", 400, w)
		return
	}
	refresh_str := auth.MakeRefreshToken()
	refresh, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: refresh_str, UserID: userdb.ID})
	if err != nil {
		err_mes("couldn't add refresh token", 400, w)
		return
	}
	user := User{
		ID:         userdb.ID,
		CreatedAt:  userdb.CreatedAt,
		UpdatedAt:  userdb.UpdatedAt,
		Email:      userdb.Email,
		Token:      tok,
		Refresh:    refresh.Token,
		Avatar:     userdb.Avatar,
		SystemName: userdb.SystemName,
		Theme:      userdb.Theme,
		Font:       userdb.Font,
	}
	marsh(user, 200, w)
}

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	// if err := cfg.dbc.Ping(); err != nil {
	// 	fmt.Printf("Database is not reachable: %v", err)
	// }
	// fmt.Println("Database connection pool configured and ping successful.")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No authorization", 400, w)
		return
	}
	dbtoken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	if dbtoken.RevokedAt.Valid {
		err_mes("Token revoked", 401, w)
		return
	}
	if dbtoken.ExpiresAt.Before(time.Now()) {
		err_mes("Token expired", 401, w)
		return
	}

	refreshed, err := auth.MakeJWT(dbtoken.UserID, cfg.JWTSecret, time.Hour)
	if err != nil {
		err_mes("Something went wrong with the jwt token", 400, w)
		return
	}

	type RefreshedToken struct {
		// User_id uuid.UUID `json:"user_id"`
		Token string `json:"token"`
	}
	rt := RefreshedToken{Token: refreshed}
	marsh(rt, 200, w)
}

func (cfg *apiConfig) revoke_token(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_mes("No authorization", 400, w)
		return
	}
	dbtoken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		err_mes("No token", 401, w)
		return
	}
	cfg.db.RevokeRefreshToken(r.Context(), dbtoken.Token)
	w.WriteHeader(204)
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	log.Fatalf("%v", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("DB not opening")
	}
	db.SetMaxOpenConns(25)                 // максимум открытых соединений
	db.SetMaxIdleConns(5)                  // максимум простаивающих соединений
	db.SetConnMaxLifetime(5 * time.Minute) // время жизни соединения
	db.SetConnMaxIdleTime(1 * time.Minute) // время простоя соединения

	// Проверка подключения
	if err = db.Ping(); err != nil {
		fmt.Printf("Failed to ping DB: %v", err)
	}
	fmt.Println("everything ok")
	dbQueries := database.New(db)
	apiCfg := apiConfig{}
	apiCfg.platform = os.Getenv("PLATFORM")
	apiCfg.db = dbQueries
	apiCfg.dbc = db
	jwtSecret := os.Getenv("JWT_SECRET")
	apiCfg.JWTSecret = jwtSecret

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	//app
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("/web")))))

	//readiness
	mux.HandleFunc("GET /api/healthz", Readiness)

	//security
	mux.HandleFunc("POST /api/login", apiCfg.login)
	mux.HandleFunc("POST /api/refresh", apiCfg.refresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.revoke_token)

	//users
	mux.HandleFunc("POST /api/users", apiCfg.post_user)
	mux.HandleFunc("GET /api/users/{userId}", apiCfg.get_user)
	mux.HandleFunc("PUT /api/users", apiCfg.edit_userinfo)
	mux.HandleFunc("PUT /api/userinfo", apiCfg.edit_userinfo_settings)
	//avatar
	mux.HandleFunc("PUT /api/users/avatar", apiCfg.edit_user_avatar)
	mux.HandleFunc("DELETE /api/users/avatar", apiCfg.delete_user_avatar)

	//alters
	mux.HandleFunc("POST /api/alters", apiCfg.post_alter)
	mux.HandleFunc("GET /api/alters", apiCfg.get_alters)
	mux.HandleFunc("GET /api/alters_wo_diaries", apiCfg.get_alters_without_diary)
	mux.HandleFunc("GET /api/alters/{alterId}", apiCfg.get_alter)
	mux.HandleFunc("PUT /api/alters/{alterId}", apiCfg.edit_alter)
	mux.HandleFunc("DELETE /api/alters/{alterId}", apiCfg.delete_alter)
	//avatar
	mux.HandleFunc("PUT /api/alters/avatar/{alterId}", apiCfg.edit_alter_avatar)
	mux.HandleFunc("DELETE /api/alters/avatar/{alterId}", apiCfg.delete_alter_avatar)

	//fronts
	mux.HandleFunc("PATCH /api/alters/{alterId}", apiCfg.change_front)
	mux.HandleFunc("GET /api/alters/fronting/{userId}", apiCfg.get_fronting_alters)
	mux.HandleFunc("GET /api/fronts/time", apiCfg.get_fronts_by_time)

	//friends

	//requests

	mux.HandleFunc("POST /api/friends/requests/{receiverId}", apiCfg.post_friend_request)
	mux.HandleFunc("GET /api/friends/requests/sender", apiCfg.get_friend_requests_as_sender)
	mux.HandleFunc("GET /api/friends/requests/reciever", apiCfg.get_friend_requests_as_receiver)
	mux.HandleFunc("PUT /api/friends/requests", apiCfg.put_answer_friend_request)

	//friends

	mux.HandleFunc("GET /api/friends", apiCfg.get_friends)
	mux.HandleFunc("DELETE /api/friends/{friendId}", apiCfg.delete_friend_request)

	//diaries

	//diaries
	mux.HandleFunc("POST /api/diaries/{alterId}", apiCfg.post_diary)
	mux.HandleFunc("GET /api/diaries", apiCfg.get_diaries)
	mux.HandleFunc("GET /api/diary/{alterId}", apiCfg.get_diary_by_alter)
	mux.HandleFunc("GET /api/diaries/{diaryId}", apiCfg.get_diary)
	mux.HandleFunc("PUT /api/diaries/{diaryId}", apiCfg.edit_diary)
	mux.HandleFunc("DELETE /api/diaries/{diaryId}", apiCfg.delete_diary)

	//entries

	mux.HandleFunc("POST /api/diaries/entries/{diaryId}", apiCfg.post_diary_entry)
	mux.HandleFunc("GET /api/diaries/entries/{diaryId}", apiCfg.get_diary_entries)
	mux.HandleFunc("GET /api/diaries/entries/{diaryId}/{entryId}", apiCfg.get_diary_entry)
	mux.HandleFunc("PUT /api/diaries/entries/{diaryId}/{entryId}", apiCfg.edit_diary_entry)
	mux.HandleFunc("DELETE /api/diaries/entries/{diaryId}/{entryId}", apiCfg.delete_diary_entry)

	//files

	mux.HandleFunc("POST /api/diaries/files/{entryId}", apiCfg.post_diary_file)
	mux.HandleFunc("GET /api/diaries/files/{entryId}", apiCfg.get_diary_files_by_entry)
	mux.HandleFunc("GET /api/diaries/files/{entryId}/{fileId}", apiCfg.get_diary_entry_file)
	mux.HandleFunc("DELETE /api/diaries/files/{fileId}", apiCfg.delete_diary_entry_file)

	//for newbies

	mux.HandleFunc("POST /api/for_newbies", apiCfg.post_for_newbies)
	mux.HandleFunc("GET /api/for_newbies", apiCfg.get_for_newbies)
	mux.HandleFunc("PUT /api/for_newbies", apiCfg.edit_for_newbies)

	srv.ListenAndServe()
}
