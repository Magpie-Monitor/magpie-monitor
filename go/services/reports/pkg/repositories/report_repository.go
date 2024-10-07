package repositories

import (
	"context"

	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NodeReport struct {
	Node         string          `bson:"node" json:"node"`
	Precision    string          `bson:"precision" json:"precision"`
	CustomPrompt string          `bson:"customPrompt" json:"customPrompt"`
	Incidents    []*NodeIncident `bson:"-" json:"incidents"`
	IncidentIds  []string        `bson:"incidentIds" json:"-"`
}

type ApplicationReport struct {
	ApplicationName string                 `bson:"name" json:"applicationName"`
	Precision       string                 `bson:"precision" json:"precision"`
	CustomPrompt    string                 `bson:"customPrompt" json:"customPrompt"`
	Incidents       []*ApplicationIncident `bson:"-" json:"incidents"`
	IncidentIds     []string               `bson:"incidentIds" json:"-"`
}

type ReportState string

const (
	ReportState_FailedToGenerate   ReportState = "failed_to_generate"
	ReportState_AwaitingGeneration ReportState = "awaiting_generation"
	ReportState_Generated          ReportState = "generated"
)

type Urgency int

const (
	_ Urgency = iota
	Urgency_Low
	Urgency_Medium
	Urgency_High
)

type Report struct {
	Id                      string               `bson:"_id,omitempty" json:"id"`
	Status                  ReportState          `bson:"status" json:"status"`
	Cluster                 string               `bson:"cluster" json:"cluster"`
	FromDateNs              int64                `bson:"fromDateNs" json:"fromDateNs"`
	ToDateNs                int64                `bson:"toDateNs" json:"toDateNs"`
	RequestedAtNs           int64                `bson:"requestedAtNs" json:"requestedAtNs"`
	ScheduledGenerationAtMs int64                `bson:"scheduledGenerationAtNs" json:"scheduledGenerationAtNs"`
	Title                   string               `bson:"title" json:"title,omitempty"`
	NodeReports             []*NodeReport        `bson:"nodeReports" json:"nodeReports,omitempty"`
	ApplicationReports      []*ApplicationReport `bson:"applicationReports" json:"applicationReports,omitempty"`
	TotalApplicationEntries int                  `bson:"totalApplicationEntries" json:"totalApplicationEntries"`
	TotalNodeEntries        int                  `bson:"totalNodeEntries" json:"totalNodeEntries"`
	Urgency                 Urgency              `bson:"urgency" json:"urgency,omitempty"`

	// Present only if report is pending
	ScheduledApplicationInsights *ScheduledApplicationInsights `bson:"scheduledApplicationInsights" json:"scheduledApplicationInsights,omitempty"`
	ScheduledNodeInsights        *ScheduledNodeInsights        `bson:"scheduledNodeInsights" json:"scheduledNodeInsights,omitempty"`
}

type ApplicationInsightConfiguration struct {
	ApplicationName string `json:"applicationName"`
	Precision       string `json:"precision"`
	CustomPrompt    string `json:"customPrompt"`
}

type NodeInsightConfiguration struct {
	NodeName     string `json:"nodeName"`
	Precision    string `json:"precision"`
	CustomPrompt string `json:"customPrompt"`
}

type ScheduledApplicationInsights struct {
	Id                       string                             `json:"id"`
	CreatedAt                int64                              `json:"created_at"`
	ExpiresAt                int64                              `json:"expires_at"`
	CompletedAt              int64                              `json:"completed_at"`
	FromDateNs               int64                              `bson:"fromDateNs" json:"fromDateNs"`
	ToDateNs                 int64                              `bson:"toDateNs" json:"toDateNs"`
	Cluster                  string                             `bson:"cluster" json:"cluster"`
	ApplicationConfiguration []*ApplicationInsightConfiguration `json:"applicationConfiguration"`
}

type ScheduledNodeInsights struct {
	Id                string                      `json:"id"`
	CreatedAt         int64                       `json:"created_at"`
	ExpiresAt         int64                       `json:"expires_at"`
	CompletedAt       int64                       `json:"completed_at"`
	FromDateNs        int64                       `bson:"fromDateNs" json:"fromDateNs"`
	ToDateNs          int64                       `bson:"toDateNs" json:"toDateNs"`
	Cluster           string                      `bson:"cluster" json:"cluster"`
	NodeConfiguration []*NodeInsightConfiguration `json:"nodeConfiguration"`
}

func MapApplicationNameToConfiguration(configurations []*ApplicationInsightConfiguration) map[string]*ApplicationInsightConfiguration {
	groupedConfigurations := make(map[string]*ApplicationInsightConfiguration)
	for _, conf := range configurations {
		groupedConfigurations[conf.ApplicationName] = conf
	}

	return groupedConfigurations
}

func MapNodeNameToConfiguration(configurations []*NodeInsightConfiguration) map[string]*NodeInsightConfiguration {
	groupedConfigurations := make(map[string]*NodeInsightConfiguration)
	for _, conf := range configurations {
		groupedConfigurations[conf.NodeName] = conf
	}

	return groupedConfigurations
}

const (
	REPORTS_DB_NAME              = "reports"
	SCHEDULED_REPORTS_COLLECTION = "scheduled_reports"
	REPORTS_COLLECTION           = "reports"
)

type ReportRepositoryErrorKind string

const (
	ReportNotFound  ReportRepositoryErrorKind = "REPORT_NOT_FOUND"
	InvalidReportId ReportRepositoryErrorKind = "INVALID_REPORT_ID"
	InternalError   ReportRepositoryErrorKind = "INTERNAL_ERROR"
)

type ReportRepositoryError struct {
	msg  string
	kind ReportRepositoryErrorKind
}

func (e *ReportRepositoryError) Error() string {
	return e.msg
}

func (e *ReportRepositoryError) Kind() ReportRepositoryErrorKind {
	return e.kind
}

func NewReportNotFoundError(err error) *ReportRepositoryError {
	return &ReportRepositoryError{
		msg:  fmt.Sprintf("Report does not exists: %s", err),
		kind: ReportNotFound,
	}
}

func NewInvalidReportIdError(err error) *ReportRepositoryError {
	return &ReportRepositoryError{
		msg:  fmt.Sprintf("Invalid report id: %s", err),
		kind: InvalidReportId,
	}
}

func NewReportInternalError(err error) *ReportRepositoryError {
	return &ReportRepositoryError{
		msg:  err.Error(),
		kind: InternalError,
	}
}

type FilterParams struct {
	Cluster  *string
	FromDate *int64
	ToDate   *int64
}

type ReportRepository interface {
	GetAllReports(ctx context.Context, filter FilterParams) ([]*Report, *ReportRepositoryError)
	InsertReport(ctx context.Context, report *Report) (*Report, *ReportRepositoryError)
	GetSingleReport(ctx context.Context, id string) (*Report, *ReportRepositoryError)
	UpdateReport(ctx context.Context, report *Report) *ReportRepositoryError
	InsertApplicationIncidents(ctx context.Context, reports []*ApplicationReport) error
	InsertNodeIncidents(ctx context.Context, reports []*NodeReport) error
}

type MongoDbReportRepository struct {
	mongoDbClient                  *mongo.Client
	logger                         *zap.Logger
	applicationIncidentsRepository ApplicationIncidentRepository
	nodeIncidentsRepository        NodeIncidentRepository
}

func (r *MongoDbReportRepository) GetSingleReport(ctx context.Context, id string) (*Report, *ReportRepositoryError) {

	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(REPORTS_COLLECTION)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Failed to decode report id", zap.Error(err))
		return nil, NewInvalidReportIdError(err)
	}

	documents := coll.FindOne(ctx, bson.M{"_id": objectId})

	var result *Report
	err = documents.Decode(&result)
	if err != nil {
		r.logger.Error("Failed to decode reports from mongodb", zap.Error(err))
		return nil, NewReportNotFoundError(err)
	}

	for _, applicationReport := range result.ApplicationReports {
		incidents, err := r.applicationIncidentsRepository.GetIncidentsByIds(ctx, applicationReport.IncidentIds)
		if err != nil {
			r.logger.Error("Failed to join application incidents to report", zap.Error(err))
		}
		applicationReport.Incidents = incidents
	}

	for _, nodeReport := range result.NodeReports {
		incidents, err := r.nodeIncidentsRepository.GetIncidentsByIds(ctx, nodeReport.IncidentIds)
		if err != nil {
			r.logger.Error("Failed to join application incidents to report", zap.Error(err))
		}
		nodeReport.Incidents = incidents
	}

	return result, nil
}

