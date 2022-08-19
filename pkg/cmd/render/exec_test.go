// Copyright Red Hat
package render

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestOptions_Complete(t *testing.T) {
	type fields struct {
		Header     string
		Paths      []string
		ValuesPath string
		Values     map[string]interface{}
		OutputFile string
		SortOnKind bool
		OutputDir  string
		Excluded   []string
	}
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "read value file succees",
			fields: fields{
				ValuesPath: "../../../test/unit/resources/scenario/values.yaml",
			},
			wantErr: false,
		},
		{
			name: "read value file not found",
			fields: fields{
				ValuesPath: "file_not_found.yaml",
			},
			wantErr: true,
		},
		{
			name:    "read value stdin",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				Header:     tt.fields.Header,
				Paths:      tt.fields.Paths,
				ValuesPath: tt.fields.ValuesPath,
				Values:     tt.fields.Values,
				OutputFile: tt.fields.OutputFile,
				SortOnKind: tt.fields.SortOnKind,
				OutputDir:  tt.fields.OutputDir,
				Exclude:    tt.fields.Excluded,
			}
			var fileIn *os.File
			var err error
			if len(o.ValuesPath) == 0 {
				fileIn, err = ioutil.TempFile("", "stdin")
				if err != nil {
					t.Error(err)
				}
				defer os.Remove(fileIn.Name())
				err = ioutil.WriteFile(fileIn.Name(), []byte("ServiceAccount: my-sa\n"), 0600)
				if err != nil {
					t.Error(err)
				}
				os.Stdin = fileIn
			}
			if err := o.Complete(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Options.Complete() error = %v, wantErr %v", err, tt.wantErr)
			}
			switch tt.name {
			case "read value file succees":
				iSimple, ok := o.Values["Simple"]
				if !ok {
					t.Error("'Simple' not found in value")
				}
				simple := iSimple.(map[string]interface{})
				iSA, ok := simple["ServiceAccount"]
				if !ok {
					t.Error("'ServiceAccount' not found in value")
				}
				sa := iSA.(string)
				if sa != "my-sa" {
					t.Errorf("'Expected 'my-sa' got %s", sa)
				}
				if o.OutputFile != os.Stdout.Name() {
					t.Errorf("'Expected %s got %s", os.Stdout.Name(), o.OutputFile)
				}
			case "read value stdin":
				iSA, ok := o.Values["ServiceAccount"]
				if !ok {
					t.Error("'ServiceAccount' not found in value")
				}
				sa := iSA.(string)
				if sa != "my-sa" {
					t.Errorf("'Expected 'my-sa' got %s", sa)
				}
			}
		})
	}
}

func TestOptions_Validate(t *testing.T) {
	type fields struct {
		Header     string
		Paths      []string
		ValuesPath string
		Values     map[string]interface{}
		OutputFile string
		SortOnKind bool
		OutputDir  string
		Excluded   []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "directory succees",
			fields: fields{
				Header: "../../../test/unit/resources/scenario/musttemplateasset/header.txt",
				Paths:  []string{"../../../test/unit/resources/scenario/musttemplateasset"},
			},
			wantErr: false,
		},
		{
			name: "directory failed",
			fields: fields{
				Header: "../../../test/unit/resources/scenario/musttemplateasset/header.txt",
				Paths:  []string{"wrong_dir"},
			},
			wantErr: true,
		},
		{
			name:    "empty failed",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				Header:     tt.fields.Header,
				Paths:      tt.fields.Paths,
				ValuesPath: tt.fields.ValuesPath,
				Values:     tt.fields.Values,
				OutputFile: tt.fields.OutputFile,
				SortOnKind: tt.fields.SortOnKind,
				OutputDir:  tt.fields.OutputDir,
				Exclude:    tt.fields.Excluded,
			}
			if err := o.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Options.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOptions_Run(t *testing.T) {
	type fields struct {
		Header     string
		Paths      []string
		ValuesPath string
		Values     map[string]interface{}
		OutputFile string
		SortOnKind bool
		OutputDir  string
		Excluded   []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "header no outputdir",
			fields: fields{
				Header:     "../../../test/unit/resources/scenario/musttemplateasset/header.txt",
				Paths:      []string{"../../../test/unit/resources/scenario/musttemplateasset/body_for_header.txt"},
				ValuesPath: "../../../test/unit/resources/scenario/values.yaml",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				Header:     tt.fields.Header,
				Paths:      tt.fields.Paths,
				ValuesPath: tt.fields.ValuesPath,
				Values:     tt.fields.Values,
				OutputFile: tt.fields.OutputFile,
				SortOnKind: tt.fields.SortOnKind,
				OutputDir:  tt.fields.OutputDir,
				Exclude:    tt.fields.Excluded,
			}
			if err := o.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Options.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
