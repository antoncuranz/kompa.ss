package attachments

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"travel-planner/internal/entity"
	"travel-planner/internal/repo"
)

type UseCase struct {
	repo repo.AttachmentsRepo
}

func New(r repo.AttachmentsRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) GetAttachments(ctx context.Context, tripID int32) ([]entity.Attachment, error) {
	return uc.repo.GetAttachments(ctx, tripID)
}

func (uc *UseCase) GetAttachmentByID(ctx context.Context, tripID int32, id int32) (entity.Attachment, error) {
	return uc.repo.GetAttachmentByID(ctx, tripID, id)
}

func (uc *UseCase) CreateAttachment(ctx context.Context, tripID int32, fileHeader *multipart.FileHeader) (entity.Attachment, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return entity.Attachment{}, fmt.Errorf("open file: %w", err)
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return entity.Attachment{}, fmt.Errorf("read file: %w", err)
	}
	return uc.repo.SaveAttachment(ctx, entity.Attachment{
		TripID: tripID,
		Name:   fileHeader.Filename,
		Blob:   fileBytes,
	})
}

func (uc *UseCase) DeleteAttachment(ctx context.Context, tripID int32, attachmentId int32) error {
	return uc.repo.DeleteAttachment(ctx, tripID, attachmentId)
}