func (r *MongoDbReportRepository) GetSingleSheduledReport(ctx context.Context, id string) (*Report, *ReportRepositoryError) {

	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(SCHEDULED_REPORTS_COLLECTION)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Failed to decode report id", zap.Error(err))
		return nil, NewInvalidReportIdError(err)
	}

	documents := coll.FindOne(ctx, bson.M{"_id": objectId})

	var result *Report
	err = documents.Decode(&result)
	if err != nil {
		r.logger.Error("Failed to decode reports from mongodb", zap.Error(err))
		return nil, NewReportNotFoundError(err)
	}

	return result, nil

}

func (r *MongoDbReportRepository) GetAllReports(ctx context.Context, filter FilterParams) ([]*Report, *ReportRepositoryError) {
	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(REPORTS_COLLECTION)

	mongoFilter := bson.M{}
	if filter.Cluster != nil {
		mongoFilter["cluster"] = *filter.Cluster
	}
	if filter.ToDate != nil {
		mongoFilter["toDateNs"] = *filter.ToDate
	}
	if filter.FromDate != nil {
		mongoFilter["fromDateNs"] = *filter.FromDate
	}

	documents, err := coll.Find(ctx, mongoFilter)
	if err != nil {
		r.logger.Error("Failed to get all reports from mongodb", zap.Error(err))
		return nil, NewReportInternalError(err)
	}

	var results []*Report
	err = documents.All(ctx, &results)
	if err != nil {
		r.logger.Error("Failed to decode reports from mongodb", zap.Error(err))
		return nil, NewReportInternalError(err)
	}

	return results, nil
}

