package model

import (
	"time"
)

// User
type User struct {
	ID string `json:"id" firestore:"id"` // "" for unidentified
	Name string `json:"name" firestore:"name"`
	Email string `json:"email" firestore:"email"`
	Sub string `json:"sub" firestore:"sub"`
	AppRole string `json:"app_role" firestore:"app_role"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

func NewUser(name string, email string, sub string, appRole string) User {
	// Creaates new user without an identity (attributed at datadabase insertion)
	return User{ "", name, email, sub, appRole, time.Now(), time.Now() }
}

func (u User) HasIdentity() bool {
	return u.ID != ""
}

type Workspace struct {
}

type Source struct {
	ID string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"` 
}

type Destination struct {
}
