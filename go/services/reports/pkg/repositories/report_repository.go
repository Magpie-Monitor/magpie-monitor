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

type NodeIncidentSource struct {
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
	NodeName  string `bson:"nodeName" json:"nodeName"`
	Content   string `bson:"content" json:"content"`
}

type NodeIncident struct {
	Id             string               `bson:"_id,omitempty" json:"id"`
	Category       string               `bson:"category" json:"category"`
	Summary        string               `bson:"summary" json:"summary"`
	Recommendation string               `bson:"recommendation" json:"recommendation"`
	Urgency        Urgency              `bson:"urgency" json:"urgency"`
	Sources        []NodeIncidentSource `bson:"sources" json:"sources"`
}

type ApplicationIncidentSource struct {
	Timestamp     int64  `bson:"timestamp" json:"timestamp"`
	PodName       string `bson:"podName" json:"podName"`
	ContainerName string `bson:"containerName" json:"containerName"`
	Image         string `bson:"image" json:"image"`
	Content       string `bson:"content" json:"content"`
}

type ApplicationIncident struct {
	Id             string                      `bson:"_id,omitempty" json:"id"`
	Category       string                      `bson:"category" json:"category"`
	Summary        string                      `bson:"summary" json:"summary"`
	Recommendation string                      `bson:"recommendation" json:"recommendation"`
	Urgency        Urgency                     `bson:"urgency" json:"urgency"`
	Sources        []ApplicationIncidentSource `bson:"sources" json:"sources"`
}

type NodeReport struct {
	Node         string         `bson:"node" json:"node"`
	Precision    string         `bson:"precision" json:"precision"`
	CustomPrompt string         `bson:"customPrompt" json:"customPrompt"`
	Incidents    []NodeIncident `bson:"incidents" json:"incidents"`
}

type ApplicationReport struct {
	ApplicationName string                `bson:"name" json:"applicationName"`
	Precision       string                `bson:"precision" json:"precision"`
	CustomPrompt    string                `bson:"customPrompt" json:"customPrompt"`
	Incidents       []ApplicationIncident `bson:"incidents" json:"incidents"`
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
	Id                      string              `bson:"_id,omitempty" json:"id"`
	Status                  ReportState         `bson:"status" json:"status"`
	Cluster                 string              `bson:"cluster" json:"cluster"`
	FromDateNs              int64               `bson:"fromDateNs" json:"fromDateNs"`
	ToDateNs                int64               `bson:"toDateNs" json:"toDateNs"`
	RequestedAtNs           int64               `bson:"requestedAtNs" json:"requestedAtNs"`
	ScheduledGenerationAtMs int64               `bson:"scheduledGenerationAtNs" json:"scheduledGenerationAtNs"`
	Title                   string              `bson:"title" json:"title,omitempty"`
	NodeReports             []NodeReport        `bson:"nodeReports" json:"nodeReports,omitempty"`
	ApplicationReports      []ApplicationReport `bson:"applicationReports" json:"applicationReports,omitempty"`
	TotalApplicationEntries int                 `bson:"totalApplicationEntries" json:"totalApplicationEntries"`
	TotalNodeEntries        int                 `bson:"totalNodeEntries" json:"totalNodeEntries"`
	Urgency                 Urgency             `bson:"urgency" json:"urgency,omitempty"`

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
}

type MongoDbReportRepository struct {
	mongoDbClient *mongo.Client
	logger        *zap.Logger
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
	ReportsDbMongoClient *mongo.Client
	Logger               *zap.Logger
}

func NewMongoDbReportRepository(p Params) *MongoDbReportRepository {

	return &MongoDbReportRepository{
		mongoDbClient: p.ReportsDbMongoClient,
		logger:        p.Logger,
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
