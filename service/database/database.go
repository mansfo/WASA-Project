/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// Error messages
var ErrSelfLike error = errors.New("User can't like his own photos")

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// Creates a new user in the database, returns the username and the userid of the user, true if the user is new and false if it's not new and also error
	NewUser(u string) (User, bool, error)

	// Check if the user exists
	CheckUser(u Userid) (bool, error)

	// Returns the username, the path to the user's profile picture and an error
	GetUserInfos(u Userid) (Username, string, error)

	// Check if a username exists in the database and if it is taken by someone else, regardless of letter case: if so, returns true
	CheckIfUsernameIsTaken(id Userid, username string) (bool, error)

	// User u sets string s as his/her new username, it returns the username and an error
	SetMyUserName(u Userid, un string) error

	// User u sets his/her profile picture, returns an error
	SetMyProfilePict(u Userid, newImage ImageId) error

	// Returns the profile of a the user v and an error
	GetProfile(u Userid, v Userid) (UserProfile, error)

	// Returns the list of the users that match searchedName and an error
	SearchUser(searcher Userid, searchedName Userid) ([]UserSearched, error)

	// Returns the stream of the user and an error
	GetMyStream(u Userid) ([]Photo, error)

	// Returns the infos of the users that like the photo and an error
	GetLikeOnPhoto(u Userid, p PhotoId) (int, []UserSearched, error)

	// User u puts a like on photo p, returns an error
	LikePhoto(u Userid, p PhotoId) error

	// returns how many likes there are on the photo and an error
	CountLikes(u Userid, p PhotoId) (int, error)

	// Check if user u likes the photo p, returns 1 if so, 0 if not
	CheckLike(u Userid, p PhotoId) (int, error)

	// User u removes his/her like from photo p, returns an error
	UnlikePhoto(u Userid, p PhotoId) error

	// Check if user u can like the photo, i. e. check if the photo wasn't uploaded by user u
	SelfLike(u Userid, p PhotoId) (bool, error)

	// Returns the comments on the photo and an error
	GetComments(u Userid, p PhotoId) (int, []Comment, error)

	// Returns how many comments are published on the photo and an error
	CountComments(u Userid, p PhotoId) (int, error)

	// User u posts a string s as a comment on photo p, returns an error
	CommentPhoto(u Userid, p PhotoId, s string, dt string) (Comment, error)

	// Remove a comment, returns an error
	UncommentPhoto(c CommentId) error

	// Check if user u is the user who has written  the comment c
	IsCommentAuthor(u Userid, c CommentId) (bool, error)

	// Creates a new image in the database, returns the identifier of the image and an error
	NewImage(dt string) (int64, error)

	// Sets the new image's path in the db
	SetImagePath(imid ImageId, path string) error

	// Returns the photos published by user v (if v has not banned u) and an error
	GetUsersPhotos(u Userid, v Userid) (int, []Photo, error)

	// Returns how many photos were published by user v
	CountPhotos(u Userid, v Userid) (int, error)

	// User u uploads an image i, returns the id of the photo and an error
	UploadPhoto(u Userid, i ImageId, d string) (int64, error)

	// User u removes photo p, returns an error
	RemovePhoto(p PhotoId) error

	// Returns an array containing the userid of all the user that follow v (minus the users banned by u and the users that banned u) and an error
	GetFollowers(u Userid, v Userid) (int, []UserSearched, error)

	// Returns how many users are following v minus the users banned by u and the users that banned u
	CountFollowers(u Userid, v Userid) (int, error)

	// Returns an array containing the userid of all the user followed by v (minus the users banned by u and the users that banned u) and an error
	GetFollowing(u Userid, v Userid) (int, []UserSearched, error)

	// Returns how many users are followed by v minus the users banned by u and the users that banned u
	CountFollowing(u Userid, v Userid) (int, error)

	// User u follows user v (if v has not banned u), returns an error
	FollowUser(u Userid, v Userid) error

	// User u follows user v, returns an error
	UnfollowUser(u Userid, v Userid) error

	// Check if v is followed by u
	IsFollowed(u Userid, v Userid) (bool, error)

	// Returns an array containing the userid of all the user banned by u and an error
	GetBanned(u Userid) ([]UserSearched, error)

	// User u bans user v, returns an error
	BanUser(u Userid, v Userid) error

	// User u unbans user v, returns an error
	UnbanUser(u Userid, v Userid) error

	// Check if user u is banned by user v, returns also an error
	CheckBan(u Userid, v Userid) (bool, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// Creates the sql tables for the app
func createDatabase(db *sql.DB) error {
	tables := [7]string{
		`CREATE TABLE IF NOT EXISTS images(
			imageid INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			uploadtime DATETIME,
			pathimage VARCHAR(99999)
		);`,
		`CREATE TABLE IF NOT EXISTS users(
			userid VARCHAR(16) NOT NULL PRIMARY KEY,
			username VARCHAR(16) NOT NULL UNIQUE,
			profilePicture INTEGER,
			FOREIGN KEY (profilePicture) references images(imageid) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS photos(
			photoid INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			user VARCHAR(16) NOT NULL,
			image INTEGER NOT NULL,
			date DATETIME NOT NULL,
			FOREIGN KEY (user) references users(userid) ON DELETE CASCADE,
			FOREIGN KEY (image) references images(imageid) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS likes(
			user VARCHAR(16) NOT NULL,
			photo INTEGER NOT NULL,
			PRIMARY KEY (user, photo),
			FOREIGN KEY (user) references users(userid) ON DELETE CASCADE,
			FOREIGN KEY (photo) references photos(photoid) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS comments(
			commentid INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			photo INTEGER NOT NULL,
			user VARCHAR(16),
			text VARCHAR(300) NOT NULL,
			date DATETIME NOT NULL,
			FOREIGN KEY (user) references users(userid) ON DELETE CASCADE,
			FOREIGN KEY (photo) references photos(photoid) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS followers(
			follower VARCHAR(16) NOT NULL,
			followed VARCHAR(16) NOT NULL,
			PRIMARY KEY (follower, followed),
			FOREIGN KEY (follower) references users(userid) ON DELETE CASCADE,
			FOREIGN KEY (followed) references users(userid) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS bans(
			banner VARCHAR(16) NOT NULL,
			banned VARCHAR(16) NOT NULL,
			PRIMARY KEY (banner, banned),
			FOREIGN KEY (banner) references users(userid) ON DELETE CASCADE,
			FOREIGN KEY (banned) references users(userid) ON DELETE CASCADE
		);`,
	}

	// Loop to create the tables
	for i := 0; i < len(tables); i++ {
		sqlStmt := tables[i]
		_, err := db.Exec(sqlStmt)
		if err != nil {
			return err
		}
	}
	// Adding a path to an image in the db to use it as a default profile picture
	img_path := "../assets/images/noPicture.jpg"
	_, err := db.Exec("INSERT INTO images(pathimage) VALUES (?)", img_path)
	if err != nil {
		return err
	}
	return nil
}
