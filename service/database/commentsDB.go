package database

import (
	"time"
)

// Returns the comments published on the photo p and how many comments there are (minus the comments published by users banned by u or users that banned u)
func (db *appdbimpl) GetComments(u Userid, p PhotoId) (int, []Comment, error) {
	// Returns the comments published on the photo (minus the comments published by users banned by u or users that banned u)
	rows, err := db.c.Query("SELECT * FROM comments WHERE photo = ? AND user NOT IN (SELECT banner FROM bans WHERE banned = ?) AND user NOT IN (SELECT banned FROM bans WHERE banner = ?)", p.IdPhoto, u.IdUser, u.IdUser)
	if err != nil {
		// Error executing the query
		return -1, nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []Comment

	// Add every row returned by the query in the result array
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.IdComment.IdComment, &comment.IdPhoto.IdPhoto, &comment.User.Usr.IdUser.IdUser, &comment.Text, &comment.Date)
		// Save in comment the content of one of the rows
		if err != nil {
			// Checking for errors
			return -1, nil, err
		}

		var usr UserSearched
		usr.Usr.IdUser = comment.User.Usr.IdUser

		name, propic, err := db.GetUserInfos(comment.User.Usr.IdUser)
		if err != nil {
			// Checking for errors
			return -1, nil, err
		}
		usr.Usr.Name = name
		usr.Img = propic

		comment.User.Usr.Name.Name = usr.Usr.Name.Name
		comment.User.Img = usr.Img

		// Adds the comment in the array
		res = append(res, comment)
	}

	// Check for error in rows
	err = rows.Err()
	if err != nil {
		return -1, nil, err
	}

	// Returs how many comments are published on the photo (minus the comments published by users banned by u or users that banned u)
	count, err := db.CountComments(u, p)
	if err != nil {
		// Checking for errors
		return -1, nil, err
	}
	// Returns the comment counter, the array containing the comments and an error (that should be nil)
	return count, res, err
}

// Returns how many comments are published on photo p (minus the comments published by users banned by u or users that banned u)
func (db *appdbimpl) CountComments(u Userid, p PhotoId) (int, error) {
	var counter int
	err := db.c.QueryRow("SELECT COUNT(*) FROM comments WHERE photo = ? AND user NOT IN (SELECT banner FROM bans WHERE banned = ?) AND user NOT IN (SELECT banned FROM bans WHERE banner = ?)", p.IdPhoto, u.IdUser, u.IdUser).Scan(&counter)
	if err != nil {
		// Checking for errors
		return -1, err
	}
	return counter, nil
}

// Add a comment on the photo p
func (db *appdbimpl) CommentPhoto(u Userid, p PhotoId, s string, dt string) (Comment, error) {
	res, err := db.c.Exec("INSERT INTO comments (photo, user, text, date) VALUES (?, ?, ?, ?)", p.IdPhoto, u.IdUser, s, dt)
	if err != nil {
		// Checking for errors
		return Comment{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		// Checking for errors
		return Comment{}, err
	}

	var usr UserSearched
	usr.Usr.IdUser = u
	name, propic, err := db.GetUserInfos(u)
	if err != nil {
		// Checking for errors
		return Comment{}, err
	}

	usr.Usr.Name = name
	usr.Img = propic

	date, err := time.Parse(time.RFC3339, dt)
	if err != nil {
		return Comment{}, err
	}
	var comment Comment
	comment.IdPhoto = p
	comment.IdComment.IdComment = id
	comment.Text = s
	comment.Date = date
	comment.User = usr

	return comment, nil
}

// Removes a comment from the photo p
func (db *appdbimpl) UncommentPhoto(c CommentId) error {
	_, err := db.c.Exec("DELETE FROM comments WHERE commentid = ?", c.IdComment)
	if err != nil {
		// Checking for errors
		return err
	}

	return nil
}

// Function that checks if the user u is the one that has written the comment c, returns also an error
func (db *appdbimpl) IsCommentAuthor(u Userid, c CommentId) (bool, error) {
	var num int
	err := db.c.QueryRow("SELECT COUNT(*) FROM comments WHERE user = ? AND commentid = ?", u.IdUser, c.IdComment).Scan(&num)
	if err != nil {
		return false, err
	}
	if num != 1 {
		return false, nil
	}
	return true, nil
}
