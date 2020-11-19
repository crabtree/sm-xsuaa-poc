package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Peripli/service-manager-cli/pkg/httputil"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/crabtree/sm-xsuaa-poc/internal/params"
	"github.com/google/uuid"
)

const (
	offeringName = "xsuaa"
	planName     = "application"
)

func main() {
	// as params pass username, password, and service manager URL
	// APP_USERNAME=username APP_PASSWORD=passwd APP_BASE_URL=https://smurl.domain.local go run main.go"
	params, err := params.Parse()
	exitOnError(err)

	// reusing the client from smctl, but we use http client with
	// transport supporting basic auth credentials to be able to
	// use credentials from the platform
	httpClient := &http.Client{}
	httpClient.Transport = &BasicAuthTransport{
		Username: params.Username,
		Password: params.Password,
		Rt:       http.DefaultTransport,
	}

	sm := smclient.NewClient(context.Background(), httpClient, params.BaseURL)

	// list the offergins filtering by the xsuaa name
	offerings, err := sm.ListOfferings(&query.Parameters{FieldQuery: []string{fmt.Sprintf("name eq '%s'", offeringName)}})
	exitOnError(err)

	if len(offerings.ServiceOfferings) == 0 {
		log.Fatalf("There is no offering named %s\n", offeringName)
	}
	offering := offerings.ServiceOfferings[0]

	// list the plans filtering by th application plan
	plans, err := sm.ListPlans(&query.Parameters{FieldQuery: []string{
		fmt.Sprintf("name eq '%s'", planName),
		fmt.Sprintf("service_offering_id eq '%s'", offering.ID)}})
	exitOnError(err)

	if len(plans.ServicePlans) == 0 {
		log.Fatalf("There is no plan named %s for offering %s\n", planName, offeringName)
	}
	plan := plans.ServicePlans[0]

	// creating the xsuaa service instance
	id := uuid.New().String()
	spaceID := uuid.New().String()
	orgID := uuid.New().String()
	instanceBody := map[string]interface{}{
		"service_id":        offering.CatalogID,
		"plan_id":           plan.CatalogID,
		"space_id":          spaceID,
		"organization_guid": orgID,
		"context": map[string]string{
			"platform": "kubernetes",
		},
	}
	data, err := json.Marshal(instanceBody)
	exitOnError(err)

	// as we are using the platform credentials
	// the request must be issued directly to the OSB API
	// the smctl client does not expose it as a method,
	// so we use generic Call
	res, err := sm.Call("PUT",
		fmt.Sprintf("/v1/osb/%s/v2/service_instances/%s", offering.BrokerID, id),
		bytes.NewBuffer(data),
		&query.Parameters{})
	exitOnError(err)

	if res.StatusCode >= 300 {
		log.Fatalf("Got status code %d when creating service instance", res.StatusCode)
	}

	if res.StatusCode == 202 {
		// queued, do the status check
	} else if res.StatusCode == 201 {
		// created, we can get the last operation
		res, err := sm.Call("GET",
			fmt.Sprintf("/v1/osb/%s/v2/service_instances/%s/last_operation", offering.BrokerID, id),
			nil,
			&query.Parameters{})
		exitOnError(err)

		var lastOp map[string]interface{}
		err = httputil.UnmarshalResponse(res, &lastOp)
		exitOnError(err)

		// create binding
		bindingID := uuid.New().String()
		bindingBody := map[string]interface{}{
			"service_id": offering.CatalogID,
			"plan_id":    plan.CatalogID,
			"context": map[string]string{
				"platform": "kubernetes",
			},
		}
		data, err := json.Marshal(bindingBody)
		exitOnError(err)

		res, err = sm.Call("PUT",
			fmt.Sprintf("/v1/osb/%s/v2/service_instances/%s/service_bindings/%s", offering.BrokerID, id, bindingID),
			bytes.NewBuffer(data),
			&query.Parameters{})
		exitOnError(err)

		var bindingOp map[string]interface{}
		err = httputil.UnmarshalResponse(res, &bindingOp)
		exitOnError(err)

		// print the binding details
		fmt.Printf("%+v\n", bindingOp)
	} else {
		log.Fatalf("Got status code %d", res.StatusCode)
	}
}

func exitOnError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

type BasicAuthTransport struct {
	Username string
	Password string

	Rt http.RoundTripper
}

func (b *BasicAuthTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if b.Username != "" && b.Password != "" {
		request.SetBasicAuth(b.Username, b.Password)
	}

	return b.Rt.RoundTrip(request)
}
