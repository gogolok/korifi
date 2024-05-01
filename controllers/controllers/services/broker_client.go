package services

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BrokerClient interface {
	GetCatalog(context.Context, *korifiv1alpha1.CFServiceBroker) (*Catalog, error)

	ProvisionServiceInstance(context.Context, *korifiv1alpha1.CFServiceInstance) error
	GetServiceInstanceLastOperation(context.Context, *korifiv1alpha1.CFServiceInstance) (LastOperation, error)
	DeprovisionServiceInstance(context.Context, *korifiv1alpha1.CFServiceInstance) error

	BindService(context.Context, *korifiv1alpha1.CFServiceBinding) error
	GetServiceBinding(context.Context, *korifiv1alpha1.CFServiceBinding) (ServiceBinding, error)
	GetServiceBindingLastOperation(context.Context, *korifiv1alpha1.CFServiceBinding) (LastOperation, error)
	UnbindService(context.Context, *korifiv1alpha1.CFServiceBinding) error
}

type Catalog struct {
	Services []Service `json:"services"`
}

type Service struct {
	Id                   string                 `json:"id"`
	Name                 string                 `json:"name"`
	Description          string                 `json:"description"`
	Bindable             bool                   `json:"bindable"`
	InstancesRetrievable bool                   `json:"instances_retrievable"`
	BindingsRetrievable  bool                   `json:"bindings_retrievable"`
	PlanUpdateable       bool                   `json:"plan_updateable"`
	AllowContextUpdates  bool                   `json:"allow_context_updates"`
	Tags                 []string               `json:"tags"`
	Requires             []string               `json:"requires"`
	Metadata             map[string]interface{} `json:"metadata"`
	DashboardClient      struct {
		Id          string `json:"id"`
		Secret      string `json:"secret"`
		RedirectUri string `json:"redirect_url"`
	} `json:"dashboard_client"`
	Plans []Plan `json:"plans"`
}

type Plan struct {
	Id               string                            `json:"id"`
	Name             string                            `json:"name"`
	Description      string                            `json:"description"`
	Metadata         map[string]interface{}            `json:"metadata"`
	Free             bool                              `json:"free"`
	Bindable         bool                              `json:"bindable"`
	BindingRotatable bool                              `json:"binding_rotatable"`
	PlanUpdateable   bool                              `json:"plan_updateable"`
	Schemas          korifiv1alpha1.ServicePlanSchemas `json:"schemas"`
}

type ServiceBinding struct {
	Credentials map[string]any `json:"credentials"`
}

type LastOperation struct {
	Exists      bool
	State       string
	Description string
}

type brokerClient struct {
	k8sClient     client.Client
	rootNamespace string
}

func NewBrokerClient(
	k8sClient client.Client,
	rootNamespace string,
) BrokerClient {
	return &brokerClient{
		k8sClient:     k8sClient,
		rootNamespace: rootNamespace,
	}
}

func (c *brokerClient) GetCatalog(ctx context.Context, broker *korifiv1alpha1.CFServiceBroker) (*Catalog, error) {
	_, resp, err := c.newBrokerRequester().forBroker(broker).sendRequest(ctx, "/v2/catalog", http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("get catalog request failed: %w", err)
	}

	catalog := &Catalog{}
	err = json.Unmarshal(resp, catalog)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal catalog: %w", err)
	}

	return catalog, nil
}

func (c *brokerClient) ProvisionServiceInstance(ctx context.Context, cfServiceInstance *korifiv1alpha1.CFServiceInstance) error {
	plan, err := c.getServicePlan(ctx, cfServiceInstance)
	if err != nil {
		return err
	}

	offering, err := c.getServiceOffering(ctx, plan)
	if err != nil {
		return err
	}

	provisionRequest := map[string]any{
		"service_id": offering.Spec.Broker_catalog.Id,
		"plan_id":    plan.Spec.Broker_catalog.Id,
	}

	if cfServiceInstance.Spec.Parameters != nil {
		paramsMap := map[string]any{}
		err = json.Unmarshal(cfServiceInstance.Spec.Parameters.Raw, &paramsMap)
		if err != nil {
			return fmt.Errorf("failed to unmarshal service instance parameters: %w", err)
		}
		provisionRequest["parameters"] = paramsMap
	}

	provisionUrl, err := url.JoinPath("v2", "service_instances", cfServiceInstance.Name)
	if err != nil {
		return fmt.Errorf("failed to construct service broker provision url: %w", err)
	}

	_, _, err = c.newBrokerRequester().forPlan(plan).async().sendRequest(ctx, provisionUrl, http.MethodPut, provisionRequest)
	if err != nil {
		return fmt.Errorf("provision request failed: %w", err)
	}

	return nil
}

