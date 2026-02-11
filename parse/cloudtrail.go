package parse

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/LeRedTeam/iampg/policy"
)

// CloudTrailLog represents a CloudTrail log file structure.
type CloudTrailLog struct {
	Records []CloudTrailRecord `json:"Records"`
}

// CloudTrailRecord represents a single CloudTrail event.
type CloudTrailRecord struct {
	EventSource string              `json:"eventSource"`
	EventName   string              `json:"eventName"`
	AWSRegion   string              `json:"awsRegion"`
	Resources   []CloudTrailResource `json:"resources"`
	RequestParameters json.RawMessage `json:"requestParameters"`
	ErrorCode   string              `json:"errorCode"`
}

// CloudTrailResource represents a resource in a CloudTrail event.
type CloudTrailResource struct {
	ARN  string `json:"ARN"`
	Type string `json:"type"`
}

// ParseCloudTrail parses CloudTrail JSON and returns observed calls.
func ParseCloudTrail(data []byte) ([]policy.ObservedCall, error) {
	var log CloudTrailLog
	if err := json.Unmarshal(data, &log); err != nil {
		// Try parsing as array directly
		var records []CloudTrailRecord
		if err2 := json.Unmarshal(data, &records); err2 != nil {
			return nil, fmt.Errorf("invalid CloudTrail format: %w", err)
		}
		log.Records = records
	}

	var calls []policy.ObservedCall
	for _, record := range log.Records {
		call := recordToCall(record)
		if call != nil {
			calls = append(calls, *call)
		}
	}

	return calls, nil
}

func recordToCall(record CloudTrailRecord) *policy.ObservedCall {
	// Extract service from eventSource (e.g., "s3.amazonaws.com" -> "s3")
	service := strings.TrimSuffix(record.EventSource, ".amazonaws.com")

	// Get resource ARN
	resource := "*"
	if len(record.Resources) > 0 && record.Resources[0].ARN != "" {
		resource = record.Resources[0].ARN
	} else {
		// Try to extract from request parameters
		resource = extractResourceFromParams(service, record.RequestParameters)
	}

	return &policy.ObservedCall{
		Service:  service,
		Action:   record.EventName,
		Resource: resource,
		Region:   record.AWSRegion,
	}
}

func extractResourceFromParams(service string, params json.RawMessage) string {
	if len(params) == 0 {
		return "*"
	}

	var data map[string]interface{}
	if err := json.Unmarshal(params, &data); err != nil {
		return "*"
	}

	switch service {
	case "s3":
		bucket, _ := data["bucketName"].(string)
		key, _ := data["key"].(string)
		if bucket != "" {
			if key != "" {
				return "arn:aws:s3:::" + bucket + "/" + key
			}
			return "arn:aws:s3:::" + bucket
		}
	case "dynamodb":
		if table, ok := data["tableName"].(string); ok {
			return "arn:aws:dynamodb:*:*:table/" + table
		}
	case "lambda":
		if fn, ok := data["functionName"].(string); ok {
			return "arn:aws:lambda:*:*:function:" + fn
		}
	case "sqs":
		if url, ok := data["queueUrl"].(string); ok {
			// Parse queue URL to get queue name
			parts := strings.Split(url, "/")
			if len(parts) >= 2 {
				return "arn:aws:sqs:*:" + parts[len(parts)-2] + ":" + parts[len(parts)-1]
			}
		}
	}

	return "*"
}
