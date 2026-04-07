// Copyright (C) 2026 LeRedTeam
// SPDX-License-Identifier: AGPL-3.0-or-later

package parse

import (
	"testing"
)

func TestParseCloudTrail_SingleRecord(t *testing.T) {
	data := []byte(`{
		"Records": [{
			"eventSource": "s3.amazonaws.com",
			"eventName": "GetObject",
			"awsRegion": "us-east-1",
			"requestParameters": {
				"bucketName": "my-bucket",
				"key": "data.json"
			}
		}]
	}`)

	calls, err := ParseCloudTrail(data)

	if err != nil {
		t.Fatalf("ParseCloudTrail failed: %v", err)
	}

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}

	call := calls[0]
	if call.Service != "s3" {
		t.Errorf("expected service s3, got %s", call.Service)
	}
	if call.Action != "GetObject" {
		t.Errorf("expected action GetObject, got %s", call.Action)
	}
	if call.Region != "us-east-1" {
		t.Errorf("expected region us-east-1, got %s", call.Region)
	}
	if call.Resource != "arn:aws:s3:::my-bucket/data.json" {
		t.Errorf("expected resource arn:aws:s3:::my-bucket/data.json, got %s", call.Resource)
	}
}

func TestParseCloudTrail_MultipleRecords(t *testing.T) {
	data := []byte(`{
		"Records": [
			{
				"eventSource": "s3.amazonaws.com",
				"eventName": "GetObject",
				"awsRegion": "us-east-1",
				"requestParameters": {"bucketName": "bucket1", "key": "file1"}
			},
			{
				"eventSource": "dynamodb.amazonaws.com",
				"eventName": "PutItem",
				"awsRegion": "us-west-2",
				"requestParameters": {"tableName": "Users"}
			}
		]
	}`)

	calls, err := ParseCloudTrail(data)

	if err != nil {
		t.Fatalf("ParseCloudTrail failed: %v", err)
	}

	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}

	// Check S3 call
	if calls[0].Service != "s3" || calls[0].Action != "GetObject" {
		t.Errorf("first call should be s3:GetObject, got %s:%s", calls[0].Service, calls[0].Action)
	}

	// Check DynamoDB call
	if calls[1].Service != "dynamodb" || calls[1].Action != "PutItem" {
		t.Errorf("second call should be dynamodb:PutItem, got %s:%s", calls[1].Service, calls[1].Action)
	}
}

func TestParseCloudTrail_ArrayFormat(t *testing.T) {
	// Some CloudTrail exports are just arrays without the Records wrapper
	data := []byte(`[{
		"eventSource": "lambda.amazonaws.com",
		"eventName": "InvokeFunction",
		"awsRegion": "eu-west-1",
		"requestParameters": {"functionName": "MyFunc"}
	}]`)

	calls, err := ParseCloudTrail(data)

	if err != nil {
		t.Fatalf("ParseCloudTrail failed: %v", err)
	}

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}

	if calls[0].Service != "lambda" {
		t.Errorf("expected service lambda, got %s", calls[0].Service)
	}
}

func TestParseCloudTrail_WithResourceARN(t *testing.T) {
	data := []byte(`{
		"Records": [{
			"eventSource": "sqs.amazonaws.com",
			"eventName": "SendMessage",
			"awsRegion": "us-east-1",
			"resources": [
				{"ARN": "arn:aws:sqs:us-east-1:123456789:my-queue"}
			]
		}]
	}`)

	calls, err := ParseCloudTrail(data)

	if err != nil {
		t.Fatalf("ParseCloudTrail failed: %v", err)
	}

	if calls[0].Resource != "arn:aws:sqs:us-east-1:123456789:my-queue" {
		t.Errorf("expected ARN from resources, got %s", calls[0].Resource)
	}
}

func TestParseCloudTrail_InvalidJSON(t *testing.T) {
	data := []byte(`not valid json`)

	_, err := ParseCloudTrail(data)

	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseCloudTrail_Empty(t *testing.T) {
	data := []byte(`{"Records": []}`)

	calls, err := ParseCloudTrail(data)

	if err != nil {
		t.Fatalf("ParseCloudTrail failed: %v", err)
	}

	if len(calls) != 0 {
		t.Errorf("expected 0 calls, got %d", len(calls))
	}
}
