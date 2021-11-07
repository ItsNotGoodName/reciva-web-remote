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
	cfg.DB = path.Join(os.TempDir(), "test.db")
	os.Remove(cfg.DB)

	s, err := NewStore(cfg)
	if err != nil {
		t.Error(err)
	}
	if len(s.URLS) != 0 {
		t.Error("URLS should be empty")
	}

	cfg.URLS = []string{"http://example.com/01.m3u", "http://example.com/02.m3u"}
	cfg.DB = path.Join(os.TempDir(), "test2.db")
	os.Remove(cfg.DB)
	s, err = NewStore(cfg)
	if err != nil {
		t.Error(err)
	}
	if len(s.URLS) != 2 {
		t.Error("Presets should be empty")
	}
	if s.URLS[0] != "http://example.com/01.m3u" {
		t.Errorf("URL should be true, got %s", s.URLS[0])
	}
	if s.URLS[1] != "http://example.com/02.m3u" {
		t.Errorf("URL should be true, got %s", s.URLS[1])
	}
}

func TestPreset(t *testing.T) {
	testPreset := Preset{
		URL: "http://example.com/01.m3u",
	}

	cfg := config.NewConfig()
	cfg.DB = path.Join(os.TempDir(), "test.db")
	os.Remove(cfg.DB)

	s, err := NewStore(cfg)
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get presets test
	presets, err := s.ReadPresets(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(presets) != 0 {
		t.Errorf("Got %d presets, expected 0", len(presets))
	}

	// Create preset test
	err = s.CreatePreset(ctx, &testPreset)
	if err != nil {
		t.Error(err)
	}
	if testPreset.URL != "http://example.com/01.m3u" {
		t.Errorf("URL is was changed, got %s expected %s", testPreset.URL, "http://example.com/01.m3u")
	}
	if testPreset.SID != 0 {
		t.Errorf("SID is not set, got %d", testPreset.SID)
	}

	// Get presets test
	presets, err = s.ReadPresets(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(presets) != 1 {
		t.Errorf("Got %d presets, expected 1", len(presets))
	}
	if presets[0].URL != testPreset.URL {
		t.Errorf("Got preset with URL %s, expected %s", presets[0].URL, testPreset.URL)
	}
	if presets[0].SID != 0 {
		t.Errorf("Got preset with SID %d, expected 0", presets[0].SID)
	}

	// Get preset by URL and check if it's the same
	preset, err := s.ReadPreset(ctx, testPreset.URL)
	if err != nil {
		t.Error(err)
	}
	if preset.URL != testPreset.URL {
		t.Errorf("Got preset with URL %s, expected %s", preset.URL, testPreset.URL)
	}
	if preset.SID != 0 {
		t.Errorf("Got preset with SID %d, expected 0", preset.SID)
	}

	// Delete all preset test
	preset2 := testPreset
	preset2.URL = "http://example.com/02.m3u"
	err = s.CreatePreset(ctx, &preset2)
	if err != nil {
		t.Error(err)
	}
	presets, err = s.ReadPresets(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(presets) != 2 {
		t.Errorf("Got %d presets, expected 2", len(presets))
	}
}

func TestStream(t *testing.T) {
	TestPreset := Preset{
		URL: "http://example.com/01.m3u",
	}
	TestStream := Stream{
		Name:    "test",
		Content: "test",
	}

	cfg := config.NewConfig()
	cfg.DB = path.Join(os.TempDir(), "test.db")
	os.Remove(cfg.DB)

	s, err := NewStore(cfg)
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create preset test
	testPreset := TestPreset
	err = s.CreatePreset(ctx, &testPreset)
	if err != nil {
		t.Error(err)
	}

	// Get streams test
	streams, err := s.ReadStreams(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(streams) != 0 {
		t.Errorf("Got %d streams, expected 0", len(streams))
	}

	// Create stream test
	testStream := TestStream
	err = s.CreateStream(ctx, &testStream)
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
	ok, err := s.UpdatePreset(ctx, &testPreset)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Update failed")
	}
	preset, err := s.ReadPreset(ctx, testPreset.URL)
	if err != nil {
		t.Error(err)
	}
	if preset.SID != testPreset.SID {
		t.Errorf("Got preset with SID %d, expected %d", preset.SID, testPreset.SID)
	}

	// Get stream by preset test
	stream, err := s.ReadStream(ctx, testPreset.SID)
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
	stream, err = s.ReadStream(ctx, preset.SID)
	if err == nil {
		t.Errorf("Got stream with ID %d, expected nil", stream.ID)
	}

	// Get stream by id test
	stream, err = s.ReadStream(ctx, testStream.ID)
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
	streams, err = s.ReadStreams(ctx)
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
	ok, err = s.UpdateStream(ctx, &testStream)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Update stream failed")
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
	ok, err = s.DeleteStream(ctx, &testStream)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Stream was not deleted")
	}
	stream, err = s.ReadStream(ctx, testStream.ID)
	if err == nil {
		t.Error("Got stream with ID ", stream.ID, ", expected error")
	}

	// Make sure preset's SID is 0
	preset, err = s.ReadPreset(ctx, testPreset.URL)
	if err != nil {
		t.Error(err)
	}
	if preset.SID != 0 {
		t.Errorf("Got preset with SID %d, expected 0", preset.SID)
	}
}
