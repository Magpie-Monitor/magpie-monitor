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

type NodeIncident struct {
	Category       string `bson:"category"`
	Summary        string `bson:"summary"`
	Recommendation string `bson:"recommendation"`
	Source         string `bson:"source"`
	Timestamp      int64  `bson:"timestamp"`
}

type ApplicationIncident struct {
	Category       string `bson:"category"`
	Summary        string `bson:"summary"`
	Recommendation string `bson:"recommendation"`
	Source         string `bson:"source"`
	Timestamp      int64  `bson:"timestamp"`
	PodName        string `bson:"podName"`
	ContainerName  string `bson:"containerName"`
}

type NodeReport struct {
	Host         string         `bson:"host"`
	Precision    string         `bson:"precision"`
	CustomPrompt string         `bson:"customPrompt"`
	Incidents    []NodeIncident `bson:"incidents"`
}

type ApplicationReport struct {
	ApplicationName string                `bson:"name"`
	Precision       string                `bson:"precision"`
	CustomPrompt    string                `bson:"customPrompt"`
	Incidents       []ApplicationIncident `bson:"incidents"`
}

type ReportState string

const (
	ReportState_FailedToGenerate   ReportState = "failed_to_generate"
	ReportState_AwaitingGeneration ReportState = "awaiting_generation"
	ReportState_Generated          ReportState = "generated"
)

type Report struct {
	Id                      string              `bson:"_id,omitempty"`
	Status                  ReportState         `bson:"status"`
	RequestedAtNs           int64               `bson:"requestedAtMs"`
	GeneratedAtNs           int64               `bson:"generatedAtMs"`
	ScheduledGenerationAtMs int64               `bson:"scheduledGenerationAtMs"`
	Title                   string              `bson:"title"`
	FromDateNs              int64               `bson:"fromDateMs"`
	ToDateNs                int64               `bson:"toDateMs"`
	NodeReports             []NodeReport        `bson:"nodeReports"`
	ApplicationReports      []ApplicationReport `bson:"applicationReports"`
}

var REPORTS_DB_NAME = "reports"
var REPORTS_COLLECTION = "reports"

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

type ReportRepository interface {
	GetAllReports(ctx context.Context) ([]*Report, *ReportRepositoryError)
	InsertReport(ctx context.Context, report *Report) (*Report, *ReportRepositoryError)
	GetSingleReport(ctx context.Context, id string) (*Report, *ReportRepositoryError)
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

func (r *MongoDbReportRepository) GetAllReports(ctx context.Context) ([]*Report, *ReportRepositoryError) {
	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(REPORTS_COLLECTION)
	documents, err := coll.Find(ctx, bson.D{{}})
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
