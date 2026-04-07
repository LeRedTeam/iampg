// Copyright (C) 2026 LeRedTeam
// SPDX-License-Identifier: AGPL-3.0-or-later

package capture

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/LeRedTeam/iampg/policy"
)

var (
	// Match AWS hostnames like s3.amazonaws.com, dynamodb.us-east-1.amazonaws.com
	awsHostRegex = regexp.MustCompile(`^([a-z0-9-]+)\.(?:([a-z0-9-]+)\.)?amazonaws\.com$`)
	// Match S3 bucket-style hosts like bucket.s3.amazonaws.com
	s3BucketRegex = regexp.MustCompile(`^(.+)\.s3\.([a-z0-9-]+\.)?amazonaws\.com$`)
)

// ParseAWSRequest extracts AWS API call information from an HTTP request.
func ParseAWSRequest(req *http.Request, body []byte) *policy.ObservedCall {
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}

	// Check if this is an AWS request
	if !strings.HasSuffix(host, ".amazonaws.com") && !strings.HasSuffix(host, ".amazonaws.com.cn") {
		return nil
	}

	call := &policy.ObservedCall{}

	// Extract service and region from host
	call.Service, call.Region = parseAWSHost(host)
	if call.Service == "" {
		return nil
	}

	// Extract action based on service type
	call.Action = parseAction(call.Service, req, body)
	if call.Action == "" {
		call.Action = "Unknown"
	}

	// Extract resource
	call.Resource = parseResource(call.Service, req, body, call.Region)

	return call
}

func parseAWSHost(host string) (service, region string) {
	// Handle S3 bucket-style URLs
	if matches := s3BucketRegex.FindStringSubmatch(host); matches != nil {
		return "s3", strings.TrimSuffix(matches[2], ".")
	}

	// Handle standard AWS URLs
	if matches := awsHostRegex.FindStringSubmatch(host); matches != nil {
		service = matches[1]
		region = matches[2]

		// Some services include region in first part
		if strings.HasPrefix(service, "s3-") || strings.HasPrefix(service, "s3.") {
			region = strings.TrimPrefix(strings.TrimPrefix(service, "s3-"), "s3.")
			service = "s3"
		}
		return service, region
	}

	return "", ""
}

func parseAction(service string, req *http.Request, body []byte) string {
	switch service {
	case "s3":
		return parseS3Action(req)
	case "dynamodb":
		return parseDynamoDBAction(req)
	case "lambda":
		return parseLambdaAction(req)
	case "sqs":
		return parseSQSAction(req, body)
	case "sns":
		return parseSNSAction(req, body)
	case "sts":
		return parseSTSAction(req, body)
	case "iam":
		return parseIAMAction(req, body)
	default:
		return parseGenericAction(req, body)
	}
}

func parseS3Action(req *http.Request) string {
	method := req.Method
	path := req.URL.Path
	query := req.URL.Query()

	// Check for specific query parameters that indicate actions
	if _, ok := query["versioning"]; ok {
		if method == "GET" {
			return "GetBucketVersioning"
		}
		return "PutBucketVersioning"
	}
	if _, ok := query["lifecycle"]; ok {
		if method == "GET" {
			return "GetLifecycleConfiguration"
		}
		return "PutLifecycleConfiguration"
	}
	if _, ok := query["policy"]; ok {
		if method == "GET" {
			return "GetBucketPolicy"
		}
		return "PutBucketPolicy"
	}
	if _, ok := query["acl"]; ok {
		if method == "GET" {
			return "GetObjectAcl"
		}
		return "PutObjectAcl"
	}
	if _, ok := query["uploads"]; ok {
		return "ListMultipartUploads"
	}

	// Basic operations based on method and path structure
	parts := strings.Split(strings.Trim(path, "/"), "/")
	hasKey := len(parts) > 1 && parts[1] != ""

	switch method {
	case "GET":
		if hasKey {
			return "GetObject"
		}
		return "ListBucket"
	case "PUT":
		if hasKey {
			return "PutObject"
		}
		return "CreateBucket"
	case "DELETE":
		if hasKey {
			return "DeleteObject"
		}
		return "DeleteBucket"
	case "HEAD":
		if hasKey {
			return "HeadObject"
		}
		return "HeadBucket"
	case "POST":
		return "PutObject" // Multipart upload
	}

	return ""
}

