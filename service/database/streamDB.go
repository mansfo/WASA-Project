package database

// Function that returns the stream of the user u and an error
func (db *appdbimpl) GetMyStream(u Userid) ([]Photo, error) {
	// Stores in rows all the photos published by users followed by u in reverse chronological order
	rows, err := db.c.Query("SELECT p.photoid, p.user, p.date, i.pathimage FROM photos as p, images as i WHERE p.image = i.imageid AND p.user IN (SELECT followed FROM followers WHERE follower = ?) ORDER BY p.date DESC", u.IdUser)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []Photo
	// Reads the result of the query
	for rows.Next() {
		var photo Photo
		// Stores in the variable photo the content of a row
		err := rows.Scan(&photo.IdPhoto.IdPhoto, &photo.User.Usr.IdUser.IdUser, &photo.Date, &photo.ImgPath)
		if err != nil {
			return nil, err
		}

		name, propic, err := db.GetUserInfos(photo.User.Usr.IdUser)
		if err != nil {
			// Error getting infos
			return nil, err
		}
		photo.User.Usr.Name = name
		photo.User.Img = propic

		// Returns the like on the photo
		like_count, likes, err := db.GetLikeOnPhoto(u, photo.IdPhoto)
		if err != nil {
			return nil, err
		}
		photo.Likes = likes
		photo.LikeCounter = like_count

		// Returns the comment on the photo
		comm_count, comments, err := db.GetComments(u, photo.IdPhoto)
		if err != nil {
			return nil, err
		}
		photo.Comments = comments
		photo.CommentCounter = comm_count

		// Check if user u likes the photo
		liked, err := db.CheckLike(u, photo.IdPhoto)
		if err != nil {
			return nil, err
		}
		if liked == 1 {
			photo.IsLiked = true
		} else {
			photo.IsLiked = false
		}

		// Adds the photo with all of its components to the result array
		res = append(res, photo)
	}

	// Check for error in rows
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, err
}
