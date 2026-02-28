package auth

const (
	sqlCreateUser = `
INSERT INTO users (user_id, role, name, email, password_hash, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5, now(), now())
RETURNING id, created_at, updated_at;
`

	sqlGetByUserID = `
SELECT id, user_id, email, password_hash, role
FROM users
WHERE user_id = $1
LIMIT 1;
`

	sqlGetByEmail = `
SELECT id, user_id, email, password_hash, role
FROM users
WHERE email = $1
LIMIT 1;
`

	sqlUpdatePassword = `
UPDATE users SET password_hash = $1, updated_at = now() WHERE id = $2;
`
)