func parseDynamoDBAction(req *http.Request) string {
	// DynamoDB uses X-Amz-Target header
	target := req.Header.Get("X-Amz-Target")
	if target != "" {
		// Format: DynamoDB_20120810.GetItem
		parts := strings.Split(target, ".")
		if len(parts) == 2 {
			return parts[1]
		}
	}
	return ""
}

func parseLambdaAction(req *http.Request) string {
	path := req.URL.Path
	method := req.Method

	if strings.Contains(path, "/invocations") {
		return "InvokeFunction"
	}
	if strings.Contains(path, "/functions") {
		switch method {
		case "GET":
			return "GetFunction"
		case "POST":
			return "CreateFunction"
		case "DELETE":
			return "DeleteFunction"
		case "PUT":
			return "UpdateFunctionCode"
		}
	}
	return ""
}

func parseSQSAction(req *http.Request, body []byte) string {
	// SQS uses Action query parameter or form data
	if action := req.URL.Query().Get("Action"); action != "" {
		return action
	}
	return parseFormAction(body)
}

func parseSNSAction(req *http.Request, body []byte) string {
	if action := req.URL.Query().Get("Action"); action != "" {
		return action
	}
	return parseFormAction(body)
}

func parseSTSAction(req *http.Request, body []byte) string {
	if action := req.URL.Query().Get("Action"); action != "" {
		return action
	}
	return parseFormAction(body)
}

func parseIAMAction(req *http.Request, body []byte) string {
	if action := req.URL.Query().Get("Action"); action != "" {
		return action
	}
	return parseFormAction(body)
}

func parseGenericAction(req *http.Request, body []byte) string {
	// Try X-Amz-Target header first
	if target := req.Header.Get("X-Amz-Target"); target != "" {
		parts := strings.Split(target, ".")
		if len(parts) >= 2 {
			return parts[len(parts)-1]
		}
	}

	// Try Action query parameter
	if action := req.URL.Query().Get("Action"); action != "" {
		return action
	}

	// Try form data
	if action := parseFormAction(body); action != "" {
		return action
	}

	return ""
}

func parseFormAction(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return ""
	}
	return values.Get("Action")
}

func parseResource(service string, req *http.Request, body []byte, region string) string {
	switch service {
	case "s3":
		return parseS3Resource(req)
	case "dynamodb":
		return parseDynamoDBResource(body, region)
	case "lambda":
		return parseLambdaResource(req.URL.Path, region)
	case "sqs":
		return parseSQSResource(req, region)
	default:
		return "*"
	}
}

func parseS3Resource(req *http.Request) string {
	path := strings.Trim(req.URL.Path, "/")
	host := req.Host

	var bucket, key string

	// Check for bucket-style hostname
	if matches := s3BucketRegex.FindStringSubmatch(host); matches != nil {
		bucket = matches[1]
		key = path
	} else {
		// Path-style
		parts := strings.SplitN(path, "/", 2)
		if len(parts) > 0 {
			bucket = parts[0]
		}
		if len(parts) > 1 {
			key = parts[1]
		}
	}

	if bucket == "" {
		return "arn:aws:s3:::*"
	}
	if key == "" {
		return "arn:aws:s3:::" + bucket
	}
	return "arn:aws:s3:::" + bucket + "/" + key
}

func parseDynamoDBResource(body []byte, region string) string {
	if len(body) == 0 {
		return "*"
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "*"
	}

	if tableName, ok := data["TableName"].(string); ok {
		// We don't have account ID, so use wildcard
		return "arn:aws:dynamodb:" + region + ":*:table/" + tableName
	}
	return "*"
}

func parseLambdaResource(path, region string) string {
	// Path format: /2015-03-31/functions/FUNCTION_NAME/...
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if p == "functions" && i+1 < len(parts) {
			funcName := parts[i+1]
			return "arn:aws:lambda:" + region + ":*:function:" + funcName
		}
	}
	return "*"
}

func parseSQSResource(req *http.Request, region string) string {
	// SQS URL format: https://sqs.region.amazonaws.com/account/queue-name
	path := strings.Trim(req.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) >= 2 {
		account := parts[0]
		queueName := parts[1]
		return "arn:aws:sqs:" + region + ":" + account + ":" + queueName
	}
	return "*"
}
