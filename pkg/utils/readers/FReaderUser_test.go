package readers

import (
	"cli-project/internal/domain/models"
	"reflect"
	"testing"
)

func TestFReaderUser(t *testing.T) {
	type args struct {
		f    string
		flag int
	}
	var tests []struct {
		name string
		args args
		want []models.User
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FReaderUser(tt.args.f, tt.args.flag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FReaderUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
