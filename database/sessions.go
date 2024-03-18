package database

// func (db DB) CreateSession(username string, token string, role string, createdAt time.Time, expiresAt time.Time) error {

// 	query := "INSERT INTO sessions (username, token, role, created_at, expired_at) VALUES ($1, $2, $3, $4,$5)"
// 	values := []interface{}{username, token, role, createdAt, expiresAt}
// 	log.Println("query:CreateSession:", query, values)
// 	_, err := db.pool.Exec(context.Background(), query, values...)
// 	if err == pgx.ErrNoRows {
// 		log.Println("err:CreateSession:pgx.ErrNoRows:", err)
// 		return nil
// 	}
// 	if err != nil {
// 		log.Println("err:CreateSession:", err)
// 		return err
// 	}

// 	log.Println("sussess:CreateSession", username)

// 	return err
// }

// func (db *DB) GetSessionByUsername(username string) (*model.Session, error) {
// 	var session model.Session
// 	err := db.pool.QueryRow(context.Background(), `
//         SELECT user, token, role, created_at,  expires_at
//         FROM sessions
//         WHERE username = $1
//     `, username).Scan(&session.User, &session.Role, &session.CreatedAt, &session.ExpiresAt)

// 	if err == pgx.ErrNoRows {
// 		log.Println("err:GetSessionByUsername:pgx.ErrNoRows:", err)
// 		return nil, nil
// 	}
// 	if err != nil {
// 		log.Println("err:GetSessionByUsername:", err)
// 		return nil, err
// 	}

// 	return &session, nil
// }

// func (db *DB) DeleteSession(sessionToken string) error {
// 	_, err := db.pool.Exec(context.Background(), `
//         DELETE FROM sessions WHERE session_token = $1
//     `, sessionToken)

// 	if err == pgx.ErrNoRows {
// 		log.Println("err:DeleteSession:pgx.ErrNoRows:", err)
// 		return nil
// 	}
// 	if err != nil {
// 		log.Println("err:DeleteSession:", err)
// 		return err
// 	}

// 	log.Println("sussess:DeleteSession", sessionToken)
// 	return err
// }