func (r *MongoDbReportRepository) InsertApplicationIncidents(ctx context.Context, reports []*ApplicationReport) error {

	for _, report := range reports {
		ids, err := r.applicationIncidentsRepository.InsertIncidents(ctx, report.Incidents)
		if err != nil {
			r.logger.Error("Failed to insert incident incidents for a report", zap.Error(err))
			return err
		}
		report.IncidentIds = ids
	}

	return nil
}

func (r *MongoDbReportRepository) InsertNodeIncidents(ctx context.Context, reports []*NodeReport) error {

	for _, report := range reports {
		ids, err := r.nodeIncidentsRepository.InsertIncidents(ctx, report.Incidents)
		if err != nil {
			r.logger.Error("Failed to insert incident incidents for a report", zap.Error(err))
			return err
		}
		report.IncidentIds = ids
	}

	return nil
}

func (r *MongoDbReportRepository) InsertReport(ctx context.Context, report *Report) (*Report, *ReportRepositoryError) {
	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(REPORTS_COLLECTION)
	id, err := coll.InsertOne(context.TODO(), report)
	if err != nil {
		r.logger.Error("Failed to insert a report", zap.Error(err))
		return nil, NewReportInternalError(err)
	}

	var resultReport Report

	err = coll.FindOne(ctx, bson.M{"_id": id.InsertedID}).Decode(&resultReport)
	if err != nil {
		r.logger.Error("Failed to get inserted report", zap.Error(err))
		return nil, NewReportInternalError(err)
	}

	return &resultReport, nil
}

func (r *MongoDbReportRepository) UpdateReport(ctx context.Context, report *Report) *ReportRepositoryError {
	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(REPORTS_COLLECTION)

	hexId, err := primitive.ObjectIDFromHex(report.Id)
	if err != nil {
		r.logger.Error("Failed to get report id", zap.Error(err))
		return NewReportInternalError(err)
	}

	report.Id = ""

	_, err = coll.ReplaceOne(ctx, bson.D{{"_id", hexId}}, report)
	if err != nil {
		r.logger.Error("Failed to update a report", zap.Error(err))
		return NewReportInternalError(err)
	}
	return nil
}

type Params struct {
	fx.In
	ReportsDbMongoClient           *mongo.Client
	Logger                         *zap.Logger
	ApplicationIncidentsRepository ApplicationIncidentRepository
	NodeIncidentsRepository        NodeIncidentRepository
}

func NewMongoDbReportRepository(p Params) *MongoDbReportRepository {

	return &MongoDbReportRepository{
		mongoDbClient:                  p.ReportsDbMongoClient,
		logger:                         p.Logger,
		applicationIncidentsRepository: p.ApplicationIncidentsRepository,
		nodeIncidentsRepository:        p.NodeIncidentsRepository,
	}
}

func ProvideAsReportRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(ReportRepository)),
	)
}

// Compile-time check if MongoDbReportRepository implements
// the ReportRepository interface
var _ ReportRepository = &MongoDbReportRepository{}
