package database

import "time"

// Userid struct for the db, represents the unique identifer for a user
type Userid struct {
	IdUser string `json:"user_id"`
}

// Username struct for the db, represents the username of a user
type Username struct {
	Name string `json:"username"`
}

// User struct for the db, represents the userid and the username of a user
type User struct {
	IdUser Userid   `json:"user_id"`
	Name   Username `json:"username"`
}

// A struct for the db representing a user and the path to his/her profile picture
type UserSearched struct {
	Usr User   `json:"user"`          // The user
	Img string `json:"path_to_image"` // The path of the profile picture
}

// UserProfile struct for the db, represeting the user that owns the profile, his/her profile picture,
// his/her photos, his/her followers and followed
type UserProfile struct {
	Owner            User           `json:"user"`              // The user that owns the profile
	ProfilePict      Image          `json:"profile_picture"`   // The profile picture of this profile
	Photos           []Photo        `json:"photos"`            // The photos published by the user
	PhotoCounter     int            `json:"photo_counter"`     // How many photos there are on the profile
	Followers        []UserSearched `json:"followers"`         // The userid of the users that follow this profile
	FollowersCounter int            `json:"followers_counter"` // How many users are following this profile
	Following        []UserSearched `json:"following"`         // The userid of the users followed by this profile
	FollowingCounter int            `json:"following_counter"` // How many users are followed by this profile
}

// PhotoId struct for the db, represents the unique identifier for a photo
type PhotoId struct {
	IdPhoto int `json:"photo_id"`
}

// Photo struct for the db, represents the identifier of the photo, the identifier of the user that published it,
// the image published, the date in which it was published, an array for the likes and an array for the comments
type Photo struct {
	IdPhoto        PhotoId        `json:"photo_id"`        // The identifier of the photo
	User           UserSearched   `json:"user"`            // The identifier, the username and the propic of the user that has published the photo
	ImgPath        string         `json:"image"`           // The image published
	Date           time.Time      `json:"date"`            // The publication date of the photo
	Likes          []UserSearched `json:"likes"`           // The username, userid and proPic of the users that like the photo
	LikeCounter    int            `json:"like_counter"`    // How many users like the photo
	Comments       []Comment      `json:"comments"`        // The comments on the photo
	CommentCounter int            `json:"comment_counter"` // How many comments were published on the photo
	IsLiked        bool           `json:"like_status"`     // It's equal true if the requesting user likes this photo
}

// CommentId struct for the db, represents the unique identifier for a comment
type CommentId struct {
	IdComment int64 `json:"comment_id"`
}

// Comment struct for the db, represents the identifier of the comment, the identifier of the photo where the comment was published,
// the identifier of the user that published the comment, the text published and the publication date of the comment
type Comment struct {
	IdComment CommentId    `json:"comment_id"` // The identifier of the comment
	IdPhoto   PhotoId      `json:"photo_id"`   // The identifier of the photo
	User      UserSearched `json:"user"`       // The identifier, the username and the propic of the user
	Text      string       `json:"comment"`    // The text published as a comment
	Date      time.Time    `json:"date"`       // The publication date of the comment
}

// ImageId struct for the db, represents the unique identifier of the image
type ImageId struct {
	IdImage int64 `json:"image_id"` // The unique identifier of the image
}

// Image struct for the db, represents the identifier of the image and the image in binary format
type Image struct {
	IdImage  ImageId `json:"image_id"`  // The unique identifier of the image
	UriImage string  `json:"uri_image"` // The URI of the image
}