func (c *brokerClient) GetServiceInstanceLastOperation(ctx context.Context, cfServiceInstance *korifiv1alpha1.CFServiceInstance) (LastOperation, error) {
	plan, err := c.getServicePlan(ctx, cfServiceInstance)
	if err != nil {
		return LastOperation{}, err
	}

	stateUrl, err := url.JoinPath("v2", "service_instances", cfServiceInstance.Name, "last_operation")
	if err != nil {
		return LastOperation{}, fmt.Errorf("failed to construct service broker last operation url: %w", err)
	}

	respCode, resp, err := c.newBrokerRequester().
		forPlan(plan).
		allowNotFound().
		sendRequest(ctx, stateUrl, http.MethodGet, nil)
	if err != nil {
		return LastOperation{}, fmt.Errorf("last operation request failed: %w", err)
	}

	if respCode == http.StatusNotFound {
		return LastOperation{
			Exists: false,
		}, nil
	}

	respMap := map[string]string{}
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return LastOperation{}, fmt.Errorf("failed to unmarshal last operation response: %w", err)
	}

	return LastOperation{
		Exists:      true,
		State:       respMap["state"],
		Description: respMap["description"],
	}, nil
}

func (c *brokerClient) DeprovisionServiceInstance(ctx context.Context, cfServiceInstance *korifiv1alpha1.CFServiceInstance) error {
	plan, err := c.getServicePlan(ctx, cfServiceInstance)
	if err != nil {
		return err
	}

	deprovisionUrl, err := url.JoinPath("v2", "service_instances", cfServiceInstance.Name)
	if err != nil {
		return fmt.Errorf("failed to construct service broker deprovision url: %w", err)
	}

	_, _, err = c.newBrokerRequester().forPlan(plan).async().sendRequest(ctx, deprovisionUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("deprovision request failed: %w", err)
	}

	return nil
}

func (c *brokerClient) BindService(ctx context.Context, cfServiceBinding *korifiv1alpha1.CFServiceBinding) error {
	cfServiceInstance, err := c.getCFServiceInstance(ctx, cfServiceBinding)
	if err != nil {
		return err
	}
	plan, err := c.getServicePlan(ctx, cfServiceInstance)
	if err != nil {
		return err
	}

	offering, err := c.getServiceOffering(ctx, plan)
	if err != nil {
		return err
	}

	bindRequest := map[string]any{
		"service_id": offering.Spec.Broker_catalog.Id,
		"plan_id":    plan.Spec.Broker_catalog.Id,
	}

	bindUrl, err := url.JoinPath("v2", "service_instances", cfServiceInstance.Name, "service_bindings", cfServiceBinding.Name)
	if err != nil {
		return fmt.Errorf("failed to construct service broker provision url: %w", err)
	}

	_, _, err = c.newBrokerRequester().forPlan(plan).async().sendRequest(ctx, bindUrl, http.MethodPut, bindRequest)
	if err != nil {
		return fmt.Errorf("bind request failed: %w", err)
	}

	return nil
}

func (c *brokerClient) GetServiceBinding(ctx context.Context, cfServiceBinding *korifiv1alpha1.CFServiceBinding) (ServiceBinding, error) {
	cfServiceInstance, err := c.getCFServiceInstance(ctx, cfServiceBinding)
	if err != nil {
		return ServiceBinding{}, err
	}
	plan, err := c.getServicePlan(ctx, cfServiceInstance)
	if err != nil {
		return ServiceBinding{}, err
	}

	offering, err := c.getServiceOffering(ctx, plan)
	if err != nil {
		return ServiceBinding{}, err
	}

	getRequest := map[string]any{
		"service_id": offering.Spec.Broker_catalog.Id,
		"plan_id":    plan.Spec.Broker_catalog.Id,
	}

	getUrl, err := url.JoinPath("v2", "service_instances", cfServiceInstance.Name, "service_bindings", cfServiceBinding.Name)
	if err != nil {
		return ServiceBinding{}, fmt.Errorf("failed to construct service broker binding url: %w", err)
	}

	_, resp, err := c.newBrokerRequester().forPlan(plan).sendRequest(ctx, getUrl, http.MethodGet, getRequest)
	if err != nil {
		return ServiceBinding{}, fmt.Errorf("get binding request failed: %w", err)
	}

	serviceBinding := ServiceBinding{}
	err = json.Unmarshal(resp, &serviceBinding)
	if err != nil {
		return ServiceBinding{}, fmt.Errorf("failed to unmarshal binding response: %w", err)
	}
	return serviceBinding, nil
}

