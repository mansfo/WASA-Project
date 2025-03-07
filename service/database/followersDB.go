package database

// Function that returns v's followers minus the ones banned by u or that have banned u
func (db *appdbimpl) GetFollowers(u Userid, v Userid) (int, []UserSearched, error) {
	rows, err := db.c.Query("SELECT follower FROM followers WHERE followed = ? AND follower NOT IN (SELECT banner FROM bans WHERE banned = ?) AND follower NOT IN (SELECT banned FROM bans WHERE banner = ?)", v.IdUser, u.IdUser, u.IdUser)
	if err != nil {
		return -1, nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []UserSearched
	// Reading the result of the query
	for rows.Next() {
		var userID Userid
		// Stores the content of a row in the variable user
		err = rows.Scan(&userID.IdUser)
		if err != nil {
			return -1, nil, err
		}

		var user UserSearched
		user.Usr.IdUser = userID

		// Returns the user's infos
		name, propic, err := db.GetUserInfos(userID)
		if err != nil {
			return -1, nil, err
		}
		user.Usr.Name = name
		user.Img = propic

		// Adds the variable user in the result
		res = append(res, user)
	}

	// Check for error in rows
	err = rows.Err()
	if err != nil {
		return -1, nil, err
	}

	// Returns how many users follow v
	counter, err := db.CountFollowers(u, v)
	if err != nil {
		return -1, nil, err
	}
	return counter, res, nil
}

// Function that counts how many users are following user v minus the users banned by u or that have banned u
func (db *appdbimpl) CountFollowers(u Userid, v Userid) (int, error) {
	var result int
	err := db.c.QueryRow("SELECT COUNT(*) FROM followers WHERE followed = ? AND follower NOT IN (SELECT banner FROM bans WHERE banned = ?) AND follower NOT IN (SELECT banned FROM bans WHERE banner = ?)", v.IdUser, u.IdUser, u.IdUser).Scan(&result)
	if err != nil {
		return -1, err
	}
	return result, nil
}

// Function that returns v's following minus the ones banned by u or that have banned u
func (db *appdbimpl) GetFollowing(u Userid, v Userid) (int, []UserSearched, error) {
	rows, err := db.c.Query("SELECT followed FROM followers WHERE follower = ? AND followed NOT IN (SELECT banner FROM bans WHERE banned = ?) AND followed NOT IN (SELECT banned FROM bans WHERE banner = ?)", v.IdUser, u.IdUser, u.IdUser)
	if err != nil {
		return -1, nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []UserSearched
	// Reading the result of the query
	for rows.Next() {
		var userID Userid
		// Stores the content of a row in the variable user
		err = rows.Scan(&userID.IdUser)
		if err != nil {
			return -1, nil, err
		}

		var user UserSearched
		user.Usr.IdUser = userID

		// Returns the user's infos
		name, propic, err := db.GetUserInfos(userID)
		if err != nil {
			return -1, nil, err
		}
		user.Usr.Name = name
		user.Img = propic

		// Adds the variable user in the result
		res = append(res, user)
	}

	// Check for error in rows
	err = rows.Err()
	if err != nil {
		return -1, nil, err
	}

	// Returns how many users are followed by v
	counter, err := db.CountFollowing(u, v)
	if err != nil {
		return -1, nil, err
	}
	return counter, res, nil
}

// Function that count how many users are followed by v minus the ones that are banned by u or that have banned u
func (db *appdbimpl) CountFollowing(u Userid, v Userid) (int, error) {
	var result int
	err := db.c.QueryRow("SELECT COUNT(*) FROM followers WHERE follower = ? AND followed NOT IN (SELECT banner FROM bans WHERE banned = ?) AND followed NOT IN (SELECT banned FROM bans WHERE banner = ?)", v.IdUser, u.IdUser, u.IdUser).Scan(&result)
	if err != nil {
		return -1, err
	}
	return result, nil
}

// Function that adds u in v's followers
func (db *appdbimpl) FollowUser(u Userid, v Userid) error {
	_, err := db.c.Exec("INSERT INTO followers (follower, followed) VALUES (?, ?)", u.IdUser, v.IdUser)
	if err != nil {
		return err
	}
	return nil
}

// Function that removes u from v's followers
func (db *appdbimpl) UnfollowUser(u Userid, v Userid) error {
	_, err := db.c.Exec("DELETE FROM followers WHERE follower = ? AND followed = ?", u.IdUser, v.IdUser)
	if err != nil {
		return err
	}
	return nil
}

// Function that checks if u is one of v's followers
func (db *appdbimpl) IsFollowed(u Userid, v Userid) (bool, error) {
	var num int
	err := db.c.QueryRow("SELECT COUNT(*) FROM followers WHERE follower = ? AND followed = ?", u.IdUser, v.IdUser).Scan(&num)
	if err != nil {
		return false, err
	}
	// If num == 1 there is an instance of followers in which follower = u and followed = v, so u is one of v's followers
	if num != 1 {
		return false, nil
	}
	return true, nil
}
