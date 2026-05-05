package otpRepo

import (
	"context"
	"otp_service/internal/db"
)

type OtpTemplate interface {
	GetOtpTemplates(ctx context.Context, clientID, templateName string) ([]OtpTemplateModel, error)
}

type otpTemplates struct {
	db db.DB
}

func NewOtpTemplatesRepo(db db.DB) OtpTemplate {
	return &otpTemplates{db: db}
}

func (o *otpTemplates) GetOtpTemplates(ctx context.Context, clientID, templateName string) ([]OtpTemplateModel, error) {
	query := `SELECT id, template_name, template_body, vendor_template_id, vendor, channel, email_subject, is_active, created_at, updated_at
			  FROM otp_templates
			  WHERE client_id = $1 AND template_name = $2 AND is_active = TRUE`

	var templates []OtpTemplateModel
	err := o.db.SelectContext(ctx, &templates, query, clientID, templateName)
	if err != nil {
		return nil, err
	}
	return templates, nil
}
