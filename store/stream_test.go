package store

import (
	"testing"
)

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

func compareAddTest(t *testing.T, name string, addTest *AddTest, st *Stream) {
	if st.SID != addTest.sid {
		t.Fatalf("%s(st.SID) = %d, want %d", name, st.SID, addTest.sid)
	}
	if st.Name != addTest.name {
		t.Fatalf("%s(st.Name) = %q, want %q", name, st.Name, addTest.name)
	}
	if st.Content != addTest.content {
		t.Fatalf("%s(st.Content) = %q, want %q", name, st.Content, addTest.content)
	}
}

func TestStream(t *testing.T) {
	s := testStore(t)

	for _, addTest := range addTests {
		// AddTest
		st, ok := s.AddStream(addTest.name, addTest.content)
		if ok != addTest.ret {
			t.Fatalf("TestStream(%q, %q) = %t, want %t", addTest.name, addTest.content, ok, addTest.ret)
		}
		if ok {
			compareAddTest(t, "AddTest", &addTest, st)
			// GetTest
			st, ok = s.GetStream(st.SID)
			if !ok {
				t.Fatal("GetTest(ok) = false, want true")
			}
			if ok {
				compareAddTest(t, "GetTest", &addTest, st)
			}
		}
	}
}
