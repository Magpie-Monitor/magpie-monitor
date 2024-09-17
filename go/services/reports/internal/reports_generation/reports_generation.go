package reportsgeneration

type ReportGenerator[T any] interface {
	Generate(logs []T)
}
