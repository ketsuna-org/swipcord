package controllers

import (
	"api/models"
	"api/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func DiscordOauth2(w http.ResponseWriter, r *http.Request) {
	// let's get the envs variables
	clientID := utils.GetEnv("DISCORD_CLIENT_ID", "")
	// format the URL
	urlFormatted := fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=http://localhost:4000/discord/callback&response_type=code&scope=identify&prompt=none", clientID)

	http.Redirect(w, r, urlFormatted, http.StatusSeeOther)
}

type DiscordResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func DiscordCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "No code provided")
		return
	}
	// let's get the envs variables
	clientID := utils.GetEnv("DISCORD_CLIENT_ID", "")
	clientSecret := utils.GetEnv("DISCORD_CLIENT_SECRET", "")

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", "http://localhost:4000/discord/callback")

	req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	resp, err := client.Do(req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()
	// let's print the body of the response

	discordResponse := DiscordResponse{}
	err = json.NewDecoder(resp.Body).Decode(&discordResponse)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// let's do a new request to get the user information
	req, err = http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", discordResponse.AccessToken))
	resp, err = client.Do(req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	discordUser := models.DiscordUser{}
	err = json.NewDecoder(resp.Body).Decode(&discordUser)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// retrieve database from the context
	db := models.GetDatabase(r)
	user := models.User{}
	// let's modify expiration time.
	// we will add the current time to the expiration time
	expireTime := time.Now().Add(time.Duration(discordResponse.ExpiresIn) * time.Second)
	// check if the user already exists
	db.Where("discord_id = ?", discordUser.ID).First(&user)
	if user.ID == 0 {
		// create the user
		user = models.User{
			DiscordID:       discordUser.ID,
			Coins:           0,
			DiscordIdentity: &discordUser,
			AccessToken:     discordResponse.AccessToken,
			RefreshToken:    discordResponse.RefreshToken,
			ExpiresIn:       expireTime,
		}
		db.Create(&user)
	} else {
		// update the user
		user.DiscordIdentity = &discordUser
		user.AccessToken = discordResponse.AccessToken
		user.RefreshToken = discordResponse.RefreshToken
		user.ExpiresIn = expireTime
		db.Save(&user)
	}

	utils.RespondWithJSON(w, map[string]string{"message": "Successfully authenticated"}, http.StatusOK)
}
