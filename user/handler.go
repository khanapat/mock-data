package user

import "github.com/gofiber/fiber/v2"

type userHandler struct {
	UserRepository UserRepository
}

func NewUserHandler(userRepository UserRepository) *userHandler {
	return &userHandler{
		UserRepository: userRepository,
	}
}

// GetUser
// @Summary Get User
// @Description get user by id or first name
// @Tags User
// @Accept json
// @Produce json
// @Param id query int false "ID"
// @Param firstName query string false "FirstName"
// @Success 200 {array} user.User "Success"
// @Router /user [get]
func (s *userHandler) GetUser(c *fiber.Ctx) error {
	var req GetUserRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
	}
	m := make(map[string]interface{})
	if req.ID != nil {
		m["id"] = req.ID
	}
	if req.FirstName != nil {
		m["first_name"] = req.FirstName
	}
	users, err := s.UserRepository.QueryUser(c.Context(), m)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"data": users,
	})
}
