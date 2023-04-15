package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

//Authenticate will authenticate the user with the incoming credentials
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request){
	//recover the body
	var reqPayload struct{
		Email string `json:"email"`
		Passoword string `json:"password"`
	}

	err:= app.readJSON(w,r,&reqPayload)
	if err!= nil{
		app.errorJSON(w, errors.New("invalid body sent"),http.StatusBadRequest)
		log.Println("invalid body error")
		return
	}

	// check if the user exist
	user, err:= app.Models.User.GetByEmail(reqPayload.Email);
	if err!= nil{
		app.errorJSON(w, errors.New("invalid credentials"),http.StatusBadRequest)
		log.Println("get by email error")
		return
	}

	// match the password
	valid,err:=user.PasswordMatches(reqPayload.Passoword)

	if err!= nil || !valid{
		app.errorJSON(w, errors.New("invalid credentials"),http.StatusBadRequest)
		log.Println("password mismatch error")
		return
	}
	// log it to logger

	err= app.logRequest(w,"Authentication", fmt.Sprintf("user %s is authenticate",user.Email))
	if err!= nil {
		app.errorJSON(w, errors.New("unable to log to logger"),http.StatusBadRequest)
		return
	}

	// send back the response

	response:= jsonResponse{
		Error: false,
		Message:fmt.Sprintf("Logged in User with email %s",user.Email),
		Data: user,
	}

	_=app.writeJSON(w,http.StatusAccepted,response)
}

func  (app *Config) logRequest(w http.ResponseWriter,name string, data string ) error{
	var entry struct{
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name=name
	entry.Data=data

	jsonData,_:=json.MarshalIndent(entry,"","\t")

	request,err:=http.NewRequest("POST","http://logger-service/log",bytes.NewBuffer(jsonData))
	if err!= nil{
		return err
	}

	client:=&http.Client{}

	_,err= client.Do(request)
	if err!= nil{
		return err
	}

	return nil


}