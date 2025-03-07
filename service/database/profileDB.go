package database

import (
	"strings"
)

// Function that returns the username, the path to the user's profile picture and an error
func (db *appdbimpl) GetUserInfos(u Userid) (Username, string, error) {
	var name Username
	var profPictId int64
	err := db.c.QueryRow("SELECT username, profilePicture FROM users WHERE userid = ?", u.IdUser).Scan(&name.Name, &profPictId)
	if err != nil {
		// Checking for errors
		return Username{}, "", err
	}
	var imgPath string
	// Getting the image path related to the image id
	err = db.c.QueryRow("SELECT pathimage FROM images WHERE imageid = ?", profPictId).Scan(&imgPath)
	if err != nil {
		// Checking for errors
		return Username{}, "", err
	}
	return name, imgPath, nil
}

// Function that checks if a user exists in the database
func (db *appdbimpl) CheckUser(u Userid) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE userid = ?", u.IdUser).Scan(&count)
	if err != nil {
		return true, err
	}
	// If count == 1 there is a userid that matches the input
	if count == 1 {
		return true, nil
	}
	return false, nil
}

// Creates a new user in the db
func (db *appdbimpl) NewUser(u string) (User, bool, error) {
	// Check if u is not a userid or a username
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE userid = ? OR LOWER(username) = ?", strings.ToLower(u), strings.ToLower(u)).Scan(&count)
	if err != nil {
		return User{}, false, err
	}
	if count < 1 {
		// This is a new user
		_, err = db.c.Exec("INSERT INTO users(userid, username, profilePicture) VALUES (?, ?, 1)", strings.ToLower(u), u)
		if err != nil {
			return User{}, false, err
		}
		var id Userid
		id.IdUser = strings.ToLower(u)
		var name Username
		name.Name = u
		return User{IdUser: id, Name: name}, true, err
	} else {
		// There is a user whose userid is u or whose username is u (regardless of letter case)
		isUid, err := db.CheckUser(Userid{IdUser: strings.ToLower(u)})
		if err != nil {
			return User{}, false, err
		}
		if isUid {
			// There is a user whose userid is u
			var usr Username
			err = db.c.QueryRow("SELECT username FROM users WHERE userid = ?", strings.ToLower(u)).Scan(&usr.Name)
			if err != nil {
				return User{}, false, err
			}
			return User{IdUser: Userid{IdUser: strings.ToLower(u)}, Name: usr}, false, nil
		} else {
			// There is a user whose username is u (regardless of lettercase)
			var uid Userid
			err = db.c.QueryRow("SELECT userid FROM users WHERE LOWER(username) = ?", strings.ToLower(u)).Scan(&uid.IdUser)
			if err != nil {
				return User{}, false, err
			}
			// Get the real username of this user
			var realName Username
			err = db.c.QueryRow("SELECT username FROM users WHERE userid = ?", uid.IdUser).Scan(&realName.Name)
			if err != nil {
				return User{}, false, err
			}
			return User{IdUser: uid, Name: realName}, false, nil
		}
	}
}

// Check if a username exists in the database, regardless of letter case: if so, returns true
func (db *appdbimpl) CheckIfUsernameIsTaken(id Userid, username string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(username) FROM users WHERE LOWER(username) = ? AND userid <> ?", strings.ToLower(username), id.IdUser).Scan(&count)
	if err != nil {
		return true, err
	}
	if count > 0 {
		// There is a match between an existing username and the new username
		return true, nil
	}
	// There is no match
	return false, nil
}

// Updates the username of a user
func (db *appdbimpl) SetMyUserName(u Userid, un string) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE userid = ?", un, u.IdUser)
	return err
}

// Updates the profile picture of a user
func (db *appdbimpl) SetMyProfilePict(u Userid, newImage ImageId) error {
	// Set a new profile picture for the user
	_, err := db.c.Exec("UPDATE users SET profilePicture = ? WHERE userid = ?", newImage.IdImage, u.IdUser)
	if err != nil {
		return err
	}
	return nil
}

