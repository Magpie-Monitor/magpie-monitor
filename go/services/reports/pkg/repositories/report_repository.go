package repositories

import (
	"context"
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewReportCollection(log *zap.Logger, client *mongo.Client) *repositories.MongoDbCollection[*Report] {
	return &repositories.MongoDbCollection[*Report]{Log: log, Db: "reports", Col: "reports", Client: client}
}

type NodeReport struct {
	Node         string            `bson:"node" json:"node"`
	Accuracy     insights.Accuracy `bson:"accuracy" json:"accuracy"`
	CustomPrompt string            `bson:"customPrompt" json:"customPrompt"`
	Incidents    []*NodeIncident   `bson:"incidents" json:"incidents"`
}

type ApplicationReport struct {
	ApplicationName string                 `bson:"name" json:"applicationName"`
	Accuracy        insights.Accuracy      `bson:"accuracy" json:"accuracy"`
	CustomPrompt    string                 `bson:"customPrompt" json:"customPrompt"`
	Incidents       []*ApplicationIncident `bson:"incidents" json:"incidents"`
}

type ReportState string

const (
	ReportState_FailedToGenerate        ReportState = "failed_to_generate"
	ReportState_AwaitingGeneration      ReportState = "awaiting_generation"
	ReportState_AwaitingIncidentMerging ReportState = "awaiting_incident_merging"
	ReportState_Generated               ReportState = "generated"
)

type ScheduledIncidentMergerJob struct {
	Id string
}

type Report struct {
	Id                      string               `bson:"_id,omitempty" json:"id"`
	CorrelationId           string               `bson:"correlationId" json:"correlationId"`
	Status                  ReportState          `bson:"status" json:"status"`
	ClusterId               string               `bson:"clusterId" json:"clusterId"`
	SinceMs                 int64                `bson:"sinceMs" json:"sinceMs"`
	ToMs                    int64                `bson:"toMs" json:"toMs"`
	RequestedAtMs           int64                `bson:"requestedAtMs" json:"requestedAtMs"`
	ScheduledGenerationAtMs int64                `bson:"scheduledGenerationAtMs" json:"scheduledGenerationAtMs"`
	Title                   string               `bson:"title" json:"title"`
	NodeReports             []*NodeReport        `bson:"nodeReports" json:"nodeReports"`
	ApplicationReports      []*ApplicationReport `bson:"applicationReports" json:"applicationReports"`
	TotalApplicationEntries int                  `bson:"totalApplicationEntries" json:"totalApplicationEntries"`
	TotalNodeEntries        int                  `bson:"totalNodeEntries" json:"totalNodeEntries"`
	Urgency                 insights.Urgency     `bson:"urgency" json:"urgency"`

	// Filled only when report is scheduled
	ScheduledApplicationInsights *insights.ScheduledApplicationInsights `bson:"scheduledApplicationInsights" json:"scheduledApplicationInsights"`
	ScheduledNodeInsights        *insights.ScheduledNodeInsights        `bson:"scheduledNodeInsights" json:"scheduledNodeInsights"`
	AnalyzedApplications         int                                    `bson:"analyzedApplications" json:"analyzedApplications"`
	AnalyzedNodes                int                                    `bson:"analyzedNodes" json:"analyzedNodes"`

	// Filled only when report is generated and awaiting summarization
	ScheduledApplicationIncidentMergerJobs []*ScheduledIncidentMergerJob `json:"-"`
	ScheduledNodeIncidentMergerJobs        []*ScheduledIncidentMergerJob `json:"-"`
}

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
	ClusterId *string
	SinceMs   *int64
	ToMs      *int64
}

type ReportRepository interface {
	GetAllReports(ctx context.Context, filter FilterParams) ([]*Report, *ReportRepositoryError)
	InsertReport(ctx context.Context, report *Report) (*Report, *ReportRepositoryError)
	GetSingleReport(ctx context.Context, id string) (*Report, *ReportRepositoryError)
	UpdateReport(ctx context.Context, report *Report) *ReportRepositoryError
	InsertApplicationIncidents(ctx context.Context, incidents []*ApplicationIncident) ([]*ApplicationIncident, error)
	InsertNodeIncidents(ctx context.Context, reports []*NodeIncident) ([]*NodeIncident, error)
	GetPendingGenerationReports(ctx context.Context) ([]*Report, error)
	GetPendingIncidentMergingReports(ctx context.Context) ([]*Report, error)
}

