package handler

import (
	"mini_jira/contract"
	"mini_jira/dto"
	errs "mini_jira/pkg/error"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailVerificationController struct {
	service contract.EmailVerificationService
}

func (c *EmailVerificationController) InitService(s *contract.Service) {
	if s == nil || s.EmailVerification == nil {
		return
	}
	c.service = s.EmailVerification
}

// GET /verify-email?token=xxxxx
func (c *EmailVerificationController) VerifyEmail(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		HandleError(ctx, errs.BadRequest("token is required"))
		return
	}

	if err := c.service.VerifyEmail(token); err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Email verified successfully! You can now login.",
	})
}

// POST /resend-verification
func (c *EmailVerificationController) ResendVerification(ctx *gin.Context) {
	var payload dto.ResendVerificationRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		HandleError(ctx, errs.BadRequest("Invalid request body"))
		return
	}

	if err := c.service.ResendVerification(payload.Email); err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Verification email sent. Please check your inbox.",
	})
}
