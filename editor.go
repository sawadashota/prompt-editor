package prompteditor

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// Editor .
type Editor struct {
	// editor command
	editor string

	// originalBytes body before editing
	originalBytes []byte

	// editedBytes body after editing
	editedBytes []byte
}

// OriginalData .
type OriginalData func(e *Editor) error

// New Editor
func New(editor string, original OriginalData) (*Editor, error) {
	e := &Editor{
		editor: editor,
	}

	if err := original(e); err != nil {
		return nil, err
	}

	return e, nil
}

// MarshalJSON form struct
func MarshalJSON(v interface{}) OriginalData {
	return func(e *Editor) error {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}

		var indented bytes.Buffer
		if err := json.Indent(&indented, b, "", "  "); err != nil {
			return err
		}

		e.originalBytes = indented.Bytes()
		return nil
	}
}

// edit file using editor
var edit = func(editor string, file *os.File) error {
	cmd := exec.Command(editor, file.Name())

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

// TempFileName .
var TempFileName = "prompteditor"

// Open editor to edit
func (e *Editor) Open() error {
	tempFile, err := ioutil.TempFile("", TempFileName)
	if err != nil {
		return errors.Wrap(err, "create tempfile")
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(e.originalBytes); err != nil {
		return errors.Wrap(err, "write tempfile")
	}

	if err := edit(e.editor, tempFile); err != nil {
		return errors.Wrap(err, "edit")
	}

	body, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return errors.Wrap(err, "read edited file")
	}

	e.editedBytes = body
	return nil
}

// UnmarshalEdited object to given arg
func (e *Editor) UnmarshalEdited(v interface{}) error {
	if len(e.editedBytes) == 0 {
		return errors.New("editedBytes is empty")
	}

	if err := json.Unmarshal(e.editedBytes, v); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	return nil
}
