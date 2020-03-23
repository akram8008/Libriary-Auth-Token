package users

const selectLoginById = `SELECT id FROM users WHERE login = $1;`



const insertInUsers = `INSERT INTO users(name, login, password, role) VALUES ($1, $2, $3, $4);`


const selectUserByLogin = `SELECT id,name,password FROM users WHERE login = $1;`



type Payload struct {
	Id    int64    `json:"id"`
	Exp   int64    `json:"exp"`
	Roles []string `json:"roles"`
}

type RequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseDTO struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Token    string    `json:"token"`
}

type ErrorDTO struct {
	Error string `json:"error"`
}