// Function that returns the profile of the user v and an error
func (db *appdbimpl) GetProfile(searcher Userid, searchedName Userid) (UserProfile, error) {
	// Returns the User whose userid is searchedName if searcher isn't banned by him/her or if searcher hasn't banned him/her
	rows, err := db.c.Query("SELECT * FROM users WHERE userid = ? AND userid NOT IN (SELECT banner FROM bans WHERE banned = ?)", searchedName.IdUser, searcher.IdUser)
	if err != nil {
		// Error executing the query
		return UserProfile{}, err
	}
	defer func() { _ = rows.Close() }()

	var res UserProfile
	var profPictId int64

	// Setting some fields of the result
	for rows.Next() {
		err2 := rows.Scan(&res.Owner.IdUser.IdUser, &res.Owner.Name.Name, &profPictId)
		if err2 != nil {
			// Checking for errors
			return UserProfile{}, err2
		}
	}

	// Check for error in rows
	err = rows.Err()
	if err != nil {
		return UserProfile{}, err
	}
	var prof string
	err = db.c.QueryRow("SELECT pathimage FROM images WHERE imageid = ?", profPictId).Scan(&prof)
	if err != nil {
		// Checking for errors
		return UserProfile{}, err
	}
	res.ProfilePict.IdImage.IdImage = profPictId
	res.ProfilePict.UriImage = prof

	banned, err := db.CheckBan(searchedName, searcher)
	if err != nil {
		// Checking for errors
		return UserProfile{}, err
	}
	if banned {
		// If is banned, return only userid, username and propic
		return res, nil
	}

	// Returns the photos published by the searched user
	num_ph, photos, err := db.GetUsersPhotos(searcher, searchedName)
	if err != nil {
		// Checking for errors
		return UserProfile{}, err
	}
	// The result returned by the function is assigned to the correct field
	res.Photos = photos
	res.PhotoCounter = num_ph

	// Returns the users that are following the user searched
	num_fwrs, followers, err := db.GetFollowers(searcher, searchedName)
	if err != nil {
		// Checking for errors
		return UserProfile{}, err
	}
	// The result returned by the function is assigned to the correct field
	res.Followers = followers
	res.FollowersCounter = num_fwrs

	// Returns the users followed by the user searched
	num_fwng, following, err := db.GetFollowing(searcher, searchedName)
	if err != nil {
		// Checking for errors
		return UserProfile{}, err
	}
	// The result returned by the function is assigned to the correct field
	res.Following = following
	res.FollowingCounter = num_fwng

	// Returns the complete profile of the user searched

	return res, err
}

// Search a user using his/her username or his /her userid, returns the profile of the users that match the patterne written and an error
func (db *appdbimpl) SearchUser(searcher Userid, searchedName Userid) ([]UserSearched, error) {
	var name Userid
	name.IdUser = searchedName.IdUser + "%"
	// Returns the users whose userid matches searchedName if searcher isn't banned by him/her or if searcher hasn't banned him/her
	rows, err := db.c.Query("SELECT * FROM users WHERE userid LIKE ? AND userid NOT IN (SELECT banned FROM bans WHERE banner = ?) AND userid NOT IN (SELECT banner FROM bans WHERE banned = ?)", name.IdUser, searcher.IdUser, searcher.IdUser)
	if err != nil {
		// Error executing the query
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []UserSearched

	// Scanning the results of the query
	for rows.Next() {
		var user UserSearched
		var profPictId int64
		// Stores the content of a row in the variable user and in the variable profPictId
		err2 := rows.Scan(&user.Usr.IdUser.IdUser, &user.Usr.Name.Name, &profPictId)
		if err2 != nil {
			// Checking for errors
			return nil, err2
		}

		// Getting the image path related to the image id
		err = db.c.QueryRow("SELECT pathimage FROM images WHERE imageid = ?", profPictId).Scan(&user.Img)
		if err != nil {
			// Checking for errors
			return nil, err
		}

		// Adds the user to the result array
		res = append(res, user)
	}

	// Check for error in rows
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Now it does the same with the users whose username matches searchedame

	// Returns the users whose username matches searchedName if searcher isn't banned by him/her or if searcher hasn't banned him/her and if this user is not in the list yet
	rows2, err := db.c.Query("SELECT * FROM users WHERE username LIKE ? AND userid NOT LIKE ? AND userid NOT IN (SELECT banned FROM bans WHERE banner = ?) AND userid NOT IN (SELECT banner FROM bans WHERE banned = ?)", name.IdUser, name.IdUser, searcher.IdUser, searcher.IdUser)
	if err != nil {
		// Error executing the query
		return nil, err
	}
	defer func() { _ = rows2.Close() }()
	// Scanning the results of the query
	for rows2.Next() {
		var user UserSearched
		var profPictId int64
		// Stores the content of a row in the variable user and in the variable profPictId
		err2 := rows2.Scan(&user.Usr.IdUser.IdUser, &user.Usr.Name.Name, &profPictId)
		if err2 != nil {
			// Checking for errors
			return nil, err2
		}

		// Getting the image path related to the image id
		err = db.c.QueryRow("SELECT pathimage FROM images WHERE imageid = ?", profPictId).Scan(&user.Img)
		if err != nil {
			// Checking for errors
			return nil, err
		}

		// Adds the user to the result array
		res = append(res, user)
	}

	// Check for error in rows
	err = rows2.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}
