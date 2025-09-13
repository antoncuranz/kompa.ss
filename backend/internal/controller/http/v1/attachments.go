package v1

import (
	"fmt"
	"kompass/internal/usecase"
	"kompass/pkg/logger"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AttachmentsV1 struct {
	uc  usecase.Attachments
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Get all attachments
// @ID          getAttachments
// @Tags  	    attachments
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Success     200 {object} []entity.Attachment
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/attachments [get]
func (r *AttachmentsV1) getAttachments(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %w", err))
	}

	attachments, err := r.uc.GetAttachments(ctx.UserContext(), int32(tripID))
	if err != nil {
		return errorResponse(ctx, fmt.Errorf("get attachments: %w", err))
	}

	return ctx.Status(http.StatusOK).JSON(attachments)
}

// @Summary     Download attachment by ID
// @ID          downloadAttachment
// @Tags  	    attachments
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       attachment_id path int true "Attachment ID"
// @Success     200 {object} entity.Attachment
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/attachments/{attachment_id} [get]
func (r *AttachmentsV1) downloadAttachment(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %w", err))
	}
	attachmentID, err := ctx.ParamsInt("attachment_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse attachment_id: %w", err))
	}

	attachment, err := r.uc.GetAttachmentByID(ctx.UserContext(), int32(tripID), int32(attachmentID))
	if err != nil {
		return errorResponse(ctx, fmt.Errorf("get attachment [id=%d]: %w", attachmentID, err))
	}

	ctx.Set("Content-Type", determineContentType(attachment.Name))
	ctx.Set("Content-Disposition", "inline; filename="+attachment.Name)
	return ctx.Send(attachment.Blob)
}

func determineContentType(filename string) string {
	ext := filepath.Ext(filename)
	if contentType := mime.TypeByExtension(ext); contentType != "" {
		return contentType
	}
	return "application/octet-stream"
}

type AttachmentsParam struct {
	Attachments []string `format:"binary"`
}

// @Summary     Add attachment
// @ID          postAttachment
// @Tags  	    attachments
// @Accept      multipart/form-data
// @Param       trip_id path int true "Trip ID"
// @Param       attachments formData AttachmentsParam true "attachment"
// @Success     204
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/attachments [post]
func (r *AttachmentsV1) postAttachment(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %w", err))
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["attachments"]

	for _, file := range files {
		fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

		_, err = r.uc.CreateAttachment(ctx.UserContext(), int32(tripID), file)
		if err != nil {
			return errorResponse(ctx, fmt.Errorf("create attachment: %w", err))
		}
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Delete attachment
// @ID          deleteAttachment
// @Tags  	    attachments
// @Param       trip_id path int true "Trip ID"
// @Param       attachment_id path int true "Attachment ID"
// @Success     204
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/attachments/{attachment_id} [delete]
func (r *AttachmentsV1) deleteAttachment(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %w", err))
	}
	attachmentID, err := ctx.ParamsInt("attachment_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse attachment_id: %w", err))
	}

	if err := r.uc.DeleteAttachment(ctx.UserContext(), int32(tripID), int32(attachmentID)); err != nil {
		return errorResponse(ctx, fmt.Errorf("delete attachment with id %d: %w", attachmentID, err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}
