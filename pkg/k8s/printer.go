package k8s

import (
	"io"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
)

type dummyPrinter struct{}

func (dummyPrinter) PrintObj(runtime.Object, io.Writer) error {
	return nil
}

func dummyPrinterGetter(string) (printers.ResourcePrinter, error) {
	return dummyPrinter{}, nil
}
