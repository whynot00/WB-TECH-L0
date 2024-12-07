package httpHandler

import (
	"kaf-interface/internal/orders/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetOrderByID(c *gin.Context) {

	order, err := h.service.GetOrderByID(c.Param("id")) // обращение к сервису
	if err != nil {
		switch err {
		case models.KeyNotExistsError:
			// в случае если запрос коректен, но по данному uid нет информации
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "The request is correct. There is no data in the database."})
			return
		default:
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

	}

	c.HTML(http.StatusOK, "index.html", order)
}