func (r *MongoDbReportRepository) GetPendingGenerationReports(ctx context.Context) ([]*Report, error) {

	result, err := r.mongoDbCollection.GetDocuments(primitive.D{{Key: "status", Value: ReportState_AwaitingGeneration}}, primitive.D{})
	if err != nil {
		r.logger.Error("Failed to get all pending reports", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *MongoDbReportRepository) GetPendingIncidentMergingReports(ctx context.Context) ([]*Report, error) {

	result, err := r.mongoDbCollection.GetDocuments(primitive.D{{Key: "status", Value: ReportState_AwaitingIncidentMerging}}, primitive.D{})
	if err != nil {
		r.logger.Error("Failed to get all reports awaiting incident merging", zap.Error(err))
		return nil, err
	}

	return result, nil
}

type MongoDbReportRepository struct {
	logger                         *zap.Logger
	applicationIncidentsRepository IncidentRepository[ApplicationIncident]
	nodeIncidentsRepository        IncidentRepository[NodeIncident]
	mongoDbCollection              *repositories.MongoDbCollection[*Report]
}

func (r *MongoDbReportRepository) GetSingleReport(ctx context.Context, id string) (*Report, *ReportRepositoryError) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Failed to decode report id", zap.Error(err))
		return nil, NewInvalidReportIdError(err)
	}

	result, err := r.mongoDbCollection.GetDocument(bson.D{{Key: "_id", Value: objectId}}, bson.D{})
	if err != nil {
		r.logger.Error("Failed to decode reports from mongodb", zap.Error(err))
		return nil, NewReportNotFoundError(err)
	}

	return result, nil
}

func (r *MongoDbReportRepository) GetAllReports(ctx context.Context, filter FilterParams) ([]*Report, *ReportRepositoryError) {

	mongoFilter := bson.A{}
	if filter.ClusterId != nil {
		mongoFilter = append(mongoFilter,
			bson.D{
				{Key: "clusterId", Value: bson.D{
					{Key: "$eq", Value: filter.ClusterId}}},
			})
	}
	if filter.ToMs != nil {
		mongoFilter = append(mongoFilter,
			bson.D{
				{Key: "toMs", Value: bson.D{
					{Key: "$lte", Value: filter.ToMs}}},
			})
	}
	if filter.SinceMs != nil {
		mongoFilter = append(mongoFilter,
			bson.D{
				{Key: "sinceMs", Value: bson.D{
					{Key: "$gte", Value: filter.SinceMs}}},
			})
	}

	finalFilter := bson.D{}

	if len(mongoFilter) > 0 {
		finalFilter = append(finalFilter, bson.E{
			Key: "$and", Value: mongoFilter,
		})
	}

	reports, err := r.mongoDbCollection.GetDocuments(finalFilter, primitive.D{})

	if err != nil {
		r.logger.Error("Failed to get all reports from mongodb", zap.Error(err))
		return nil, NewReportInternalError(err)
	}

	return reports, nil
}

func (r *MongoDbReportRepository) InsertApplicationIncidents(ctx context.Context, incidents []*ApplicationIncident) ([]*ApplicationIncident, error) {

	ids, err := r.applicationIncidentsRepository.InsertIncidents(ctx, incidents)
	if err != nil {
		r.logger.Error("Failed to insert incident incidents for a report", zap.Error(err))
		return nil, err
	}

	insertedIncidents, err := r.applicationIncidentsRepository.GetIncidentsByIds(ctx, ids)
	if err != nil {
		r.logger.Error("Failed to retrieve inserted node incidents", zap.Error(err))
		return nil, err
	}

	return insertedIncidents, nil
}

func (r *MongoDbReportRepository) InsertNodeIncidents(ctx context.Context, incidents []*NodeIncident) ([]*NodeIncident, error) {

	ids, err := r.nodeIncidentsRepository.InsertIncidents(ctx, incidents)
	if err != nil {
		r.logger.Error("Failed to insert incidents from a report", zap.Error(err))
		return nil, err
	}

	insertedIncidents, err := r.nodeIncidentsRepository.GetIncidentsByIds(ctx, ids)
	if err != nil {
		r.logger.Error("Failed to retrieve inserted node incidents", zap.Error(err))
		return nil, err
	}

	return insertedIncidents, nil
}

func (r *MongoDbReportRepository) InsertReport(ctx context.Context, report *Report) (*Report, *ReportRepositoryError) {

	id, err := r.mongoDbCollection.InsertDocument(report)
	if err != nil {
		r.logger.Error("Failed to insert a report", zap.Error(err))
		return nil, NewReportInternalError(err)
	}

	resultReport, repErr := r.GetSingleReport(ctx, id.Hex())
	if repErr != nil {
		r.logger.Error("Failed to insert a report", zap.Error(repErr))
		return nil, NewReportInternalError(repErr)
	}

	return resultReport, nil
}

func (r *MongoDbReportRepository) UpdateReport(ctx context.Context, report *Report) *ReportRepositoryError {

	id, err := primitive.ObjectIDFromHex(report.Id)
	if err != nil {
		r.logger.Error("Failed to encode report id", zap.Error(err))
		return NewInvalidReportIdError(err)
	}
	report.Id = ""

	err = r.mongoDbCollection.ReplaceDocument(ctx, id, report)
	if err != nil {
		r.logger.Error("Failed to get report id", zap.Error(err))
		return NewReportInternalError(err)
	}

	report.Id = id.Hex()

	return nil
}

type Params struct {
	fx.In
	ReportsDbMongoColl             *repositories.MongoDbCollection[*Report]
	Logger                         *zap.Logger
	ApplicationIncidentsRepository IncidentRepository[ApplicationIncident]
	NodeIncidentsRepository        IncidentRepository[NodeIncident]
}

func NewMongoDbReportRepository(p Params) *MongoDbReportRepository {

	return &MongoDbReportRepository{
		mongoDbCollection:              p.ReportsDbMongoColl,
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
