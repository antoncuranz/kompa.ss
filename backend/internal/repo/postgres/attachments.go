package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"kompass/internal/entity"
	"kompass/pkg/postgres"
	"kompass/pkg/sqlc"
)

type AttachmentsRepo struct {
	Db      *pgxpool.Pool
	Queries *sqlc.Queries
}

func NewAttachmentsRepo(pg *postgres.Postgres) *AttachmentsRepo {
	return &AttachmentsRepo{
		pg.Pool,
		sqlc.New(pg.Pool),
	}
}

func (r *AttachmentsRepo) GetAttachments(ctx context.Context, tripID int32) ([]entity.Attachment, error) {
	attachments, err := r.Queries.GetAttachments(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("get all attachments [tripID=%d]: %w", tripID, err)
	}

	return mapAttachments(attachments), nil
}

func (r *AttachmentsRepo) GetAttachmentByID(ctx context.Context, tripID int32, attachmentID int32) (entity.Attachment, error) {
	attachment, err := r.Queries.GetAttachmentByID(ctx, sqlc.GetAttachmentByIDParams{TripID: tripID, ID: attachmentID})
	if err != nil {
		return entity.Attachment{}, fmt.Errorf("get attachment [id=%d]: %w", attachmentID, err)
	}

	return mapAttachment(attachment), nil
}

func (r *AttachmentsRepo) SaveAttachment(ctx context.Context, attachment entity.Attachment) (entity.Attachment, error) {
	attachmentId, err := r.Queries.InsertAttachment(ctx, sqlc.InsertAttachmentParams{
		TripID: attachment.TripID,
		Name:   attachment.Name,
		Blob:   attachment.Blob,
	})
	if err != nil {
		return entity.Attachment{}, fmt.Errorf("insert attachment: %w", err)
	}

	return r.GetAttachmentByID(ctx, attachment.TripID, attachmentId)
}

func (r *AttachmentsRepo) DeleteAttachment(ctx context.Context, tripID int32, attachmentID int32) error {
	return r.Queries.DeleteAttachmentByID(ctx, sqlc.DeleteAttachmentByIDParams{TripID: tripID, ID: attachmentID})
}

func mapAttachments(attachments []sqlc.Attachment) []entity.Attachment {
	result := []entity.Attachment{}
	for _, attachment := range attachments {
		result = append(result, mapAttachment(attachment))
	}
	return result
}

func mapAttachment(attachment sqlc.Attachment) entity.Attachment {
	return entity.Attachment{
		ID:     attachment.ID,
		TripID: attachment.TripID,
		Name:   attachment.Name,
		Blob:   attachment.Blob,
	}
}
