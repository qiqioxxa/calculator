package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCalc(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantRes    float64
		wantErr    bool
	}{
		{"add", http.MethodPost, `{"op":"add","a":4,"b":6}`, http.StatusOK, 10, false},
		{"div by zero", http.MethodPost, `{"op":"div","a":10,"b":0}`, http.StatusOK, 0, true},
		{"sqrt", http.MethodPost, `{"op":"sqrt","a":16,"b":0}`, http.StatusOK, 4, false},
		{"unknown op", http.MethodPost, `{"op":"bad","a":1,"b":2}`, http.StatusOK, 0, true},
		{"wrong method", http.MethodGet, ``, http.StatusMethodNotAllowed, 0, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/api/calc", bytes.NewBufferString(tc.body))
			rec := httptest.NewRecorder()

			handleCalc(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d; want %d", rec.Code, tc.wantStatus)
			}
			if tc.wantStatus != http.StatusOK {
				return
			}

			var resp calcResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if tc.wantErr && resp.Error == "" {
				t.Errorf("expected error in response, got empty")
			}
			if !tc.wantErr {
				if resp.Error != "" {
					t.Errorf("unexpected error: %s", resp.Error)
				}
				if resp.Result != tc.wantRes {
					t.Errorf("result = %v; want %v", resp.Result, tc.wantRes)
				}
			}
		})
	}
}

func TestHandleMemory(t *testing.T) {
	calcMem.Clear()

	steps := []struct {
		op   string
		val  float64
		want float64
	}{
		{"m+", 25, 25},
		{"m-", 5, 20},
		{"mr", 0, 20},
		{"mc", 0, 0},
	}

	for _, step := range steps {
		payload, _ := json.Marshal(memoryRequest{Op: step.op, Value: step.val})
		req := httptest.NewRequest(http.MethodPost, "/api/memory", bytes.NewReader(payload))
		rec := httptest.NewRecorder()

		handleMemory(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("op %s: status = %d; want 200", step.op, rec.Code)
		}

		var resp memoryResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("op %s: decode error: %v", step.op, err)
		}
		if resp.Memory != step.want {
			t.Errorf("op %s: memory = %v; want %v", step.op, resp.Memory, step.want)
		}
	}
}
