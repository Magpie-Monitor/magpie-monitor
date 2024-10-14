package repositories

import (
	"context"

	"fmt"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
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
	Node         string          `bson:"node" json:"node"`
	Precision    string          `bson:"precision" json:"precision"`
	CustomPrompt string          `bson:"customPrompt" json:"customPrompt"`
	Incidents    []*NodeIncident `bson:"incidents" json:"incidents"`
}

type ApplicationReport struct {
	ApplicationName string                 `bson:"name" json:"applicationName"`
	Precision       string                 `bson:"precision" json:"precision"`
	CustomPrompt    string                 `bson:"customPrompt" json:"customPrompt"`
	Incidents       []*ApplicationIncident `bson:"incidents" json:"incidents"`
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
	ClusterId               string               `bson:"clusterId" json:"clusterId"`
	SinceNano               int64                `bson:"sinceNano" json:"sinceNano"`
	ToNano                  int64                `bson:"toNano" json:"toNano"`
	RequestedAtNs           int64                `bson:"requestedAtNs" json:"requestedAtNs"`
	ScheduledGenerationAtMs int64                `bson:"scheduledGenerationAtNs" json:"scheduledGenerationAtNs"`
	Title                   string               `bson:"title" json:"title"`
	NodeReports             []*NodeReport        `bson:"nodeReports" json:"nodeReports"`
	ApplicationReports      []*ApplicationReport `bson:"applicationReports" json:"applicationReports"`
	TotalApplicationEntries int                  `bson:"totalApplicationEntries" json:"totalApplicationEntries"`
	TotalNodeEntries        int                  `bson:"totalNodeEntries" json:"totalNodeEntries"`
	Urgency                 Urgency              `bson:"urgency" json:"urgency"`

	// Present only if report is pending
	ScheduledApplicationInsights *ScheduledApplicationInsights `bson:"scheduledApplicationInsights" json:"scheduledApplicationInsights"`
	ScheduledNodeInsights        *ScheduledNodeInsights        `bson:"scheduledNodeInsights" json:"scheduledNodeInsights"`
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
	SinceNano                int64                              `bson:"sinceNano" json:"sinceNano"`
	ToNano                   int64                              `bson:"toNano" json:"toNano"`
	ClusterId                string                             `bson:"clusterId" json:"clusterId"`
	ApplicationConfiguration []*ApplicationInsightConfiguration `json:"applicationConfiguration"`
}

type ScheduledNodeInsights struct {
	Id                string                      `json:"id"`
	SinceNano         int64                       `bson:"sinceNano" json:"sinceNano"`
	ToNano            int64                       `bson:"toNano" json:"toNano"`
	ClusterId         string                      `bson:"clusterId" json:"clusterId"`
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
	SinceNano *int64
	ToNano    *int64
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
	if filter.ToNano != nil {
		mongoFilter = append(mongoFilter,
			bson.D{
				{Key: "toNano", Value: bson.D{
					{Key: "$lte", Value: filter.ToNano}}},
			})
	}
	if filter.SinceNano != nil {
		mongoFilter = append(mongoFilter,
			bson.D{
				{Key: "sinceNano", Value: bson.D{
					{Key: "$gte", Value: filter.SinceNano}}},
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

func (r *MongoDbReportRepository) InsertApplicationIncidents(ctx context.Context, reports []*ApplicationReport) error {

	for _, report := range reports {
		ids, err := r.applicationIncidentsRepository.InsertIncidents(ctx, report.Incidents)
		if err != nil {
			r.logger.Error("Failed to insert incident incidents for a report", zap.Error(err))
			return err
		}

		insertedIncidents, err := r.applicationIncidentsRepository.GetIncidentsByIds(ctx, ids)
		if err != nil {
			r.logger.Error("Failed to retrieve inserted node incidents", zap.Error(err))
			return err
		}
		report.Incidents = insertedIncidents
	}

	return nil
}

func (r *MongoDbReportRepository) InsertNodeIncidents(ctx context.Context, reports []*NodeReport) error {

	for _, report := range reports {
		ids, err := r.nodeIncidentsRepository.InsertIncidents(ctx, report.Incidents)
		if err != nil {
			r.logger.Error("Failed to insert incidents from a report", zap.Error(err))
			return err
		}

		insertedIncidents, err := r.nodeIncidentsRepository.GetIncidentsByIds(ctx, ids)
		if err != nil {
			r.logger.Error("Failed to retrieve inserted node incidents", zap.Error(err))
			return err
		}
		report.Incidents = insertedIncidents
	}

	return nil
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
