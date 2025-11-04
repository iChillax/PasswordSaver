package models

type BasicAuth struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Username    string `bson:"username" json:"username"`
	Password    string `bson:"password" json:"password"`
}
