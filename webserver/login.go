// Copyright 2016 Graeme Dykes.

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type LoginTemplateParams struct {
	ErrorMessage string
	EmailAddress string
}

func loginHandlerFavicon(w http.ResponseWriter, r *http.Request) {

}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("loginHandler")

	switch r.Method {

	case "GET":

		sessionToken, err := r.Cookie("sessionTokenStaff")
		sessionIsValid, _, _ := sessionList.CheckSessionToken(sessionToken, err)
		if sessionIsValid {
			sessionList.invalidateSessionToken(sessionToken.Value)
		}

		returnLoginHtml(w, r, "", "")
		return

	case "POST":

		// userName := r.FormValue("Email")
		// password := r.FormValue("Password")

		/*
			userLoginMessage := dynWebComms.Message_UserLogin{UserName: userName, Password: password}

			responseBytes, err := messageQueue.SendMessageSynch(dynWebComms.EncodeToJsonBinary(userLoginMessage), dynWebComms.WEB_TO_WEBDBSRVR_STAFF_LOGIN, 20)

			if err != nil {
				returnLoginHtml(w, r, "Error during authentication.", userName)
				// http.Redirect(w, r, "/login/", http.StatusFound)
				return
			}

			var userLoginResponseMessage dynWebComms.Message_UserLoginResponse

			errJsonDecode := json.Unmarshal(responseBytes, &userLoginResponseMessage)
			if errJsonDecode != nil {
				returnLoginHtml(w, r, "Error during authentication.", userName)
				// http.Redirect(w, r, "/login/", http.StatusFound)
				return
			}

			if !userLoginResponseMessage.LogInSuccess {
				returnLoginHtml(w, r, "Invalid email address or password.", userName)
				// http.Redirect(w, r, "/login/", http.StatusFound)
				return
			}

			sessionList.successfulLogin(userLoginResponseMessage.UserId, userLoginResponseMessage.UserName, userLoginResponseMessage.SessionToken)

			// This is a session cookie. Expires when you close the browser.
			cookie := http.Cookie{Name: "sessionTokenStaff", Value: userLoginResponseMessage.SessionToken, Path: "/"}
			http.SetCookie(w, &cookie)

			http.Redirect(w, r, "/", http.StatusFound)
		*/
	}
}

func SafeHTML(args ...interface{}) template.HTML {

	// fmt.Println("SafeHTML")

	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		// s = fmt.Sprint(args...)
		fmt.Println("Error in SafeHTML")
		return ""
	}
	// fmt.Println(s)
	return template.HTML(s)
}

func returnLoginHtml(w http.ResponseWriter, r *http.Request, errorMessage string, emailAddress string) {
	w.Header().Set("Content-Type", "text/html")

	// Previously called this
	// template.ParseFiles()
	// However that always stripped out comments, and was incompatible with FuncMap

	// This should not be called everytime. Move to main.

	t := template.New("Login template")

	t = t.Funcs(template.FuncMap{"safeHTML": SafeHTML})

	bytes, readErr := ioutil.ReadFile("Login.htmx")
	if readErr != nil {
		fmt.Println("Error in returnLoginHtml -", readErr)
		return
	}
	templateString := string(bytes)

	t, err := t.Parse(templateString)
	if err != nil {
		fmt.Println("Error in returnLoginHtml -", err)
		return
	}

	t.Execute(w, &LoginTemplateParams{ErrorMessage: errorMessage, EmailAddress: emailAddress})
}
