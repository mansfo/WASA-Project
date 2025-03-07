package api

import (
	"time"
	"wasaphoto-2009711/service/database"
)

// The user identifier for APIs
type Userid struct {
	ID string `json:"user_id"` // The string representing the identifier
}

// The username of the user for APIs
type Username struct {
	NameUser string `json:"username"`
}

// A struct for APIs representing a user
type User struct {
	IdUser Userid   `json:"user_id"`  // The identifier of the user
	Name   Username `json:"username"` // The username of the user
}

// A struct for APIs representing a user and the path to his/her profile picture
type UserSearched struct {
	Usr User   `json:"user"`          // The user
	Img string `json:"path_to_image"` // The path of the profile picture
}

// A struct for APIs representing the image in binary format and its identifier
type Image struct {
	ImageId  ImageIdentifier `json:"image_id"`      // The identifier of the image
	ImageUri string          `json:"path_to_image"` // The URI of the image
}

// A struct for APIs representing the image's unique identifier
type ImageIdentifier struct {
	IdImage int64 `json:"image_id"` // The identifier of the image
}

// The identifier for a photo for APIs
type PhotoId struct {
	IdPhoto int `json:"photo_id"`
}

// The identifier for a comment for APIs
type CommentId struct {
	IdComment int64 `json:"comment_id"`
}

// A struct representing the body sent by a user that publishes a comment
type CommentBody struct {
	Comment string `json:"comment"`
}

// A struct for APIs representing a comment on a photo
type Comment struct {
	IdComment      CommentId `json:"comment_id"` // The identifier of the comment
	IdUser         Userid    `json:"user_id"`    // The identifier of the user who wrote the comment
	Photo          PhotoId   `json:"photo_id"`   // The identifier of the photo
	Saying         string    `json:"comment"`    // The text that the user has written as a comment
	PublishingDate time.Time `json:"date"`       // The date in which the comment was published
}

// A struct for APIs representing a photo published by someone
type Photo struct {
	IdPhoto        PhotoId             `json:"photo_id"`        // The identifier of the photo
	IdUser         Userid              `json:"user_id"`         // The identifier of the user
	Picture        ImageIdentifier     `json:"image"`           // The image in binary format and its identifier
	UpdateDate     time.Time           `json:"date"`            // The date in which the image was published
	Likes          []database.Username `json:"likes"`           // The username of the users that like the photo
	LikeCounter    int                 `json:"like_counter"`    // How many users like the photo
	Comments       []database.Comment  `json:"comments"`        // The comments published on the photo
	CommentCounter int                 `json:"comment_counter"` // How many comments were published on the photo
}

// A struct for APIs representing the profile of a user
type UserProfile struct {
	Owner            User              `json:"user"`              // The user that owns the profile
	ProfilePict      Image             `json:"profile_picture"`   // An image that the user chose to identify him/her
	Photos           []database.Photo  `json:"photos"`            // The photos published by the user
	PhotoCounter     int               `json:"photo_counter"`     // How many photos the user has published
	Followers        []database.Userid `json:"followers"`         // The users that follow this profile
	FollowersCounter int               `json:"followers_counter"` // How many users follow this profile
	Following        []database.Userid `json:"following"`         // The users followed by this profile
	FollowingCounter int               `json:"following_counter"` // How many users are followed by this profile
}

// Convertion of a Userid from the API package to a Userid in the database package
func (u Userid) toDatabase() database.Userid {
	return database.Userid{
		IdUser: u.ID,
	}
}

// Conversion of an image identifier from the API package to an image identifier in the database package
func (ii ImageIdentifier) toDatabase() database.ImageId {
	return database.ImageId{
		IdImage: ii.IdImage,
	}
}

// Convertion of a PhotoId from the API package to a PhotoId in the database package
func (p PhotoId) toDatabase() database.PhotoId {
	return database.PhotoId{
		IdPhoto: p.IdPhoto,
	}
}

// Convertion of a CommentId from the API package to a CommentId in the database package
func (c CommentId) toDatabase() database.CommentId {
	return database.CommentId{
		IdComment: c.IdComment,
	}
}
