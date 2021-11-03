package store

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
)

func TestStore(t *testing.T) {
	cfg := config.NewConfig()
	cfg.DBPath = path.Join(os.TempDir(), "test.db")
	os.Remove(cfg.DBPath)

	s, err := NewStore(cfg)
	if err != nil {
		t.Error(err)
	}
	if len(s.Presets) != 0 {
		t.Error("Presets should be empty")
	}

	cfg.URIS = []string{"/01.m3u", "/02.m3u"}
	cfg.DBPath = path.Join(os.TempDir(), "test2.db")
	os.Remove(cfg.DBPath)
	s, err = NewStore(cfg)
	if err != nil {
		t.Error(err)
	}
	if len(s.Presets) != 2 {
		t.Error("Presets should be empty")
	}
	if s.Presets["/01.m3u"].URI != "/01.m3u" {
		t.Errorf("URI should be /01.m3u, got %s", s.Presets["/01.m3u"].URI)
	}
	if s.Presets["/02.m3u"].URI != "/02.m3u" {
		t.Errorf("URI should be /02.m3u, got %s", s.Presets["/02.m3u"].URI)
	}
}

func TestPreset(t *testing.T) {
	testPreset := Preset{
		URI: "/01.m3u",
	}

	cfg := config.NewConfig()
	cfg.DBPath = path.Join(os.TempDir(), "test.db")
	os.Remove(cfg.DBPath)

	s, err := NewStore(cfg)
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get presets test
	presets, err := s.GetPresets(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(presets) != 0 {
		t.Errorf("Got %d presets, expected 0", len(presets))
	}

	// Add preset test
	err = s.AddPreset(ctx, &testPreset)
	if err != nil {
		t.Error(err)
	}
	if testPreset.ID == 0 {
		t.Error("ID is not set")
	}
	if testPreset.URI != "/01.m3u" {
		t.Errorf("URI is was changed, got %s expected %s", testPreset.URI, "/01.m3u")
	}
	if testPreset.SID != 0 {
		t.Errorf("SID is not set, got %d", testPreset.SID)
	}

	// Get presets test
	presets, err = s.GetPresets(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(presets) != 1 {
		t.Errorf("Got %d presets, expected 1", len(presets))
	}
	if presets[0].ID != testPreset.ID {
		t.Errorf("Got preset with ID %d, expected %d", presets[0].ID, testPreset.ID)
	}
	if presets[0].URI != testPreset.URI {
		t.Errorf("Got preset with URI %s, expected %s", presets[0].URI, testPreset.URI)
	}
	if presets[0].SID != 0 {
		t.Errorf("Got preset with SID %d, expected 0", presets[0].SID)
	}

	// Get preset by URI and check if it's the same
	preset, err := s.GetPresetByURI(ctx, testPreset.URI)
	if err != nil {
		t.Error(err)
	}
	if preset.ID != testPreset.ID {
		t.Errorf("Got preset with ID %d, expected %d", preset.ID, testPreset.ID)
	}
	if preset.URI != testPreset.URI {
		t.Errorf("Got preset with URI %s, expected %s", preset.URI, testPreset.URI)
	}
	if preset.SID != 0 {
		t.Errorf("Got preset with SID %d, expected 0", preset.SID)
	}

	// Delete all preset test
	preset2 := testPreset
	preset2.URI = "/02.m3u"
	err = s.AddPreset(ctx, &preset2)
	if err != nil {
		t.Error(err)
	}
	presets, err = s.GetPresets(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(presets) != 2 {
		t.Errorf("Got %d presets, expected 2", len(presets))
	}
	err = s.DeleteAllPresets(ctx)
	if err != nil {
		t.Error(err)
	}

	// Get presets test
	presets, err = s.GetPresets(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(presets) != 0 {
		t.Error("Expected 0 presets, got ", len(presets))
	}

	// Get preset test
	preset, err = s.GetPresetByURI(ctx, testPreset.URI)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if preset != nil {
		t.Errorf("Got preset with ID %d, expected nil", preset.ID)
	}
}

func TestStream(t *testing.T) {
	TestPreset := Preset{
		URI: "/01.m3u",
	}
	TestStream := Stream{
		Name:    "test",
		Content: "test",
	}

	cfg := config.NewConfig()
	cfg.DBPath = path.Join(os.TempDir(), "test.db")
	os.Remove(cfg.DBPath)

	s, err := NewStore(cfg)
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Add preset test
	testPreset := TestPreset
	err = s.AddPreset(ctx, &testPreset)
	if err != nil {
		t.Error(err)
	}

	// Get streams test
	streams, err := s.GetStreams(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(streams) != 0 {
		t.Errorf("Got %d streams, expected 0", len(streams))
	}

	// Add stream test
	testStream := TestStream
	err = s.AddStream(ctx, &testStream)
	if err != nil {
		t.Error(err)
	}
	if testStream.ID == 0 {
		t.Error("ID is not set")
	}
	if testStream.Name != TestStream.Name {
		t.Errorf("Name is was changed, got %s expected %s", testStream.Name, TestStream.Name)
	}
	if testStream.Content != TestStream.Content {
		t.Errorf("Content is was changed, got %s expected %s", testStream.Content, TestStream.Content)
	}

	// Update preset SID test
	testPreset.SID = testStream.ID
	err = s.UpdatePresetSID(ctx, &testPreset)
	if err != nil {
		t.Error(err)
	}
	preset, err := s.GetPresetByURI(ctx, testPreset.URI)
	if err != nil {
		t.Error(err)
	}
	if preset.SID != testPreset.SID {
		t.Errorf("Got preset with SID %d, expected %d", preset.SID, testPreset.SID)
	}

	// Get stream by preset test
	stream, err := s.GetStreamByPreset(ctx, &testPreset)
	if err != nil {
		t.Error(err)
	}
	if stream.ID != testStream.ID {
		t.Errorf("Got stream with ID %d, expected %d", stream.ID, testStream.ID)
	}

	// Clear preset test
	err = s.ClearPreset(ctx, preset)
	if err != nil {
		t.Error(err)
	}
	stream, err = s.GetStreamByPreset(ctx, preset)
	if err == nil {
		t.Errorf("Got stream with ID %d, expected nil", stream.ID)
	}

	// Get stream by id test
	stream, err = s.GetStreamByID(ctx, testStream.ID)
	if err != nil {
		t.Error(err)
	}
	if stream.ID != testStream.ID {
		t.Errorf("Got stream with ID %d, expected %d", stream.ID, testStream.ID)
	}
	if stream.Name != testStream.Name {
		t.Errorf("Got stream with Name %s, expected %s", stream.Name, testStream.Name)
	}
	if stream.Content != testStream.Content {
		t.Errorf("Got stream with Content %s, expected %s", stream.Content, testStream.Content)
	}

	// Get streams test
	streams, err = s.GetStreams(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(streams) != 1 {
		t.Errorf("Got %d streams, expected 1", len(streams))
	}
	if streams[0].ID != testStream.ID {
		t.Errorf("Got stream with ID %d, expected %d", streams[0].ID, testStream.ID)
	}
	if streams[0].Name != testStream.Name {
		t.Errorf("Got stream with name %s, expected %s", streams[0].Name, testStream.Name)
	}
	if streams[0].Content != testStream.Content {
		t.Errorf("Got stream with content %s, expected %s", streams[0].Content, testStream.Content)
	}

	// Update stream test
	testStream.Name = "test2"
	testStream.Content = "test2"
	err = s.UpdateStream(ctx, &testStream)
	if err != nil {
		t.Error(err)
	}
	if testStream.ID != stream.ID {
		t.Errorf("Got stream with ID %d, expected %d", testStream.ID, stream.ID)
	}
	if testStream.Name != "test2" {
		t.Errorf("Got stream with name %s, expected %s", testStream.Name, "test2")
	}
	if testStream.Content != "test2" {
		t.Errorf("Got stream with content %s, expected %s", testStream.Content, "test2")
	}

	// Delete stream test
	err = s.DeleteStream(ctx, &testStream)
	if err != nil {
		t.Error(err)
	}
	stream, err = s.GetStreamByID(ctx, testStream.ID)
	if err == nil {
		t.Error("Got stream with ID ", stream.ID, ", expected error")
	}

	// Make sure preset's SID is 0
	preset, err = s.GetPresetByURI(ctx, testPreset.URI)
	if err != nil {
		t.Error(err)
	}
	if preset.SID != 0 {
		t.Errorf("Got preset with SID %d, expected 0", preset.SID)
	}
}
