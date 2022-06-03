package workloads_test

import (
	"context"
	"errors"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/webhooks"
	"code.cloudfoundry.org/korifi/controllers/webhooks/workloads"
	"code.cloudfoundry.org/korifi/controllers/webhooks/workloads/fake"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("CFAppValidatingWebhook", func() {
	const (
		testAppGUID      = "test-app-guid"
		testAppName      = "test-app"
		testAppNamespace = "default"
	)

	var (
		ctx                context.Context
		duplicateValidator *fake.NameValidator
		cfApp              *korifiv1alpha1.CFApp
		validatingWebhook  *workloads.CFAppValidator
		retErr             error
	)

	BeforeEach(func() {
		ctx = context.Background()

		scheme := runtime.NewScheme()
		err := korifiv1alpha1.AddToScheme(scheme)
		Expect(err).NotTo(HaveOccurred())

		cfApp = &korifiv1alpha1.CFApp{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testAppGUID,
				Namespace: testAppNamespace,
			},
			Spec: korifiv1alpha1.CFAppSpec{
				DisplayName:  testAppName,
				DesiredState: korifiv1alpha1.StoppedState,
			},
		}

		duplicateValidator = new(fake.NameValidator)
		validatingWebhook = workloads.NewCFAppValidator(duplicateValidator)
	})

	Describe("ValidateCreate", func() {
		JustBeforeEach(func() {
			retErr = validatingWebhook.ValidateCreate(ctx, cfApp)
		})

		It("allows the request", func() {
			Expect(retErr).NotTo(HaveOccurred())
		})

		It("invokes the validator correctly", func() {
			Expect(duplicateValidator.ValidateCreateCallCount()).To(Equal(1))
			actualContext, _, namespace, name := duplicateValidator.ValidateCreateArgsForCall(0)
			Expect(actualContext).To(Equal(ctx))
			Expect(namespace).To(Equal(testAppNamespace))
			Expect(name).To(Equal(testAppName))
		})

		When("the app name is a duplicate", func() {
			BeforeEach(func() {
				duplicateValidator.ValidateCreateReturns(webhooks.ErrorDuplicateName)
			})

			It("denies the request", func() {
				Expect(retErr).To(MatchError(webhooks.ValidationError{
					Type:    workloads.DuplicateAppNameErrorType,
					Message: "App with the name '" + cfApp.Spec.DisplayName + "' already exists.",
				}.Marshal()))
			})
		})

		When("validating the app name fails", func() {
			BeforeEach(func() {
				duplicateValidator.ValidateCreateReturns(errors.New("boom"))
			})

			It("denies the request", func() {
				Expect(retErr).To(MatchError(webhooks.AdmissionUnknownErrorReason()))
			})
		})
	})

	Describe("ValidateUpdate", func() {
		var updatedCFApp *korifiv1alpha1.CFApp

		BeforeEach(func() {
			updatedCFApp = cfApp.DeepCopy()
			updatedCFApp.Spec.DisplayName = "the-new-name"
		})

		JustBeforeEach(func() {
			retErr = validatingWebhook.ValidateUpdate(ctx, cfApp, updatedCFApp)
		})

		It("allows the request", func() {
			Expect(retErr).NotTo(HaveOccurred())
		})

		It("invokes the validator correctly", func() {
			Expect(duplicateValidator.ValidateUpdateCallCount()).To(Equal(1))
			actualContext, _, namespace, oldName, newName := duplicateValidator.ValidateUpdateArgsForCall(0)
			Expect(actualContext).To(Equal(ctx))
			Expect(namespace).To(Equal(cfApp.Namespace))
			Expect(oldName).To(Equal(cfApp.Spec.DisplayName))
			Expect(newName).To(Equal(updatedCFApp.Spec.DisplayName))
		})

		When("the new app name is a duplicate", func() {
			BeforeEach(func() {
				duplicateValidator.ValidateUpdateReturns(webhooks.ErrorDuplicateName)
			})

			It("denies the request", func() {
				Expect(retErr).To(MatchError(webhooks.ValidationError{
					Type:    workloads.DuplicateAppNameErrorType,
					Message: "App with the name '" + updatedCFApp.Spec.DisplayName + "' already exists.",
				}.Marshal()))
			})
		})

		When("the update validation fails for another reason", func() {
			BeforeEach(func() {
				duplicateValidator.ValidateUpdateReturns(errors.New("boom!"))
			})

			It("denies the request", func() {
				Expect(retErr).To(MatchError(webhooks.AdmissionUnknownErrorReason()))
			})
		})
	})

	Describe("ValidateDelete", func() {
		JustBeforeEach(func() {
			retErr = validatingWebhook.ValidateDelete(ctx, cfApp)
		})

		It("allows the request", func() {
			Expect(retErr).NotTo(HaveOccurred())
		})

		It("invokes the validator correctly", func() {
			Expect(duplicateValidator.ValidateDeleteCallCount()).To(Equal(1))
			actualContext, _, namespace, name := duplicateValidator.ValidateDeleteArgsForCall(0)
			Expect(actualContext).To(Equal(ctx))
			Expect(namespace).To(Equal(cfApp.Namespace))
			Expect(name).To(Equal(cfApp.Spec.DisplayName))
		})

		When("delete validation fails", func() {
			BeforeEach(func() {
				duplicateValidator.ValidateDeleteReturns(errors.New("boom!"))
			})

			It("disallows the request", func() {
				Expect(retErr).To(MatchError(webhooks.AdmissionUnknownErrorReason()))
			})
		})
	})
})