func (c *brokerClient) GetServiceBindingLastOperation(ctx context.Context, cfServiceBinding *korifiv1alpha1.CFServiceBinding) (LastOperation, error) {
	cfServiceInstance, err := c.getCFServiceInstance(ctx, cfServiceBinding)
	if err != nil {
		return LastOperation{}, err
	}

	plan, err := c.getServicePlan(ctx, cfServiceInstance)
	if err != nil {
		return LastOperation{}, err
	}

	stateUrl, err := url.JoinPath("v2", "service_instances", cfServiceInstance.Name, "service_bindings", cfServiceBinding.Name, "last_operation")
	if err != nil {
		return LastOperation{}, fmt.Errorf("failed to construct service broker last operation url: %w", err)
	}

	respCode, resp, err := c.newBrokerRequester().
		forPlan(plan).
		allowNotFound().
		sendRequest(ctx, stateUrl, http.MethodGet, nil)
	if err != nil {
		return LastOperation{}, fmt.Errorf("last operation request failed: %w", err)
	}

	if respCode == http.StatusNotFound {
		return LastOperation{
			Exists: false,
		}, nil
	}

	respMap := map[string]string{}
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return LastOperation{}, fmt.Errorf("failed to unmarshal last operation response: %w", err)
	}

	return LastOperation{
		Exists:      true,
		State:       respMap["state"],
		Description: respMap["description"],
	}, nil
}

func (c *brokerClient) UnbindService(ctx context.Context, cfServiceBinding *korifiv1alpha1.CFServiceBinding) error {
	cfServiceInstance, err := c.getCFServiceInstance(ctx, cfServiceBinding)
	if err != nil {
		return err
	}

	plan, err := c.getServicePlan(ctx, cfServiceInstance)
	if err != nil {
		return err
	}

	unbindUrl, err := url.JoinPath("v2", "service_instances", cfServiceInstance.Name, "service_bindings", cfServiceBinding.Name)
	if err != nil {
		return fmt.Errorf("failed to construct service broker unbind url: %w", err)
	}

	_, _, err = c.newBrokerRequester().forPlan(plan).async().sendRequest(ctx, unbindUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("unbind request failed: %w", err)
	}

	return nil
}

func payloadToReader(payload map[string]any) (io.Reader, error) {
	if len(payload) == 0 {
		return nil, nil
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return bytes.NewBuffer(payloadBytes), nil
}

func (c *brokerClient) getServicePlan(ctx context.Context, cfServiceInstance *korifiv1alpha1.CFServiceInstance) (*korifiv1alpha1.CFServicePlan, error) {
	plan := &korifiv1alpha1.CFServicePlan{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: c.rootNamespace,
			Name:      cfServiceInstance.Spec.ServicePlanGUID,
		},
	}

	err := c.k8sClient.Get(ctx, client.ObjectKeyFromObject(plan), plan)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan. rootNs: %q; planGUID %q", c.rootNamespace, cfServiceInstance.Spec.ServicePlanGUID)
	}

	return plan, nil
}

func (c *brokerClient) getServiceOffering(ctx context.Context, plan *korifiv1alpha1.CFServicePlan) (*korifiv1alpha1.CFServiceOffering, error) {
	offering := &korifiv1alpha1.CFServiceOffering{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: c.rootNamespace,
			Name:      plan.Labels[korifiv1alpha1.RelServiceOfferingLabel],
		},
	}

	err := c.k8sClient.Get(ctx, client.ObjectKeyFromObject(offering), offering)
	if err != nil {
		return nil, err
	}

	return offering, nil
}

