package template

import (
	"bytes"
	"context"
	"fmt"
	"github.com/helmwave/helmwave/pkg/parallel"
	"io"
	"sync"
	"text/template"
)

type Template struct {
	templates []*template.Template
}

func New() *Template {
	return &Template{}
}

func (t *Template) AddFiles(files []string) error {
	if t.templates == nil {
		t.templates = make([]*template.Template, 0, len(files))
	}
	for _, f := range files {
		tmpl := template.New(f)
		_, err := tmpl.ParseFiles(f)
		if err != nil {
			return fmt.Errorf("failed to parse file %q: %w", f, err)
		}

		t.templates = append(t.templates, tmpl)
	}

	return nil
}

func (t *Template) Run(out io.Writer) error {
	ctx := context.TODO()
	wg := parallel.NewWaitGroup()
	wg.Add(len(t.templates))

	mutex := &sync.Mutex{}

	for _, tmpl := range t.templates {
		go func(mu *sync.Mutex, wg *parallel.WaitGroup, tmpl *template.Template, out io.Writer) {
			defer wg.Done()
			buffer := &bytes.Buffer{}
			buffer.WriteString("\n---\n")
			buffer.WriteString(fmt.Sprintf("# templated out of %q", tmpl.Name))

			err := tmpl.Execute(buffer, nil)
			if err != nil {
				wg.ErrChan() <- err
				return
			}

			mu.Lock()
			defer mu.Unlock()
			_, err = io.Copy(out, buffer)
			wg.ErrChan() <- err
		}(mutex, wg, tmpl, out)
	}

	return wg.WaitWithContext(ctx)
}
