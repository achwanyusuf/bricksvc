package model

import (
	"strings"
	"time"

	"github.com/achwanyusuf/bricksvc/src/domain/clientresponse"
	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	jsoniter "github.com/json-iterator/go"
	"github.com/lucsky/cuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	TransferKey                 string = "pubTransfer:%s"
	GetSingleByParamTransferKey string = "gspTransfer:%s"
	GetByParamTransferKey       string = "gpTransfer:%s"
	GetByParamTransferPgKey     string = "gppgTransfer:%s"
)

type TransferJob struct {
	ID      int    `json:"id"`
	JobID   string `json:"job_id"`
	APIKey  string `json:"api_key"`
	Payload string `json:"payload"`
	Status  string `json:"status"`
	BaseInformation
}

func TransformTransferJob(v entity.TransferJob) (TransferJob, error) {
	payload, err := v.Payload.MarshalJSON()
	if err != nil {
		return TransferJob{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}
	deletedBy := null.Int64{}
	if v.DeletedBy.Valid {
		deletedBy = null.Int64From(int64(v.DeletedBy.Int))
	}
	creationInfo := BaseInformation{
		CreatedBy: int64(v.CreatedBy),
		CreatedAt: v.CreatedAt,
		UpdatedBy: int64(v.UpdatedBy),
		UpdatedAt: v.UpdatedAt,
		DeletedBy: deletedBy.Int64,
		DeletedAt: v.DeletedAt.Time,
	}
	return TransferJob{
		ID:              v.ID,
		JobID:           v.JobID,
		APIKey:          v.APIKey,
		Payload:         string(payload),
		Status:          v.Status.String(),
		BaseInformation: creationInfo,
	}, nil
}

type CreateTransfer struct {
	SourceBankAccount      string    `json:"source_bank_account" query:"source_bank_account"`
	DestinationBankAccount string    `json:"destination_bank_account" query:"destination_bank_account"`
	SourceBankID           int64     `json:"source_bank_id" query:"source_bank_id"`
	DestinationBankID      int64     `json:"destination_bank_id" query:"destination_bank_id"`
	Amount                 float64   `json:"amount" query:"amount"`
	TransactionDate        time.Time `json:"transaction_time" query:"transaction_time"`
}

func (c *CreateTransfer) ToClientRes() clientresponse.Transfer {
	return clientresponse.Transfer{
		Amount:                 int(c.Amount),
		Status:                 entity.TransferstatusPending.String(),
		TransactionDate:        c.TransactionDate.Format(time.RFC3339),
		SourceBankAccount:      c.SourceBankAccount,
		DestinationBankAccount: c.DestinationBankAccount,
		SourceBankID:           int(c.SourceBankID),
		DestinationBankID:      int(c.DestinationBankID),
	}
}

func (c *CreateTransfer) ToEntity(apikey string) (entity.TransferJob, error) {
	payload, err := jsoniter.Marshal(c)
	if err != nil {
		return entity.TransferJob{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}
	return entity.TransferJob{
		JobID:   cuid.New(),
		APIKey:  apikey,
		Payload: payload,
		Status:  entity.TransferstatusPending,
	}, nil
}

type GetTransferJobByParam struct {
	ID          null.Int64  `json:"id" schema:"id" query:"id"`
	JobID       null.String `json:"job_id" schema:"job_id" query:"job_id"`
	APIKey      null.String `json:"api_key" schema:"api_key" query:"api_key"`
	Payload     null.String `json:"payload" schema:"payload" query:"payload"`
	Status      null.String `json:"status" schema:"status" query:"status"`
	CreatedAtGT null.Time   `json:"created_at_gt" schema:"created_at_gt" query:"created_at_gt"`
	CreatedAtLT null.Time   `json:"created_at_lt" schema:"created_at_lt" query:"created_at_lt"`
}

func (g *GetTransferJobByParam) GetQuery() []qm.QueryMod {
	var res []qm.QueryMod
	if g.ID.Valid {
		res = append(res, qm.Where("id=?", g.ID.Int64))
	}

	if g.JobID.Valid {
		res = append(res, qm.Where("job_id=?", g.JobID.String))
	}

	if g.APIKey.Valid {
		res = append(res, qm.Where("api_key=?", g.APIKey.String))
	}

	if g.Payload.Valid {
		res = append(res, qm.Where("payload=?", g.Payload.String))
	}

	if g.Status.Valid {
		res = append(res, qm.Where("status=?", g.Status.String))
	}

	if g.CreatedAtGT.Valid {
		res = append(res, qm.Where("created_at >", g.CreatedAtGT))
	}

	if g.CreatedAtLT.Valid {
		res = append(res, qm.Where("created_at <", g.CreatedAtLT))
	}
	return res
}

type GetTransferJobsByParam struct {
	GetTransferJobByParam
	OrderBy null.String `schema:"order_by" json:"order_by" query:"order_by"`
	Limit   int64       `schema:"limit" json:"limit" query:"limit"`
	Page    int64       `schema:"page" json:"page" query:"page"`
}

func (g *GetTransferJobsByParam) GetQuery() []qm.QueryMod {
	var res []qm.QueryMod
	if g.ID.Valid {
		res = append(res, qm.Where("id=?", g.ID.Int64))
	}

	if g.JobID.Valid {
		res = append(res, qm.Where("job_id=?", g.JobID.String))
	}

	if g.APIKey.Valid {
		res = append(res, qm.Where("api_key=?", g.APIKey.String))
	}

	if g.Payload.Valid {
		res = append(res, qm.Where("payload=?", g.Payload.String))
	}

	if g.Status.Valid {
		res = append(res, qm.Where("status=?", g.Status.String))
	}

	if g.OrderBy.Valid {
		order := strings.Split(g.OrderBy.String, ",")
		for _, o := range order {
			res = append(res, qm.OrderBy(o))
		}
	}

	return res
}

func TransformPSQLSingleTransferJob(v *entity.TransferJob) (TransferJob, error) {
	creationInfo := BaseInformation{
		CreatedBy: int64(v.CreatedBy),
		CreatedAt: v.CreatedAt,
		UpdatedBy: int64(v.UpdatedBy),
		UpdatedAt: v.UpdatedAt,
		DeletedBy: int64(v.DeletedBy.Int),
		DeletedAt: v.DeletedAt.Time,
	}

	payload, err := v.Payload.MarshalJSON()
	if err != nil {
		return TransferJob{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}

	return TransferJob{
		ID:              v.ID,
		JobID:           v.JobID,
		APIKey:          v.APIKey,
		Payload:         string(payload),
		Status:          v.Status.String(),
		BaseInformation: creationInfo,
	}, nil
}

func TransformPSQLTransferJob(role *entity.TransferJobSlice) ([]TransferJob, error) {
	var res []TransferJob
	for _, v := range *role {
		creationInfo := BaseInformation{
			CreatedBy: int64(v.CreatedBy),
			CreatedAt: v.CreatedAt,
			UpdatedBy: int64(v.UpdatedBy),
			UpdatedAt: v.UpdatedAt,
			DeletedBy: int64(v.DeletedBy.Int),
			DeletedAt: v.DeletedAt.Time,
		}

		payload, err := v.Payload.MarshalJSON()
		if err != nil {
			return []TransferJob{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
		}

		res = append(res, TransferJob{
			ID:              v.ID,
			JobID:           v.JobID,
			APIKey:          v.APIKey,
			Payload:         string(payload),
			Status:          v.Status.String(),
			BaseInformation: creationInfo,
		})
	}

	return res, nil
}

type GetTransferByParam struct {
	ID                     string `json:"id" query:"id"`
	Amount                 int    `json:"amoun" query:"amoun"`
	Status                 string `json:"status" query:"status"`
	TransactionDate        string `json:"transaction_date" query:"transaction_date"`
	SourceBankAccount      string `json:"source_bank_account" query:"source_bank_account"`
	DestinationBankAccount string `json:"destination_bank_account" query:"destination_bank_account"`
	SourceBankID           int    `json:"source_bank_id" query:"source_bank_id"`
	DestinationBankID      int    `json:"destination_bank_id" query:"destination_bank_id"`
}
