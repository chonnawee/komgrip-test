package adapters

import (
	"errors"
	"komgrip-test/usecases"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

func (h *BeersHandler) validateExtension(filename string, allowExtensions []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowExt := range allowExtensions {
		if ext == allowExt {
			return true
		}
	}
	return false
}

func (h *BeersHandler) getFilePath(folderName string) (string, error) {
	path := filepath.Join("storage", folderName)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return path, nil
}

func (h *BeersHandler) generateFileName(extenstion string) string {
	timestamp := time.Now().Format("2006_01_02")
	randomNumber := rand.Intn(1000)

	filename := timestamp + "_" + strconv.Itoa(randomNumber) + extenstion

	return filename
}

func (h *BeersHandler) storeFile(c *fiber.Ctx) (path string, err error) {
	file, _ := c.FormFile("beer_img")
	pathResonse := ""
	if file != nil {
		if !h.validateExtension(file.Filename, []string{".jpg", ".jpeg", ".png"}) {
			return "", errors.New("invalid image type")
		}
		path, err := h.getFilePath("beers")
		if err != nil {
			return "", err
		}
		extenstion := filepath.Ext(file.Filename)
		filename := h.generateFileName(extenstion)
		filepath := path + "/" + filename
		pathResonse = filepath
		err = c.SaveFile(file, filepath)
		if err != nil {
			return "", err
		}
	}
	return pathResonse, nil
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
	filepath, err := h.storeFile(c)
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	var request usecases.BeersRequest
	err = c.BodyParser(&request)
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	request.BeerImgPath = filepath
	err = h.beerUseCase.CreateBeer(request)
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return h.SendSuccessResponse(c, nil)
}

func (h *BeersHandler) GetBeers(c *fiber.Ctx) error {
	beerName := c.Query("beer_name")
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	responses, err := h.beerUseCase.GetBeers(usecases.GetBeersRequest{
		BeerName: beerName,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return h.SendSuccessResponse(c, responses)
}

func (h *BeersHandler) UpdateBeer(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	filepath, err := h.storeFile(c)
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	var request usecases.BeersRequest
	err = c.BodyParser(&request)
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	request.BeerImgPath = filepath
	err = h.beerUseCase.UpdateBeer(id, request)
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return h.SendSuccessResponse(c, nil)
}

func (h *BeersHandler) DeleteBeer(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	err = h.beerUseCase.DeleteBeer(uint64(id))
	if err != nil {
		return h.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return h.SendSuccessResponse(c, nil)
}
