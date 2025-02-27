package transformer

import (
	"github.com/patrykorwat/infraconv/internal/parser"
	"testing"
)

func Test_crossplaneTransformer_Transform(t *testing.T) {
	type args struct {
		cfg             *parser.Config
		directoryOutput string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := crossplaneTransformer{}
			// GOLANG 1.24 - Feat 3: Change the working directory for the current test.
			t.Chdir("/tmp")
			// GOLANG 1.24 - Feat 4: test context
			if err := c.Transform(t.Context(), tt.args.cfg, tt.args.directoryOutput); (err != nil) != tt.wantErr {
				t.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
