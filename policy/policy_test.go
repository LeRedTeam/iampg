package policy

import (
	"encoding/json"
	"testing"
)

func TestGenerate_Empty(t *testing.T) {
	doc := Generate([]ObservedCall{})

	if doc.Version != "2012-10-17" {
		t.Errorf("expected version 2012-10-17, got %s", doc.Version)
	}
	if len(doc.Statement) != 0 {
		t.Errorf("expected 0 statements, got %d", len(doc.Statement))
	}
}

func TestGenerate_SingleCall(t *testing.T) {
	calls := []ObservedCall{
		{Service: "s3", Action: "GetObject", Resource: "arn:aws:s3:::bucket/key"},
	}

	doc := Generate(calls)

	if len(doc.Statement) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(doc.Statement))
	}

	stmt := doc.Statement[0]
	if stmt.Effect != "Allow" {
		t.Errorf("expected effect Allow, got %s", stmt.Effect)
	}
	if len(stmt.Action) != 1 || stmt.Action[0] != "s3:GetObject" {
		t.Errorf("expected action s3:GetObject, got %v", stmt.Action)
	}
	if stmt.Resource != "arn:aws:s3:::bucket/key" {
		t.Errorf("expected resource arn:aws:s3:::bucket/key, got %s", stmt.Resource)
	}
}

func TestGenerate_MultipleCallsSameResource(t *testing.T) {
	calls := []ObservedCall{
		{Service: "s3", Action: "GetObject", Resource: "arn:aws:s3:::bucket/*"},
		{Service: "s3", Action: "PutObject", Resource: "arn:aws:s3:::bucket/*"},
		{Service: "s3", Action: "GetObject", Resource: "arn:aws:s3:::bucket/*"}, // duplicate
	}

	doc := Generate(calls)

	if len(doc.Statement) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(doc.Statement))
	}

	stmt := doc.Statement[0]
	if len(stmt.Action) != 2 {
		t.Errorf("expected 2 actions (deduplicated), got %d", len(stmt.Action))
	}
}

func TestGenerate_DifferentResources(t *testing.T) {
	calls := []ObservedCall{
		{Service: "s3", Action: "GetObject", Resource: "arn:aws:s3:::bucket1/*"},
		{Service: "s3", Action: "GetObject", Resource: "arn:aws:s3:::bucket2/*"},
	}

	doc := Generate(calls)

	if len(doc.Statement) != 2 {
		t.Errorf("expected 2 statements for different resources, got %d", len(doc.Statement))
	}
}

func TestGenerate_Deterministic(t *testing.T) {
	calls := []ObservedCall{
		{Service: "dynamodb", Action: "PutItem", Resource: "arn:aws:dynamodb:*:*:table/Users"},
		{Service: "s3", Action: "GetObject", Resource: "arn:aws:s3:::bucket/*"},
		{Service: "lambda", Action: "InvokeFunction", Resource: "arn:aws:lambda:*:*:function:MyFunc"},
	}

	// Generate twice and compare
	doc1 := Generate(calls)
	doc2 := Generate(calls)

	json1, _ := doc1.ToJSON()
	json2, _ := doc2.ToJSON()

	if string(json1) != string(json2) {
		t.Error("output should be deterministic")
	}
}

func TestToJSON_ValidFormat(t *testing.T) {
	calls := []ObservedCall{
		{Service: "s3", Action: "GetObject", Resource: "arn:aws:s3:::bucket/key"},
	}

	doc := Generate(calls)
	jsonBytes, err := doc.ToJSON()

	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Verify structure
	if parsed["Version"] != "2012-10-17" {
		t.Error("missing or incorrect Version field")
	}
	if _, ok := parsed["Statement"].([]interface{}); !ok {
		t.Error("missing or incorrect Statement field")
	}
}
