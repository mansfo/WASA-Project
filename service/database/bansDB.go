package database

// Returns the userid of all the users banned by u
func (db *appdbimpl) GetBanned(u Userid) ([]UserSearched, error) {
	rows, err := db.c.Query("SELECT banned FROM bans WHERE banner = ?", u.IdUser)
	if err != nil {
		// Error executing the query
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []UserSearched
	// Add every row returned by the query in the array
	for rows.Next() {
		var userID Userid
		// Stores the content of a row in the variable userID
		err = rows.Scan(&userID.IdUser)
		if err != nil {
			return nil, err
		}

		var user UserSearched
		user.Usr.IdUser = userID

		// Returns the user's infos
		name, propic, err := db.GetUserInfos(userID)
		if err != nil {
			return nil, err
		}
		user.Usr.Name = name
		user.Img = propic

		// Adds the variable user in the result
		res = append(res, user)
	}

	// Check for error in rows
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Return the result containing the userid of the users banned by u
	return res, nil
}

// Add v to the collection of users banned by u
func (db *appdbimpl) BanUser(u Userid, v Userid) error {
	_, err := db.c.Exec("INSERT INTO bans (banner, banned) VALUES (?, ?)", u.IdUser, v.IdUser)
	if err != nil {
		// Error executing the function
		return err
	}
	// remove u's likes on all v's post
	rows, err := db.c.Query("SELECT l.photo FROM likes as l, photos as p WHERE p.photoid = l.photo AND l.user = ? AND p.user = ?", u.IdUser, v.IdUser)
	if err != nil {
		// Error executing the query
		return err
	}
	defer func() { _ = rows.Close() }()
	var unlikes []PhotoId
	for rows.Next() {
		var phId PhotoId
		err = rows.Scan(&phId.IdPhoto)
		if err != nil {
			return err
		}
		unlikes = append(unlikes, phId)
	}
	// Check for error in rows
	err = rows.Err()
	if err != nil {
		return err
	}

	// removes u's comments on v's posts
	rows2, err := db.c.Query("SELECT c.commentid FROM comments as c, photos as p WHERE c.photo = p.photoid AND p.user = ? AND c.user = ?", v.IdUser, u.IdUser)
	if err != nil {
		// Error executing the query
		return err
	}
	defer func() { _ = rows2.Close() }()
	var uncomms []CommentId
	for rows2.Next() {
		var cid CommentId
		err = rows2.Scan(&cid.IdComment)
		if err != nil {
			return err
		}
		uncomms = append(uncomms, cid)
	}
	// Check for error in rows
	err = rows2.Err()
	if err != nil {
		return err
	}

	// Remove all the likes found
	for i := 0; i < len(unlikes); i++ {
		err = db.UnlikePhoto(u, unlikes[i])
		if err != nil {
			return err
		}
	}

	// Remove all the comments found
	for j := 0; j < len(uncomms); j++ {
		err = db.UncommentPhoto(uncomms[j])
		if err != nil {
			return err
		}
	}

	// remove v's likes on all u's post
	rows3, err := db.c.Query("SELECT l.photo FROM likes as l, photos as p WHERE p.photoid = l.photo AND l.user = ? AND p.user = ?", v.IdUser, u.IdUser)
	if err != nil {
		// Error executing the query
		return err
	}
	defer func() { _ = rows3.Close() }()
	var unlikes2 []PhotoId
	for rows3.Next() {
		var phId3 PhotoId
		err = rows3.Scan(&phId3.IdPhoto)
		if err != nil {
			return err
		}
		unlikes2 = append(unlikes2, phId3)
	}
	// Check for error in rows
	err = rows3.Err()
	if err != nil {
		return err
	}

	// removes v's comments on u's posts
	rows4, err := db.c.Query("SELECT c.commentid FROM comments as c, photos as p WHERE c.photo = p.photoid AND p.user = ? AND c.user = ?", u.IdUser, v.IdUser)
	if err != nil {
		// Error executing the query
		return err
	}
	defer func() { _ = rows4.Close() }()
	var uncomms2 []CommentId
	for rows4.Next() {
		var cid2 CommentId
		err = rows4.Scan(&cid2.IdComment)
		if err != nil {
			return err
		}
		uncomms2 = append(uncomms2, cid2)
	}
	// Check for error in rows
	err = rows4.Err()
	if err != nil {
		return err
	}

	// Remove all the likes found
	for i := 0; i < len(unlikes2); i++ {
		err = db.UnlikePhoto(v, unlikes2[i])
		if err != nil {
			return err
		}
	}

	// Remove all the comments found
	for j := 0; j < len(uncomms2); j++ {
		err = db.UncommentPhoto(uncomms2[j])
		if err != nil {
			return err
		}
	}

	// Check if user u follows v, if so u unfollows v
	followStatus, err := db.IsFollowed(u, v)
	if err != nil {
		return err
	}
	if followStatus {
		err = db.UnfollowUser(u, v)
		if err != nil {
			return err
		}
	}

	// Check if user v follows u, if so v unfollows u
	followStatusInv, err := db.IsFollowed(v, u)
	if err != nil {
		return err
	}
	if followStatusInv {
		err = db.UnfollowUser(v, u)
		if err != nil {
			return err
		}
	}

	return nil
}

// Remove v from the collection of the users banned by u
func (db *appdbimpl) UnbanUser(u Userid, v Userid) error {
	_, err := db.c.Exec("DELETE FROM bans WHERE banner = ? AND banned = ?", u.IdUser, v.IdUser)
	if err != nil {
		// Checking for errors
		return err
	}
	return nil
}

// Check if user u is banned by user v
func (db *appdbimpl) CheckBan(u Userid, v Userid) (bool, error) {
	if u == v {
		return false, nil
	}
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM bans WHERE banner = ? AND banned = ?", v.IdUser, u.IdUser).Scan(&count)
	if err != nil {
		// Error executing the query
		return true, err
	}
	// if count == 1 the user is banned
	if count == 1 {
		return true, nil
	}
	return false, nil
}
