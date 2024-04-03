package model

import (
	"time"
)

type User struct {
	ID string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
	Email string `json:"email" firestore:"email"`
	Sub string `json:"sub" firestore:"sub"`
	AppRole string `json:"app_role" firestore:"app_role"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

type Workspace struct {
}

type Source struct {
	ID string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"` 
}

type Destination struct {
}
