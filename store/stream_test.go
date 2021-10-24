package store

import (
	"testing"
)

func TestAddStream(t *testing.T) {
	type AddTest struct {
		sid     int
		name    string
		content string
		ret     bool
	}

	var addTests = []AddTest{
		{1, "Name", "Content", true},
		{0, "Name", "Content", false},
		{0, "Name", "tent", false},
		{2, "name", "tent", true},
		{3, "", "tent", true},
		{0, "", "tent", false},
		{4, "hi", "", true},
		{0, "hi", "hi", false},
	}

	s := testStore(t)

	for _, addTest := range addTests {
		st, ok := s.AddStream(addTest.name, addTest.content)
		if ok != addTest.ret {
			t.Fatalf("TestStream(%q, %q) = %t, want %t", addTest.name, addTest.content, ok, addTest.ret)
		}
		if ok {
			if st.SID != addTest.sid {
				t.Fatalf("TestStream(st.SID) = %d, want %d", st.SID, addTest.sid)
			}
			if st.Name != addTest.name {
				t.Fatalf("TestStream(st.Name) = %q, want %q", st.Name, addTest.name)
			}
			if st.Content != addTest.content {
				t.Fatalf("TestStream(st.Content) = %q, want %q", st.Content, addTest.content)
			}
		}

	}

	//	// Test adding a stream and getting the stream back
	//	addST, err := s.AddStream("Name", "Content")
	//	if err != nil {
	//		t.Error("could not add stream,", err)
	//	}
	//	if _, err = s.AddStream("Name", "Content"); err == nil {
	//		t.Error("duplicate add stream name")
	//	}
	//	testGetStream := func() {
	//		if getStream, ok := s.GetStream(addST.SID); !ok {
	//			t.Error("stream should exist")
	//		} else {
	//			if !reflect.DeepEqual(*addST, *getStream) {
	//				t.Error("saved stream is not equal", addST, *getStream)
	//			}
	//		}
	//	}
	//	testGetStream()
	//
	//	// Test getting the stream back after writing and reading
	//	s.WriteSettings()
	//	s.ReadSettings()
	//	testGetStream()
	//
	//	// Test deleting stream
	//	if c := s.DeleteStream(addST.SID); c != 1 {
	//		t.Error("deleting stream should return 1, got", c)
	//	}
	//	if _, ok := s.GetStream(addST.SID); ok {
	//		t.Error("stream should be deleted")
	//	}
}
