package adapters

import (
	"komgrip-test/usecases"

	"github.com/gofiber/fiber/v2"
)

type BeersHandler struct {
	beerUseCase usecases.BeersUseCase
}

type Response struct {
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewBeersHandler(beerUseCase usecases.BeersUseCase) *BeersHandler {
	return &BeersHandler{beerUseCase: beerUseCase}
}

func (h *BeersHandler) SendSuccessResponse(c *fiber.Ctx, data interface{}) error {
	response := Response{
		Status: true,
		Code:   fiber.StatusOK,
		Data:   data,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *BeersHandler) SendErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	response := Response{
		Status:  false,
		Code:    statusCode,
		Message: message,
	}
	return c.Status(statusCode).JSON(response)
}

func (h *BeersHandler) CreateBeer(c *fiber.Ctx) error {
	var requests usecases.BeersRequest
	err := c.BodyParser(&requests)
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	err = h.beerUseCase.CreateBeer(requests)
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return h.SendSuccessResponse(c, nil)
}

func (h *BeersHandler) GetBeers(c *fiber.Ctx) error {
	beerName := c.Query("beer_name")
}

func (h *BeersHandler) UpdateBeer(c *fiber.Ctx) error
func (h *BeersHandler) DeleteBeer(c *fiber.Ctx) error
