package tasks

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/poy/go-dependency-injection/pkg/injection"
)

// StorePath is the path to save the store values to. If unset, the the store
// will just be in-memory.
type StorePath string

// ProvideStorePath saves the store to the given path.
func ProvideStorePath(p string) {
	injection.Register[StorePath](
		func(ctx context.Context) StorePath {
			return StorePath(p)
		},
	)
}

func init() {
	injection.Register[Store](
		func(ctx context.Context) Store {
			storePath, ok := injection.TryResolve[StorePath](ctx)

			s := &store{}
			save := func() {}
			if ok {
				save = func() {
					f, err := os.Create(string(storePath))
					if err != nil {
						log.Fatalf("failed to create file %s: %v", storePath, err)
					}
					defer f.Close()
					if err := json.NewEncoder(f).Encode(s.tasks); err != nil {
						log.Fatalf("failed to encode tasks: %v", err)
					}
				}

				// Try loading the file. If it doesn't work, just move on.
				f, err := os.Open(string(storePath))
				if err != nil {
					if !os.IsNotExist(err) {
						log.Printf("failed to open store file: %s", err)
					}
				} else if err := json.NewDecoder(f).Decode(&s.tasks); err != nil {
					log.Printf("failed to decode store file: %s", err)
				} else {
					for i := range s.tasks {
						s.tasks[i].save = save
					}
				}

			}
			s.save = save
			return s
		},
	)
}

// Store is a store of Tasks.
type Store interface {
	// Add adds a new Task to the store.
	Add(name, description string)
	// Remove removes a Task from the store.
	Remove(name string)
	// TaskNames returns the tasks in the store.
	TaskNames() []string
	// GetTask returns the Task with the given name.
	GetTask(name string) *Task
}

// store keeps track of the tasks.
type store struct {
	tasks []*Task
	save  func()
}

// Task represents a task.
type Task struct {
	name        string
	datetime    int64
	description string
	completed   *int64
	notes       []Note
	save        func()
}

func (t *Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"name":        t.name,
		"datetime":    t.datetime,
		"description": t.description,
		"completed":   t.completed,
		"notes":       t.notes,
	})
}

func (t *Task) UnmarshalJSON(b []byte) error {
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	t.name, _ = m["name"].(string)
	t.datetime = int64(m["datetime"].(float64))
	t.description, _ = m["description"].(string)
	t.completed, _ = m["completed"].(*int64)
	t.notes, _ = m["notes"].([]Note)
	return nil
}

// Name returns the name of the Task.
func (t *Task) Name() string {
	return t.name
}

// Description returns the description of the Task.
func (t *Task) Description() string {
	return t.description
}

// Completed returns if the task has been comleted.
func (t *Task) Completed() bool {
	return t.completed != nil
}

// CompletedAt returns the time the task was completed.
func (t *Task) CompletedAt() time.Time {
	if t.completed == nil {
		return time.Time{}
	}
	return time.Unix(0, *t.completed)
}

// Complete marks the task as completed.
func (t *Task) Complete() {
	now := time.Now().Unix()
	t.completed = &now
	t.save()
}

// AddNotes adds notes to the task.
func (t *Task) AddNotes(notes ...string) {
	for _, note := range notes {
		t.notes = append(t.notes, Note{
			datetime: time.Now().UnixNano(),
			note:     note,
		})
	}
	t.save()
}

// Notes returns the notes for the task.
func (t *Task) Notes() []Note {
	return t.notes
}

// Note is a note about a task.
type Note struct {
	datetime int64  `json:"datetime"`
	note     string `json:"note"`
}

func (n *Note) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"datetime": n.datetime,
		"note":     n.note,
	})
}

func (n *Note) UnmarshalJSON(b []byte) error {
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	n.datetime = int64(m["datetime"].(float64))
	n.note, _ = m["note"].(string)
	return nil
}

// Datetime returns the time the note was created.
func (n *Note) Datetime() time.Time {
	return time.Unix(0, n.datetime)
}

// Note adds a note to the task.
func (n *Note) Note() string {
	return n.note
}

// Add adds a new Task to the store.
func (s *store) Add(name, description string) {
	if s.GetTask(name) != nil {
		return
	}

	s.tasks = append(s.tasks, &Task{
		name:        name,
		datetime:    time.Now().UnixNano(),
		description: description,
		save:        s.save,
	})
	s.save()
}

// Remove removes a Task from the store.
func (s *store) Remove(name string) {
	name = strings.ToLower(name)
	for i, t := range s.tasks {
		if strings.ToLower(t.name) == name {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			s.save()
			return
		}
	}
}

// TaskNames returns the tasks in the store.
func (s *store) TaskNames() []string {
	var names []string
	for _, t := range s.tasks {
		names = append(names, t.name)
	}
	return names
}

// GetTask returns the Task with the given name.
func (s *store) GetTask(name string) *Task {
	name = strings.ToLower(name)
	for _, t := range s.tasks {
		if strings.ToLower(t.name) == name {
			return t
		}
	}
	return nil
}
