package lambdawrap

type contextKey string

const (
	contextKeyS3Entity = contextKey("S3_ENTITY")
	contextKeySNSARN   = contextKey("SNS_ARN")
	contextKeySQSARN   = contextKey("SQS_ARN")
)
