// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package data

import "testing"

func TestStringBoolean_ProperlyUnmarshals(t *testing.T) {
	sb := StringBoolean(false)
	err := sb.UnmarshalJSON([]byte("false"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if sb != false {
		t.Errorf("Expected false, got %v", sb)
	}

	sb = StringBoolean(true)
	err = sb.UnmarshalJSON([]byte("true"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if sb != true {
		t.Errorf("Expected true, got %v", sb)
	}
}

func TestStringBoolean_UnmarshalsAnyStringToTrue(t *testing.T) {
	sb := StringBoolean(true)
	err := sb.UnmarshalJSON([]byte("yehyehyeh foobar"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if sb != true {
		t.Errorf("Expected true, got %v", sb)
	}
}

func TestWatchItem_GetMetrics(t *testing.T) {
	w := WatchItem{
		Title: "Test",
		Url:   "https://example.com",
	}
	metrics, err := w.GetMetrics()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(metrics) != 2 {
		t.Errorf("Expected 2 metrics, got %v", len(metrics))
	}
	if metrics[0] != "Test" {
		t.Errorf("Expected Test, got %v", metrics[0])
	}
	if metrics[1] != "example.com" {
		t.Errorf("Expected example.com, got %v", metrics[1])
	}
}

func TestWatchItem_GetMetrics_EmptyTitle(t *testing.T) {
	w := WatchItem{
		Url: "https://example.com",
	}
	_, err := w.GetMetrics()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestWatchItem_GetMetrics_InvalidUri(t *testing.T) {
	w := WatchItem{
		Title: "Test",
		Url:   "foo-bar-is-not-a-uri",
	}
	_, err := w.GetMetrics()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestWatchItem_GetMetrics_EmptyHost(t *testing.T) {
	w := WatchItem{
		Title: "Test",
		Url:   "http://",
	}
	_, err := w.GetMetrics()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