func (c *brokerClient) getCFServiceInstance(ctx context.Context, cfServiceBinding *korifiv1alpha1.CFServiceBinding) (*korifiv1alpha1.CFServiceInstance, error) {
	cfServiceInstance := &korifiv1alpha1.CFServiceInstance{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cfServiceBinding.Namespace,
			Name:      cfServiceBinding.Spec.Service.Name,
		},
	}

	err := c.k8sClient.Get(ctx, client.ObjectKeyFromObject(cfServiceInstance), cfServiceInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to get CFServiceInstance for binding: %w", err)
	}

	return cfServiceInstance, nil
}

type brokerRequester struct {
	k8sClient           client.Client
	rootNamespace       string
	plan                *korifiv1alpha1.CFServicePlan
	broker              *korifiv1alpha1.CFServiceBroker
	acceptsIncomplete   bool
	allowNotFoundStatus bool
}

func (c *brokerClient) newBrokerRequester() *brokerRequester {
	return &brokerRequester{k8sClient: c.k8sClient, rootNamespace: c.rootNamespace}
}

func (r *brokerRequester) forPlan(plan *korifiv1alpha1.CFServicePlan) *brokerRequester {
	r.plan = plan
	return r
}

func (r *brokerRequester) forBroker(broker *korifiv1alpha1.CFServiceBroker) *brokerRequester {
	r.broker = broker
	return r
}

func (r *brokerRequester) async() *brokerRequester {
	r.acceptsIncomplete = true
	return r
}

func (r *brokerRequester) allowNotFound() *brokerRequester {
	r.allowNotFoundStatus = true
	return r
}

func (r *brokerRequester) sendRequest(ctx context.Context, requestPath string, method string, payload map[string]any) (int, []byte, error) {
	broker, err := r.getBroker(ctx)
	if err != nil {
		return 0, nil, err
	}

	requestUrl, err := url.JoinPath(broker.Spec.URL, requestPath)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to build broker requestUrl for path %q: %w", requestPath, err)
	}
	if r.acceptsIncomplete {
		requestUrl += "?" + url.Values{"accepts_incomplete": {"true"}}.Encode()
	}

	payloadReader, err := payloadToReader(payload)
	if err != nil {
		return 0, nil, fmt.Errorf("failed create payload reader: %w", err)
	}

	req, err := http.NewRequest(method, requestUrl, payloadReader)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create new HTTP request: %w", err)
	}

	// TODO: configure whether to trust insecure brokers
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	userName, password, err := r.getCredentials(ctx, broker)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get credentials: %w", err)
	}
	authPlain := fmt.Sprintf("%s:%s", userName, password)
	auth := base64.StdEncoding.EncodeToString([]byte(authPlain))
	req.Header.Add("Authorization", "Basic "+auth)

	req.Header.Add("X-Broker-API-Version", "2.17")

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read body: %w", err)
	}

	if r.allowNotFoundStatus && resp.StatusCode == http.StatusNotFound {
		return resp.StatusCode, respBody, nil
	}

	if resp.StatusCode > 299 {
		return resp.StatusCode, respBody, fmt.Errorf("request returned non-OK status %d: %s", resp.StatusCode, string(respBody))
	}

	return resp.StatusCode, respBody, nil
}

func (r *brokerRequester) getBroker(ctx context.Context) (*korifiv1alpha1.CFServiceBroker, error) {
	if r.broker != nil {
		return r.broker, nil
	}

	if r.plan == nil {
		return nil, errors.New("neither broker, nor plan configured")
	}

	brokerName, ok := r.plan.Labels[korifiv1alpha1.RelServiceBrokerLabel]
	if !ok {
		return nil, fmt.Errorf("plan %q has no broker guid label set", r.plan.Name)
	}

	broker := &korifiv1alpha1.CFServiceBroker{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: r.rootNamespace,
			Name:      brokerName,
		},
	}

	err := r.k8sClient.Get(ctx, client.ObjectKeyFromObject(broker), broker)
	if err != nil {
		return nil, fmt.Errorf("failed to get broker: %w", err)
	}

	return broker, nil
}

func (r *brokerRequester) getCredentials(ctx context.Context, broker *korifiv1alpha1.CFServiceBroker) (string, string, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: broker.Namespace,
			Name:      broker.Spec.SecretName,
		},
	}
	err := r.k8sClient.Get(ctx, client.ObjectKeyFromObject(secret), secret)
	if err != nil {
		return "", "", err
	}

	return string(secret.Data["username"]), string(secret.Data["password"]), nil
}
