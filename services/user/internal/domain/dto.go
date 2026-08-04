package domain

type GetUsersRequest struct {
	IDs []int64
}

type User struct {
	ID int64
	FirstName string
	LastName string
	Username string
}

type GetUsersResponse struct {
	Users []User
}
