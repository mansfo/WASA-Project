package database

// Function that returns the infos of the users that like the photo p minus the one that banned u or are banned by u
func (db *appdbimpl) GetLikeOnPhoto(u Userid, p PhotoId) (int, []UserSearched, error) {
	rows, err := db.c.Query("SELECT user FROM likes WHERE photo = ? AND user NOT IN (SELECT banned FROM bans WHERE banner = ?) AND user NOT IN (SELECT banner FROM bans WHERE banned = ?)", p.IdPhoto, u.IdUser, u.IdUser)
	if err != nil {
		return -1, nil, err
	}
	defer func() { _ = rows.Close() }()
	var res []UserSearched
	// Reading the result of the query
	for rows.Next() {
		var id Userid
		// Stores the content of a row in the variable id
		err := rows.Scan(&id.IdUser)
		if err != nil {
			return -1, nil, err
		}
		var user UserSearched
		user.Usr.IdUser = id
		name, propic, err := db.GetUserInfos(id)
		if err != nil {
			return -1, nil, err
		}
		user.Usr.Name = name
		user.Img = propic
		// Adds name to the array that will be returned
		res = append(res, user)
	}

	// Check for error in rows
	err = rows.Err()
	if err != nil {
		return -1, nil, err
	}

	// Returns how many users like the photo
	counter, err := db.CountLikes(u, p)
	if err != nil {
		return -1, nil, err
	}

	return counter, res, nil
}

// Function that add u's username to the list of the users that like the photo; returns an error
func (db *appdbimpl) LikePhoto(u Userid, p PhotoId) error {
	// Checks if the user is trying to like one of his photos
	selfLike, err := db.SelfLike(u, p)
	if err != nil {
		return err
	}
	if selfLike {
		return ErrSelfLike
	}
	// Creates a new instance of likes
	_, err = db.c.Exec("INSERT INTO likes (user, photo) VALUES (?, ?)", u.IdUser, p.IdPhoto)
	if err != nil {
		return err
	}

	return nil
}

// Function that removes a like from a photo; returns an error
func (db *appdbimpl) UnlikePhoto(u Userid, p PhotoId) error {
	_, err := db.c.Exec("DELETE FROM likes WHERE (user = ? AND photo = ?)", u.IdUser, p.IdPhoto)
	if err != nil {
		return err
	}

	return nil
}

// Function that counts how many users like a photo, returns also an error
func (db *appdbimpl) CountLikes(u Userid, p PhotoId) (int, error) {
	var num int
	err := db.c.QueryRow("SELECT COUNT (*) FROM likes WHERE photo = ? AND user NOT IN (SELECT banner FROM bans WHERE banned = ?) AND user NOT IN (SELECT banned FROM bans WHERE banner = ?)", p.IdPhoto, u.IdUser, u.IdUser).Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

// Function that checks if a user likes a photo, returns also an error
func (db *appdbimpl) CheckLike(u Userid, p PhotoId) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM likes WHERE photo = ? AND user = ?", p.IdPhoto, u.IdUser).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Function that checks if a user is trying to like one of his photos, returns also an error
func (db *appdbimpl) SelfLike(u Userid, p PhotoId) (bool, error) {
	var user Userid
	err := db.c.QueryRow("SELECT user FROM photos WHERE photoid = ?", p.IdPhoto).Scan(&user.IdUser)
	if err != nil {
		return true, err
	}
	if user != u {
		return false, nil
	}
	return true, nil
}
