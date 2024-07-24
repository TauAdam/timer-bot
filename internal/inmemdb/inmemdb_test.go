package inmemdb_test

import (
	"github.com/TauAdam/timer-bot/internal/inmemdb"
	"github.com/TauAdam/timer-bot/internal/timer"
	"reflect"
	"testing"
	"time"
)

func TestInMemoryDB_AddTimer(t *testing.T) {
	type args struct {
		id string
		t  timer.Timer
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		errMessage string
	}{
		{
			name: "Add valid timer",
			args: args{
				id: "1",
				t:  timer.Timer{Duration: 10 * time.Second, StartTime: time.Now()},
			},
			wantErr: false,
		},
		{
			name: "Add timer with empty ID",
			args: args{
				id: "",
				t:  timer.Timer{Duration: 5 * time.Second, StartTime: time.Now()},
			},
			wantErr:    true,
			errMessage: "id cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := inmemdb.NewInMemoryDB()
			if err := db.AddTimer(tt.args.id, tt.args.t); (err != nil) != tt.wantErr && err.Error() != tt.errMessage {
				t.Errorf("AddTimer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemoryDB_DeleteTimer(t *testing.T) {
	db := inmemdb.NewInMemoryDB()
	_ = db.AddTimer("1", timer.Timer{Duration: 10 * time.Second, StartTime: time.Now()})
	tests := []struct {
		name       string
		id         string
		wantErr    bool
		errMessage string
	}{
		{
			name:    "Delete existing timer",
			id:      "1",
			wantErr: false,
		},
		{
			name:       "Delete non-existing timer",
			id:         "2",
			wantErr:    true,
			errMessage: "timer not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.ResetTimer(tt.id); (err != nil) != tt.wantErr && err.Error() != tt.errMessage {
				t.Errorf("ResetTimer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemoryDB_GetTimer(t *testing.T) {
	db := inmemdb.NewInMemoryDB()
	expectedTimer := timer.Timer{Duration: 10 * time.Second, StartTime: time.Now()}
	_ = db.AddTimer("1", expectedTimer)
	tests := []struct {
		name       string
		id         string
		want       timer.Timer
		wantErr    bool
		errMessage string
	}{
		{
			name:    "Get existing timer",
			id:      "1",
			want:    expectedTimer,
			wantErr: false,
		},
		{
			name:       "Get non-existing timer",
			id:         "2",
			want:       timer.Timer{},
			wantErr:    true,
			errMessage: "timer not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.GetTimer(tt.id)
			if (err != nil) != tt.wantErr && err.Error() != tt.errMessage {
				t.Errorf("GetTimer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTimer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
