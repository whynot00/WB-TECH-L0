package httpHandler

import (
	"kaf-interface/internal/orders/models"
	"log/slog"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Logging(c *gin.Context) {

	startTime := time.Now() // запись текущего времени запроса

	c.Next()

	// модель лога
	// все ошибки возникающие внутри хендлера записываютс в gin.context.Errors
	logInfo := models.Logger{
		RequestMethod: c.Request.Method,
		StatusCode:    c.Writer.Status(),
		ExecutonTime:  time.Since(startTime).String(),
		RequestID:     requestid.Get(c),
		RequestURI:    c.Request.Host + c.Request.RequestURI,
		RequestSize:   c.Writer.Size(),
		Errors:        len(c.Errors),
	}

	switch {
	// если ошибок > 0, то они логгируются отдельно, логгируя requestid и саму ошибку
	case logInfo.Errors != 0:
		for _, err := range c.Errors {
			h.logger.Error("error",
				slog.String("request_id", logInfo.RequestID),
				slog.String("error", err.Error()))
		}
		h.logger.Warn("request - with errors", slog.Any("parameters", logInfo))

	default:
		// штатное логгирование запроса
		h.logger.Info("request - success", slog.Any("parameters", logInfo))
	}

	c.Next()
}
