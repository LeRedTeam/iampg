package parse

import (
	"testing"
)

func TestParseAccessDenied_FullFormat(t *testing.T) {
	msg := "User: arn:aws:iam::123456789:user/dev is not authorized to perform: s3:GetObject on resource: arn:aws:s3:::my-bucket/secret.txt"

	calls := ParseAccessDenied(msg)

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
	if call.Resource != "arn:aws:s3:::my-bucket/secret.txt" {
		t.Errorf("expected resource arn:aws:s3:::my-bucket/secret.txt, got %s", call.Resource)
	}
}

func TestParseAccessDenied_DynamoDB(t *testing.T) {
	msg := "User: arn:aws:iam::123:user/dev is not authorized to perform: dynamodb:PutItem on resource: arn:aws:dynamodb:us-east-1:123:table/Users"

	calls := ParseAccessDenied(msg)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}

	call := calls[0]
	if call.Service != "dynamodb" {
		t.Errorf("expected service dynamodb, got %s", call.Service)
	}
	if call.Action != "PutItem" {
		t.Errorf("expected action PutItem, got %s", call.Action)
	}
}

func TestParseAccessDenied_Lambda(t *testing.T) {
	msg := "not authorized to perform: lambda:InvokeFunction on resource: arn:aws:lambda:us-east-1:123:function:MyFunc"

	calls := ParseAccessDenied(msg)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}

	call := calls[0]
	if call.Service != "lambda" {
		t.Errorf("expected service lambda, got %s", call.Service)
	}
	if call.Action != "InvokeFunction" {
		t.Errorf("expected action InvokeFunction, got %s", call.Action)
	}
	if call.Resource != "arn:aws:lambda:us-east-1:123:function:MyFunc" {
		t.Errorf("unexpected resource: %s", call.Resource)
	}
}

func TestParseAccessDenied_ShortFormat(t *testing.T) {
	msg := "AccessDenied: sqs:SendMessage on arn:aws:sqs:us-east-1:123:my-queue"

	calls := ParseAccessDenied(msg)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}

	call := calls[0]
	if call.Service != "sqs" {
		t.Errorf("expected service sqs, got %s", call.Service)
	}
	if call.Action != "SendMessage" {
		t.Errorf("expected action SendMessage, got %s", call.Action)
	}
}

func TestParseMultipleErrors(t *testing.T) {
	input := `not authorized to perform: s3:GetObject on resource: arn:aws:s3:::bucket/key1
not authorized to perform: s3:PutObject on resource: arn:aws:s3:::bucket/key2
not authorized to perform: dynamodb:Query on resource: arn:aws:dynamodb:us-east-1:123:table/Users`

	calls := ParseMultipleErrors(input)

	if len(calls) != 3 {
		t.Fatalf("expected 3 calls, got %d", len(calls))
	}
}

func TestParseAccessDenied_Empty(t *testing.T) {
	calls := ParseAccessDenied("")

	if len(calls) != 0 {
		t.Errorf("expected 0 calls for empty input, got %d", len(calls))
	}
}

func TestParseAccessDenied_NoMatch(t *testing.T) {
	calls := ParseAccessDenied("some random error message")

	if len(calls) != 0 {
		t.Errorf("expected 0 calls for non-matching input, got %d", len(calls))
	}
}
