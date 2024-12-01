package openai

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	scheduledjobs "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/scheduled_jobs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"time"
)

type OpenAiJobStatus = string

const OPENAI_DATABASE_NAME_KEY = "REPORTSDB_NAME"

const (
	OpenAiJobStatus__Enqueued   OpenAiJobStatus = "ENQUEUED"
	OpenAiJobStatus__InProgress OpenAiJobStatus = "IN_PROGRESS"
	OpenAiJobStatus__Completed  OpenAiJobStatus = "COMPLETED"
	OpenAiJobStatus__Failed     OpenAiJobStatus = "FAILED"
)

type OpenAiJob struct {
	Id                 string                        `bson:"_id,omitempty" json:"id"`
	ScheduledAt        int64                         `bson:"scheduledAt" json:"scheduledAt"`
	CompletionRequests map[string]*CompletionRequest `bson:"completionRequests" json:"completionRequests"`
	Status             OpenAiJobStatus               `bson:"status" json:"status"`
	BatchId            *string                       `bson:"batchId" json:"batchId"`
}

func (j *OpenAiJob) IsEqueued() bool {
	return j.Status == OpenAiJobStatus__Enqueued
}

func (j *OpenAiJob) IsCompleted() bool {
	return j.Status == OpenAiJobStatus__Completed
}

func (j *OpenAiJob) IsFailed() bool {
	return j.Status == OpenAiJobStatus__Failed
}

func (j *OpenAiJob) GetId() string {
	return j.Id
}

func (j *OpenAiJob) GetStatus() string {
	return j.Status
}

func (j *OpenAiJob) SetId(id string) {
	j.Id = id
}

func NewOpenAiJob(completionRequests map[string]*CompletionRequest) *OpenAiJob {
	return &OpenAiJob{
		ScheduledAt:        time.Now().UnixMilli(),
		Status:             OpenAiJobStatus__Enqueued,
		CompletionRequests: completionRequests,
	}
}

func NewOpenAiJobsCollection(log *zap.Logger, client *mongo.Client) *sharedrepositories.MongoDbCollection[*OpenAiJob] {
	envs.ValidateEnvs("Failed to get reports db name", []string{OPENAI_DATABASE_NAME_KEY})

	reportDb := os.Getenv(OPENAI_DATABASE_NAME_KEY)
	return &sharedrepositories.MongoDbCollection[*OpenAiJob]{Log: log, Db: reportDb, Col: "scheduled_openai_jobs", Client: client}
}

type OpenAiJobParams struct {
	fx.In
	OpenAiJobsDbMongoCollection *sharedrepositories.MongoDbCollection[*OpenAiJob]
	Logger                      *zap.Logger
}

func NewMongoDbOpenAiJobRepository(p OpenAiJobParams) *scheduledjobs.MongoDbScheduledJobRepository[*OpenAiJob] {

	return scheduledjobs.NewMongoDbScheduledJobRepository(p.OpenAiJobsDbMongoCollection, p.Logger)
}

func ProvideAsOpenAiJobRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(scheduledjobs.ScheduledJobRepository[*OpenAiJob])),
	)
}

// Compile-time check if MongoDbIncidentRepository implements
// the IncidentRepository interface
var _ scheduledjobs.ScheduledJobRepository[*OpenAiJob] = &scheduledjobs.MongoDbScheduledJobRepository[*OpenAiJob]{}
