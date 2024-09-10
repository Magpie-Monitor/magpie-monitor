package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Incident struct {
	Category       string `bson:"category"`
	Summary        string `bson:"summary"`
	Recommendation string `bson:"recommendation"`
	Source         string `bson:"source"`
	Timestamp      int    `bson:"timestamp"`
}

type HostReport struct {
	Host         string      `bson:"host"`
	Precision    string      `bson:"precision"`
	CustomPrompt string      `bson:"customPrompt"`
	Incidents    []*Incident `bson:"incidents"`
}

type Report struct {
	Id          int           `bson:"id"`
	Title       string        `bson:"title"`
	StartMs     int           `bson:"startMs"`
	EndMs       int           `bson:"endMs"`
	HostReports []*HostReport `bson:"hostReports"`
}

var REPORTS_DB_NAME = "reports"
var REPORTS_COLLECTION = "reports"

type ReportRepository interface {
	GetAllReports(ctx context.Context) ([]*Report, error)
	InsertReport(ctx context.Context, report *Report) error
}

type MongoDbReportRepository struct {
	mongoDbClient *mongo.Client
	logger        *zap.Logger
}

func (r *MongoDbReportRepository) GetAllReports(ctx context.Context) ([]*Report, error) {
	return nil, nil
}

func (r *MongoDbReportRepository) InsertReport(ctx context.Context, report *Report) error {
	coll := r.mongoDbClient.Database(REPORTS_DB_NAME).Collection(REPORTS_COLLECTION)
	_, err := coll.InsertOne(context.TODO(), report)
	if err != nil {
		r.logger.Error("Failed to insert a report", zap.Error(err))
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
