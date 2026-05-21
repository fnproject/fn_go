package shim

import (
	"testing"

	"github.com/oracle/oci-go-sdk/v65/functions"
)

func TestGeneratedRawDestinationAnnotationsArePreservedWhenFriendlyAnnotationsAbsent(t *testing.T) {
	annotations := map[string]interface{}{
		annotationOCIParityFnSuccessDestination: map[string]interface{}{
			"kind":     "STREAM",
			"streamId": "ocid1.stream.oc1..example",
		},
		annotationOCIParityFnFailureDestination: map[string]interface{}{
			"kind":    "NOTIFICATION",
			"topicId": "ocid1.onstopic.oc1..example",
		},
	}

	createDetails := functions.CreateFunctionDetails{}
	if err := applyGeneratedOCIParityCreateFunctionDetails(&createDetails, annotations); err != nil {
		t.Fatalf("applyGeneratedOCIParityCreateFunctionDetails() error = %v", err)
	}
	if successDestination, failureDestination, err := parseDestinationAnnotations(annotations); err != nil {
		t.Fatalf("parseDestinationAnnotations() error = %v", err)
	} else {
		if successDestination != nil {
			createDetails.SuccessDestination = successDestination
		}
		if failureDestination != nil {
			createDetails.FailureDestination = failureDestination
		}
	}
	if createDetails.SuccessDestination == nil {
		t.Fatal("expected generated success destination to remain set on create details")
	}
	if createDetails.FailureDestination == nil {
		t.Fatal("expected generated failure destination to remain set on create details")
	}

	updateDetails := functions.UpdateFunctionDetails{}
	if err := applyGeneratedOCIParityUpdateFunctionDetails(&updateDetails, annotations); err != nil {
		t.Fatalf("applyGeneratedOCIParityUpdateFunctionDetails() error = %v", err)
	}
	if successDestination, failureDestination, err := parseDestinationAnnotations(annotations); err != nil {
		t.Fatalf("parseDestinationAnnotations() error = %v", err)
	} else {
		if successDestination != nil {
			updateDetails.SuccessDestination = successDestination
		}
		if failureDestination != nil {
			updateDetails.FailureDestination = failureDestination
		}
	}
	if updateDetails.SuccessDestination == nil {
		t.Fatal("expected generated success destination to remain set on update details")
	}
	if updateDetails.FailureDestination == nil {
		t.Fatal("expected generated failure destination to remain set on update details")
	}
}
