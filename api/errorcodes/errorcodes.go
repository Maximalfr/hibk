package errorcodes

import "net/http"

type errorStruct struct {
	Status  int    `json:"status"`
	Message string `json:"message" omitempty`
}

func OK() (int, errorStruct) {
	return http.StatusOK, errorStruct{0, "ok"}
}

// Auth
func BadPassword() (int, errorStruct) {
	return http.StatusUnauthorized, errorStruct{1, "bad password"}
}

func BadCredentials() (int, errorStruct) {
	return http.StatusUnauthorized, errorStruct{2, "bad credentials"}
}

func UsernameAlreadyExists() (int, errorStruct) {
	return http.StatusUnauthorized, errorStruct{3, "username already exists"}
}

func InternalError(error string) (int, errorStruct) {
	return http.StatusInternalServerError, errorStruct{50, error}
}

// jwtUtils
func MissingAuthorizationHeader() (int, errorStruct) {
	return http.StatusUnauthorized, errorStruct{10, "missing authorization header"}
}

func ErrorToken(error string) (int, errorStruct) {
	return http.StatusUnauthorized, errorStruct{11, "error verifying JWT token: " + error}
}

func ErrorTokenExpired() (int, errorStruct) {
	return http.StatusUnauthorized, errorStruct{12, "token is either expired or not active yet"}
}

//General
func BadRequest() (int, errorStruct) {
	return http.StatusBadRequest, errorStruct{20, "bad request"}
}
