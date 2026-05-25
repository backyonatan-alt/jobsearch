package ics

import (
	"strings"
	"testing"
	"time"
)

const googleSample = "BEGIN:VCALENDAR\r\n" +
	"VERSION:2.0\r\n" +
	"PRODID:-//Google Inc//Google Calendar 70.9054//EN\r\n" +
	"BEGIN:VEVENT\r\n" +
	"UID:abc123@google.com\r\n" +
	"DTSTART:20260601T140000Z\r\n" +
	"DTEND:20260601T150000Z\r\n" +
	"SUMMARY:Anthropic onsite — system design\r\n" +
	"LOCATION:Anthropic HQ\\, 580 Howard St\\, San Francisco\r\n" +
	"DESCRIPTION:Round 2 of 4.\\nBring laptop.\r\n" +
	"ORGANIZER;CN=Recruiting:mailto:recruiting@anthropic.com\r\n" +
	"ATTENDEE;CN=Jane Doe;ROLE=REQ-PARTICIPANT:mailto:jane@anthropic.com\r\n" +
	"ATTENDEE;CN=John Smith:mailto:john@anthropic.com\r\n" +
	"END:VEVENT\r\n" +
	"END:VCALENDAR\r\n"

func TestParseGoogleStyle(t *testing.T) {
	events, err := Parse(strings.NewReader(googleSample))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("want 1 event, got %d", len(events))
	}
	e := events[0]

	want := time.Date(2026, 6, 1, 14, 0, 0, 0, time.UTC)
	if !e.StartsAt.Equal(want) {
		t.Errorf("start: want %s, got %s", want, e.StartsAt)
	}
	if e.Duration() != time.Hour {
		t.Errorf("duration: want 1h, got %s", e.Duration())
	}
	if e.UID != "abc123@google.com" {
		t.Errorf("UID: got %q", e.UID)
	}
	if e.Summary != "Anthropic onsite — system design" {
		t.Errorf("summary: got %q", e.Summary)
	}
	if e.Location != "Anthropic HQ, 580 Howard St, San Francisco" {
		t.Errorf("location unescape: got %q", e.Location)
	}
	if e.Description != "Round 2 of 4.\nBring laptop." {
		t.Errorf("description unescape: got %q", e.Description)
	}
	if e.Organizer.Email != "recruiting@anthropic.com" {
		t.Errorf("organizer email: got %q", e.Organizer.Email)
	}
	if len(e.Attendees) != 2 || e.Attendees[0].Name != "Jane Doe" || e.Attendees[0].Email != "jane@anthropic.com" {
		t.Errorf("attendees: got %+v", e.Attendees)
	}
}

func TestParseTZID(t *testing.T) {
	src := "BEGIN:VCALENDAR\nBEGIN:VEVENT\n" +
		"UID:tz@x\n" +
		"DTSTART;TZID=America/New_York:20260601T100000\n" +
		"DTEND;TZID=America/New_York:20260601T110000\n" +
		"SUMMARY:NYC interview\n" +
		"END:VEVENT\nEND:VCALENDAR\n"
	events, err := Parse(strings.NewReader(src))
	if err != nil || len(events) != 1 {
		t.Fatalf("parse: %v / %d", err, len(events))
	}
	loc, _ := time.LoadLocation("America/New_York")
	want := time.Date(2026, 6, 1, 10, 0, 0, 0, loc)
	if !events[0].StartsAt.Equal(want) {
		t.Errorf("tz start: want %s, got %s", want, events[0].StartsAt)
	}
}

func TestParseAllDay(t *testing.T) {
	src := "BEGIN:VCALENDAR\nBEGIN:VEVENT\n" +
		"UID:d@x\n" +
		"DTSTART;VALUE=DATE:20260601\n" +
		"SUMMARY:Take-home due\n" +
		"END:VEVENT\nEND:VCALENDAR\n"
	events, err := Parse(strings.NewReader(src))
	if err != nil || len(events) != 1 {
		t.Fatalf("parse: %v / %d", err, len(events))
	}
	if !events[0].AllDay {
		t.Error("expected AllDay=true")
	}
}

func TestLineFolding(t *testing.T) {
	src := "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\n" +
		"UID:fold@x\r\n" +
		"DTSTART:20260601T140000Z\r\n" +
		"SUMMARY:Very long summary that wraps acr\r\n" +
		" oss two lines\r\n" +
		"END:VEVENT\r\nEND:VCALENDAR\r\n"
	events, err := Parse(strings.NewReader(src))
	if err != nil || len(events) != 1 {
		t.Fatalf("parse: %v / %d", err, len(events))
	}
	if events[0].Summary != "Very long summary that wraps across two lines" {
		t.Errorf("fold: got %q", events[0].Summary)
	}
}

func TestVAlarmInsideVEventDoesNotTerminate(t *testing.T) {
	src := "BEGIN:VCALENDAR\nBEGIN:VEVENT\n" +
		"UID:alarm@x\n" +
		"DTSTART:20260601T140000Z\n" +
		"SUMMARY:With alarm\n" +
		"BEGIN:VALARM\nTRIGGER:-PT30M\nACTION:DISPLAY\nEND:VALARM\n" +
		"LOCATION:After alarm\n" +
		"END:VEVENT\nEND:VCALENDAR\n"
	events, err := Parse(strings.NewReader(src))
	if err != nil || len(events) != 1 {
		t.Fatalf("parse: %v / %d", err, len(events))
	}
	if events[0].Location != "After alarm" {
		t.Errorf("alarm nesting: location got %q", events[0].Location)
	}
}

func TestMultipleEvents(t *testing.T) {
	src := "BEGIN:VCALENDAR\n" +
		"BEGIN:VEVENT\nUID:a\nDTSTART:20260601T140000Z\nSUMMARY:A\nEND:VEVENT\n" +
		"BEGIN:VEVENT\nUID:b\nDTSTART:20260602T140000Z\nSUMMARY:B\nEND:VEVENT\n" +
		"END:VCALENDAR\n"
	events, err := Parse(strings.NewReader(src))
	if err != nil || len(events) != 2 {
		t.Fatalf("parse: %v / %d", err, len(events))
	}
	if events[0].UID != "a" || events[1].UID != "b" {
		t.Errorf("order/UIDs: %+v", events)
	}
}
