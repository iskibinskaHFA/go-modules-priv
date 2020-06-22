package migrate

import (
	"github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
	"math/big"
	"time"
)

// Header is gorm model for usage.headers
type Header struct {
	HeaderID       []byte    `gorm:"primary_key" json:"header_id"`
	HeaderStatusID int       `json:"header_status_id"`
	SenderHeaderID string    `json:"sender_header_id"`
	SenderID       string    `json:"sender_id"`
	SenderName     string    `json:"sender_name"`
	PeriodStart    time.Time `gorm:"type:DATE"`
	PeriodEnd      time.Time `gorm:"type:DATE"`
	*royModel
}

// UsageSummary is gorm model for usage.usage_summaries
type UsageSummary struct {
	UsageSummaryID       []byte         `gorm:"primary_key" json:"usage_summary_id"`
	SenderUsageSummaryID string         `json:"sender_usage_summary_id"`
	HeaderID             []byte         `json:"header_id"`
	ServiceID            string         `json:"service_id"`
	SalesData            postgres.Jsonb `json:"sales_data"`
	*royModel
}
type Usages struct {
	UsageID                  []byte  `gorm:"primary_key" json:"usage_id"`
	HeaderID                 []byte  `json:"sender_usage_summary_id"`
	ResourceID               []byte  `json: "resource_id"`
	UsageSummaryID           []byte  `json: "usage_summary_id"`
	Units                    int64   `json: "units"`
	DurationAdjustmentFactor float64 `json:"duration_adjustment_factor"`
	AdjustedUnits            float64 `json:"adjusted_units"`
	*royModel
}

//Revenue structure is a helper structure for base rate calculation
type Revenue struct {
	SubscriberCount      int64   `json: "subscriber_count"`
	NetServiceRevenue    float64 `json: "net_service_revenue"`
	LabelContentCost     float64 `json: "label_content_cost"`
	PerformanceRoyalties float64 `json: "performance_royalties"`
	AdjustedUnitsTotal   float64
}

// RateResult holds the result of calculation
type RateResult struct {
	DefaultRate     float64 `json: "default_rate"`
	GreaterThanFive float64 `json: "greater_than_five"`
}

//type cons struct {
//	Floor  int `json: "floor"`
//	Record int `json: "record"`
//	Rev    int `json: "rev"`
//}

// Service is gorm model for royalty.services
type Service struct {
	ServiceID   string `gorm:"primary_key;type:VARCHAR" json:"service_id" sql:"VARCHAR"`
	Name        string `json:"name"`
	Description string `json:"description"`
	*royModel
}

// RateFormula is gorm model for royalty.rate_formulas
type RateFormula struct {
	RateFormulaID   uint           `gorm:"primary_key;AUTO_INCREMENT" json:"rate_formula_id"`
	Formula         postgres.Jsonb `json:"formula"`
	FormulaMetadata postgres.Jsonb `json:"formula_metadata"`
	RateDefinitions []RateDefinition
	*royModel
}

// RateDefinition is gorm model for royalty.rate_formulas
type RateDefinition struct {
	RateDefinitionID uint `gorm:"primary_key;AUTO_INCREMENT" json:"rate_definition_id"`
	RateFormulaID    uint `gorm:"not null"`
	RateFormula      RateFormula
	ServiceID        string         `json:"service_id"`
	StartDate        time.Time      `json:"start_date"`
	EndDate          time.Time      `json:"start_end"`
	Constants        postgres.Jsonb `json:"constants"`
	*royModel
}

// CalcStepsLogDefinition is gorm model for royalty.calc_steps_logs_definitions
type CalcStepsLogDefinition struct {
	LogDefinitionID uint   `gorm:"primary_key";AUTO_INCREMENT" json:"log_definition_id"`
	ServiceID       string `gorm:"not null"`
	Params          string
	Sprintf         string
	Step            string
	Result          string
	SequenceID      uint `gorm:"not null"`
	*royModel
}

// CalcStepsLog is gorm model for royalty.calc_steps_logs
type CalcStepsLog struct {
	CalcStepsLogID  uint   `gorm:"primary_key";AUTO_INCREMENT" json:"calc_step_log_id"`
	LogDefinitionID uint   `gorm:"not null"`
	UsageSummaryID  []byte `gorm:"not null"`
	SequenceID      uint   `gorm:"not null"`
	ResultValue     float64
	Text            string
	Step            string
	*royModel
}

// Rate is gorm model for royalty.rates which holds base rates
type Rate struct {
	RateID         uint           `gorm:"primary_key" json:"rate_id"`
	UsageSummaryID []byte         `json:"usage_summary_id"`
	AmountPerUnit  postgres.Jsonb `json:"amount_per_unit"` // 2 rates, > 5 min
	*royModel
}

// Release is gorm model for usage.releases
type Release struct {
	ReleaseID       []byte    `gorm:"primary_key" json:"release_id"`
	ReleaseUUID     uuid.UUID `json:"release_uuid; default:uuid_generate_v4()"`
	SenderReleaseID string    `json:"sender_release_id"`
	Title           string    `json:"title"`
	LabelName       string    `json:"label_name"`
	Upc             string    `json:"upc"`
	*royModel
}

// Work is gorm model for usage.works
type Work struct {
	WorkID       []byte `gorm:"primary_key" json:"work_id"`
	SenderWorkID string `json:"sender_work_id"`
	Iswc         string `json:"iswc"`
	*royModel
}

// Writer is gorm model for usage.writers
type Writer struct {
	PartyID   string     `gorm:"primary_key" json:"party_id" sql:"VARCHAR"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy string     `json:"updated_by"`
}

// Resource is gorm model for usage.resources
type Resource struct {
	ResourceID               []byte    `gorm:"primary_key" json:"resource_id"`
	SenderResourceID         string    `json:"sender_resource_id"`
	OriginID                 []byte    `json:"origin_id"`
	WorkID                   []byte    `json:"work_id"`
	ProcessingStatusID       big.Int   `json:"processing_status_id"`
	HfaSongCode              string    `json:"hfa_song_code"`
	ServerFixationDate       time.Time `json:"server_fixation_date"`
	Title                    string    `json:"title"`
	Artist                   string    `json:"artist"`
	Isrc                     string    `json:"isrc"`
	PlayMinutes              uint      `json:"play_minutes"`
	PlaySeconds              uint      `json:"play_seconds"`
	EffectiveDuration        uint      `json:"effective_duration"`
	DurationAdjustmentFactor float64   `json:"duration_adjustment_factor"`
	ReleaseDate              time.Time `json:"release_date"`
	*royModel
}

// WorksWriter is gorm model for royalty.works_writers
type WorksWriter struct {
	WorkID  []byte `gorm:"primary_key" json:"work_id"`
	PartyID string `gorm:"primary_key; json:"party_id"`
	Iswc    string `json:"iswc"`
	*royModel
}

type royModel struct {
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}
