package store

import (
	"testing"
)

type AddStreamTest struct {
	sid     int
	name    string
	content string
	ret     bool
}

var addStreamTests = []AddStreamTest{
	{1, "Name", "Content", true},
	{0, "Name", "Content", false},
	{0, "Name", "tent", false},
	{2, "name", "tent", true},
	{3, "", "tent", true},
	{0, "", "tent", false},
	{4, "hi", "", true},
	{0, "hi", "hi", false},
}

func runAddStreamTest(t *testing.T, name string, at *AddStreamTest, st *Stream) {
	if st.SID != at.sid {
		t.Errorf("%s(st.SID) = %d, want %d", name, st.SID, at.sid)
	}
	if st.Name != at.name {
		t.Errorf("%s(st.Name) = %q, want %q", name, st.Name, at.name)
	}
	if st.Content != at.content {
		t.Errorf("%s(st.Content) = %q, want %q", name, st.Content, at.content)
	}
}

func TestStream(t *testing.T) {
	s := testStore(t)

	for _, at := range addStreamTests {
		// AddStream
		st, ok := s.AddStream(at.name, at.content)
		if ok != at.ret {
			t.Fatalf("AddStream(%q, %q) = %t, want %t", at.name, at.content, ok, at.ret)
		}
		if ok {
			runAddStreamTest(t, "AddStream", &at, st)
			// GetStream
			st, ok = s.GetStream(st.SID)
			if !ok {
				t.Fatal("GetStream(ok) = false, want true")
			}
			if ok {
				runAddStreamTest(t, "GetStream", &at, st)
			}
		}
	}

	// GetStreams
	sts := s.GetStreams()
	count := 0
	for _, addTest := range addStreamTests {
		if addTest.ret {
			runAddStreamTest(t, "GetStream", &addTest, &sts[count])
			count += 1
		}
	}

	// UpdateStream
	for _, at := range addStreamTests {
		if at.ret {
			mod := "UpdateStream"

			st, ok := s.GetStream(at.sid)
			if !ok {
				t.Fatal("UpdateStream: stream not found")
			}

			st.Name = st.Name + mod
			at.name = st.Name
			st.Content = st.Content + mod
			at.content = st.Content

			if !s.UpdateStream(st) {
				t.Fatal("UpdateStream: stream could not be updated")
			}

			st, ok = s.GetStream(at.sid)
			if !ok {
				t.Fatal("UpdateStream: stream not found")
			}
			runAddStreamTest(t, "UpdateStream", &at, st)

			// DeleteStream
			s.DeleteStream(st.SID)
			_, ok = s.GetStream(st.SID)
			if ok {
				t.Fatal("UpdateStream: stream shuold not exist")
			}

			break
		}
	}
}
