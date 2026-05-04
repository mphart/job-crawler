package profiles

import (
	"testing"

	"job-crawler/apps/api/internal/platform/db"
)

type fakeDispatcher struct {
	userID  string
	content string
	called  bool
}

func (f *fakeDispatcher) Enqueue(userID, resumeContentBase64 string) error {
	f.called = true
	f.userID = userID
	f.content = resumeContentBase64
	return nil
}

func TestUpdateMe_EnqueuesResumeParsingWhenResumeProvided(t *testing.T) {
	store := InMemoryStore{Inner: db.NewStore()}
	dispatcher := &fakeDispatcher{}
	service := Service{Store: store, Dispatcher: dispatcher}

	_, ok, err := service.UpdateMe("u_1", map[string]any{
		"resumeFileName":      "resume.pdf",
		"resumeContentBase64": "ZmFrZS1yZXN1bWU=",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected profile update success")
	}
	if !dispatcher.called {
		t.Fatalf("expected dispatcher to be called")
	}
}
