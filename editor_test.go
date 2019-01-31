package prompteditor

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type testStruct struct {
	String string   `json:"string"`
	Number int      `json:"number"`
	Slice  []string `json:"slice"`
}

func TestNew(t *testing.T) {
	type args struct {
		editor   string
		original OriginalData
	}

	cases := map[string]struct {
		args    args
		wantErr bool
	}{
		"normal": {
			args: args{
				editor: "vim",
				original: MarshalJSON(&testStruct{
					String: "test",
					Number: 999,
					Slice:  []string{"ele1", "ele2"},
				}),
			},
			wantErr: false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := New(c.args.editor, c.args.original)
			if (err != nil) != c.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}

func TestEditor_Open_UnmarshalEdited(t *testing.T) {
	v := &testStruct{
		String: "test",
		Number: 999,
		Slice:  []string{"ele1", "ele2"},
	}

	e, err := New("vim", MarshalJSON(v))
	if err != nil {
		t.Fatal(e)
	}

	cases := map[string]struct {
		editor     *Editor
		editedBody string
		wantErr    bool
		want       *testStruct
	}{
		"normal": {
			editor:     e,
			editedBody: `{"string": "test edited", "number": 1, "slice": ["ele1", "ele2", "ele3"]}`,
			wantErr:    false,
			want: &testStruct{
				String: "test edited",
				Number: 1,
				Slice:  []string{"ele1", "ele2", "ele3"},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			edit = func(editor string, file *os.File) error {
				return ioutil.WriteFile(file.Name(), []byte(c.editedBody), 0644)
			}

			if err := c.editor.Open(); (err != nil) != c.wantErr {
				t.Errorf("Editor.Open() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if c.wantErr {
				return
			}

			var edited testStruct
			if err := c.editor.UnmarshalEdited(&edited); err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(&edited, c.want) {
				t.Errorf("after Editor.Open(), expect edited = %v, actual = %v", &edited, c.want)
			}
		})
	}
}
