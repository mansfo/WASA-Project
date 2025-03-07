package database

// Function that creates a new image in the db and returns its identifier and an error
func (db *appdbimpl) NewImage(dt string) (int64, error) {
	res, err := db.c.Exec("INSERT INTO images(uploadtime, pathimage) VALUES (?, NULL)", dt)
	if err != nil {
		// Error executing the query
		return -1, err
	}

	// Returns the identifier of the image
	imageId, err := res.LastInsertId()
	if err != nil {
		// Error while getting id returned by last db operation
		return -1, err
	}

	return imageId, nil
}

func (db *appdbimpl) SetImagePath(imid ImageId, path string) error {
	_, err := db.c.Exec("UPDATE images SET pathimage = ? WHERE imageid = ?", path, imid.IdImage)
	if err != nil {
		// Error executing the query
		return err
	}
	return nil
}

// Function that returns the photos published by a user and an error
func (db *appdbimpl) GetUsersPhotos(u Userid, v Userid) (int, []Photo, error) {
	// Stores in the variable rows the photos published by the user in reverse chronological order
	rows, err := db.c.Query("SELECT p.photoid, p.user, p.date, i.pathimage FROM photos as p, images as i WHERE p.image = i.imageid AND p.user = ? AND p.user NOT IN (SELECT banner FROM bans WHERE banned = ?) AND p.user NOT IN (SELECT banned FROM bans WHERE banner = ?) ORDER BY p.date DESC", v.IdUser, u.IdUser, u.IdUser)
	if err != nil {
		return -1, nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []Photo
	// Reading the result of the query
	for rows.Next() {
		var ph Photo
		// Stores the content of a row in the variable ph
		err = rows.Scan(&ph.IdPhoto.IdPhoto, &ph.User.Usr.IdUser.IdUser, &ph.Date, &ph.ImgPath)
		if err != nil {
			// Error scanning the row
			return -1, nil, err
		}

		name, propic, err := db.GetUserInfos(ph.User.Usr.IdUser)
		if err != nil {
			// Error getting infos
			return -1, nil, err
		}
		ph.User.Usr.Name = name
		ph.User.Img = propic

		// Returns the likes on the photo ph
		like_count, likes, err := db.GetLikeOnPhoto(u, ph.IdPhoto)
		if err != nil {
			return -1, nil, err
		}
		ph.Likes = likes
		ph.LikeCounter = like_count

		// Returns the comments on the photo ph
		comms_count, comms, err := db.GetComments(u, ph.IdPhoto)
		if err != nil {
			return -1, nil, err
		}
		ph.Comments = comms
		ph.CommentCounter = comms_count

		// Check if user u likes the photo
		liked, err := db.CheckLike(u, ph.IdPhoto)
		if err != nil {
			return -1, nil, err
		}
		if liked == 1 {
			ph.IsLiked = true
		} else {
			ph.IsLiked = false
		}

		// Adds the photo ph with all of its components to the result array
		res = append(res, ph)
	}

	// Check for error in rows
	err = rows.Err()
	if err != nil {
		return -1, nil, err
	}

	// Returns how many photos are published by the user v
	counter, err := db.CountPhotos(u, v)
	if err != nil {
		return -1, nil, err
	}
	return counter, res, err
}

// Function that returns how many photos are published by user v and an error
func (db *appdbimpl) CountPhotos(u Userid, v Userid) (int, error) {
	var counter int
	err := db.c.QueryRow("SELECT COUNT(*) FROM photos WHERE user = ? AND user NOT IN (SELECT banner FROM bans WHERE banned = ?) AND user NOT IN (SELECT banned FROM bans WHERE banner = ?)", v.IdUser, u.IdUser, u.IdUser).Scan(&counter)
	if err != nil {
		return -1, err
	}
	return counter, nil
}

// Function that uploads a new photo in the user profile; returns the photo counter, the photos and an error
func (db *appdbimpl) UploadPhoto(u Userid, i ImageId, d string) (int64, error) {
	res, err := db.c.Exec("INSERT INTO photos (user, image, date) VALUES (?, ?, ?)", u.IdUser, i.IdImage, d)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		// Error while getting id returned by last db operation
		return -1, err
	}
	return id, nil
}

// Function that removes a photo from the db; returns an error
func (db *appdbimpl) RemovePhoto(p PhotoId) error {
	var pict ImageId
	err := db.c.QueryRow("SELECT image FROM photos WHERE photoid = ?", p.IdPhoto).Scan(&pict.IdImage)
	if err != nil {
		return err
	}
	// Removes the photo from the db
	_, err = db.c.Exec("DELETE FROM photos WHERE photoid = ?", p.IdPhoto)
	if err != nil {
		return err
	}
	// Removes the image from the db
	_, err = db.c.Exec("DELETE FROM images WHERE imageid = ?", pict.IdImage)
	if err != nil {
		return err
	}

	return nil
}
