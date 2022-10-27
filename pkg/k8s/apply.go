package k8s

import (
	"bytes"
	"fmt"
	"io"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/delete"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func ApplyManifest(manifest string) error {
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	applyFlags := apply.NewApplyFlags(f, genericclioptions.IOStreams{
		In:     bytes.NewBuffer([]byte{}),
		Out:    io.Discard,
		ErrOut: io.Discard,
	})

	applyOptions, err := GetApplyOptions(applyFlags, manifest)
	if err != nil {
		return err
	}

	err = applyOptions.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate: %w", err)
	}

	err = applyOptions.Run()
	if err != nil {
		return fmt.Errorf("failed to apply: %w", err)
	}

	return nil
}

func GetApplyOptions(flags *apply.ApplyFlags, manifest string) (*apply.ApplyOptions, error) {
	dynamicClient, err := flags.Factory.DynamicClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get dynamic client: %w", err)
	}

	// dryRunVerifier := resource.NewQueryParamVerifier(dynamicClient, flags.Factory.OpenAPIGetter(), resource.QueryParamDryRun)
	fieldValidationVerifier := resource.NewQueryParamVerifier(dynamicClient, flags.Factory.OpenAPIGetter(), resource.QueryParamFieldValidation)
	fieldManager := "chartwave"

	openAPISchema, _ := flags.Factory.OpenAPISchema()

	// validationDirective, err := cmdutil.GetValidationDirective(cmd)
	// if err != nil {
	// 	return nil, err
	// }
	validationDirective := metav1.FieldValidationStrict

	validator, err := flags.Factory.Validator(validationDirective, fieldValidationVerifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get validator: %w", err)
	}

	builder := flags.Factory.NewBuilder()
	mapper, err := flags.Factory.ToRESTMapper()
	if err != nil {
		return nil, fmt.Errorf("failed to get REST mapper: %w", err)
	}

	namespace, enforceNamespace, err := flags.Factory.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	o := &apply.ApplyOptions{
		PrintFlags: flags.PrintFlags,

		DeleteOptions: &delete.DeleteOptions{
			FilenameOptions: resource.FilenameOptions{
				Filenames: []string{manifest},
			},
		},
		ToPrinter: dummyPrinterGetter,
		// ServerSideApply: serverSideApply,
		// ForceConflicts:  forceConflicts,
		FieldManager: fieldManager,
		Selector:     flags.Selector,
		// DryRunStrategy:  dryRunStrategy,
		// DryRunVerifier:  dryRunVerifier,
		// Prune:           flags.Prune,
		// PruneResources:  flags.PruneResources,
		All:            flags.All,
		Overwrite:      flags.Overwrite,
		OpenAPIPatch:   flags.OpenAPIPatch,
		PruneWhitelist: flags.PruneWhitelist,

		Recorder:            genericclioptions.NoopRecorder{},
		Namespace:           namespace,
		EnforceNamespace:    enforceNamespace,
		Validator:           validator,
		ValidationDirective: validationDirective,
		Builder:             builder,
		Mapper:              mapper,
		DynamicClient:       dynamicClient,
		OpenAPISchema:       openAPISchema,

		IOStreams: flags.IOStreams,

		// objects:       []*resource.Info{},
		// objectsCached: false,

		VisitedUids:       sets.NewString(),
		VisitedNamespaces: sets.NewString(),
	}

	o.PostProcessorFn = o.PrintAndPrunePostProcessor()

	return o, nil
}
