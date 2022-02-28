package user

type GetUserRequest struct {
	ID        *int    `json:"id" example:"1"`
	FirstName *string `json:"firstName" example:"Khanapat"`
}
