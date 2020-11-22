package graph_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/imroc/req"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/feed/graph"
	"gitlab.slade360emr.com/go/feed/graph/feed"
	db "gitlab.slade360emr.com/go/feed/graph/feed/infrastructure/database"
	"gitlab.slade360emr.com/go/feed/graph/feed/infrastructure/messaging"
	"google.golang.org/api/idtoken"
)

const (
	testHTTPClientTimeout = 180
	intMax                = 9007199254740990
)

// these are set up once in TestMain and used by all the acceptance tests in
// this package
var srv *http.Server
var baseURL string
var serverErr error

func TestMain(m *testing.M) {
	// setup
	ctx := context.Background()
	srv, baseURL, serverErr = startTestServer(ctx) // set the globals
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
	}

	// run the tests
	code := m.Run()

	// cleanup here
	defer func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("test server shutdown error: %s", err)
		}
	}()
	os.Exit(code)
}

func TestRouter(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case - should succeed",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := graph.Router(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Router() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestHealthStatusCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	wr := httptest.NewRecorder()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful health check",
			args: args{
				w: wr,
				r: req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph.HealthStatusCheck(tt.args.w, tt.args.r)
		})
	}
}

func TestGQLHandler(t *testing.T) {
	ctx := context.Background()

	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize firebase repository: %s", err)
		return
	}

	if fr == nil {
		t.Errorf("nil firebase repository")
		return
	}

	projectID, err := base.GetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		t.Errorf("project ID not found in env var: %s", err)
		return
	}

	if projectID == "" {
		t.Errorf("nil project ID")
		return
	}

	projectNumber, err := base.GetEnvVar(base.GoogleProjectNumberEnvVarName)
	if err != nil {
		t.Errorf("project number not found in env var: %s", err)
		return
	}

	if projectNumber == "" {
		t.Errorf("nil project number")
		return
	}

	projectNumberInt, err := strconv.Atoi(projectNumber)
	if err != nil {
		t.Errorf("non int project number: %s", err)
		return
	}

	if projectNumberInt == 0 {
		t.Errorf("the project number cannot be zero")
		return
	}

	ns, err := messaging.NewPubSubNotificationService(ctx, projectID, projectNumberInt)
	if err != nil {
		t.Errorf("can't initialize notification service: %s", err)
		return
	}
	if ns == nil {
		t.Errorf("nil notification service")
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	wr := httptest.NewRecorder()

	type args struct {
		ctx context.Context
		fr  feed.Repository
		ns  feed.NotificationService
		w   *httptest.ResponseRecorder
		r   *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful setup",
			args: args{
				ctx: ctx,
				fr:  fr,
				ns:  ns,
				w:   wr,
				r:   req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := graph.GQLHandler(tt.args.ctx, tt.args.fr, tt.args.ns)
			handler.ServeHTTP(tt.args.w, tt.args.r)
			assert.Equal(t, http.StatusBadRequest, tt.args.w.Code) // no auth
		})
	}
}

func TestGraphQLProcessEvent(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation ProcessEvent($flavour: Flavour!, $event: EventInput!) {
		processEvent(flavour: $flavour, event: $event)
	}	  
	`

	gql["variables"] = map[string]interface{}{
		"flavour": "CONSUMER",
		"event": map[string]interface{}{
			"name": "TEST_EVENT",
			"context": map[string]string{
				"userID":         "user-1",
				"organizationID": "org-1",
				"locationID":     "location-1",
				"timestamp":      "2020-11-05T03:26:15+00:00",
			},
			"payload": map[string]interface{}{
				"data": map[string]interface{}{
					"some": "stuff",
					"and":  "other stuff",
				},
			},
		},
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGraphQLDeleteMessage(t *testing.T) {
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post item: %s", err)
		return
	}

	testMessage := getTestMessage()
	err = postMessage(ctx, t, uid, fl, &testMessage, baseURL, testItem.ID)
	if err != nil {
		t.Errorf("can't post message: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation DeleteMessage(
		$flavour: Flavour!, 
		$itemID: String!, 
		$messageID: String!
	  ) {
		deleteMessage(flavour: $flavour, itemID: $itemID, messageID: $messageID)
	}
	`
	gql["variables"] = map[string]interface{}{
		"flavour":   fl.String(),
		"itemID":    testItem.ID,
		"messageID": testMessage.ID,
	}
	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGraphQLPostMessage(t *testing.T) {
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post item: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation PostMessage(
		$flavour: Flavour!, 
		$itemID: String!,
		$message: MsgInput!
	  ) {
		postMessage(flavour: $flavour, itemID: $itemID, message: $message) {
		  id
		  sequenceNumber
		  text
		  replyTo
		  postedByUID
		  postedByName
		}
	}
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
		"message": map[string]string{
			"id":             ksuid.New().String(),
			"text":           ksuid.New().String(),
			"replyTo":        ksuid.New().String(),
			"postedByUID":    ksuid.New().String(),
			"postedByName":   ksuid.New().String(),
			"timestamp":      time.Now().Format(time.RFC3339),
			"sequenceNumber": fmt.Sprintf("%d", time.Now().Unix()),
		},
	}
	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGraphQLHideNudge(t *testing.T) {
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feed.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation HideNudge($flavour: Flavour!, $nudgeID: String!) {
		hideNudge(flavour: $flavour, nudgeID: $nudgeID) {
		  id
		  sequenceNumber
		  visibility
		  status
		  title
		  text
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  groups
		  users
		  links {
			id
			url
			linkType
		  }
		  notificationChannels
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"nudgeID": testNudge.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGraphQLShowNudge(t *testing.T) {
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feed.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation ShowNudge($flavour: Flavour!, $nudgeID: String!) {
		showNudge(flavour: $flavour, nudgeID: $nudgeID) {
		  id
		  sequenceNumber
		  visibility
		  status
		  title
		  text
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  groups
		  users
		  links {
			id
			url
			linkType
		  }
		  notificationChannels
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"nudgeID": testNudge.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGraphQLResolveFeedItem(t *testing.T) {
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation ResolveFeedItem($flavour: Flavour!, $itemID: String!) {
		resolveFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGraphQLUnresolveFeedItem(t *testing.T) {
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation UnresolveFeedItem($flavour: Flavour!, $itemID: String!) {
		unresolveFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
func TestGraphQLPinFeedItem(t *testing.T) {
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation PinFeedItem($flavour: Flavour!, $itemID: String!) {
		pinFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
func TestGraphQLUnpinFeedItem(t *testing.T) {
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation UnpinFeedItem($flavour: Flavour!, $itemID: String!) {
		unpinFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
func TestGraphQLHideFeedItem(t *testing.T) {
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation HideFeedItem($flavour: Flavour!, $itemID: String!) {
		hideFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
func TestGraphQLShowFeedItem(t *testing.T) {
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation ShowFeedItem($flavour: Flavour!, $itemID: String!) {
		showFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGraphQLGetFeed(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
query GetFeed(
	$flavour: Flavour!
	$persistent: BooleanFilter!
	$status: Status
	$visibility: Visibility
	$expired: BooleanFilter
	$filterParams: FilterParamsInput
	) {
	getFeed(
		flavour: $flavour
		persistent: $persistent
		status: $status
		visibility: $visibility
		expired: $expired
		filterParams: $filterParams
	) {
		uid
		flavour
		actions {
			id
			sequenceNumber
			name
			actionType
			handling
		}
		nudges {
			id
			sequenceNumber
			visibility
			status
			title
			text
			actions {
				id
				sequenceNumber
				name
				actionType
				handling
			}
			groups
			users
			links {
				id
				url
				linkType
			}
			notificationChannels
		}
		items {
			id
			sequenceNumber
			expiry
			persistent
			status
			visibility
			icon {
				id
				url
				linkType
			}
			author
			tagline
			label
			timestamp
			summary
			text
			links {
				id
				url
				linkType
			  }
			actions {
				id
				sequenceNumber
				name
				actionType
				handling
			}
			conversations {
				id
				sequenceNumber
				text
				replyTo
				postedByUID
				postedByName
			}
			users
			groups
			notificationChannels
		}
	}
}	  
	 `

	gql["variables"] = map[string]interface{}{
		"flavour":    "CONSUMER",
		"persistent": "BOTH",
		"status":     "PENDING",
		"visibility": "SHOW",
		"expired":    "FALSE",
		"filterParams": map[string]interface{}{
			"labels": []string{"a_label", "another_label"},
		},
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestRoutes(t *testing.T) {
	ctx := context.Background()
	router, err := graph.Router(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, router)

	uid := xid.New().String()
	fl := feed.FlavourConsumer
	itemID := ksuid.New().String()
	nudgeID := ksuid.New().String()
	actionID := ksuid.New().String()
	messageID := ksuid.New().String()

	type args struct {
		routeName string
		params    []string
	}
	tests := []struct {
		name    string
		args    args
		wantURL string
		wantErr bool
	}{
		{
			name: "get feed",
			args: args{
				routeName: "getFeed",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/", uid, fl.String()),
			wantErr: false,
		},
		{
			name: "get feed item",
			args: args{
				routeName: "getFeedItem",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/items/%s/", uid, fl.String(), itemID),
			wantErr: false,
		},
		{
			name: "get nudge",
			args: args{
				routeName: "getNudge",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/nudges/%s/", uid, fl.String(), nudgeID),
			wantErr: false,
		},
		{
			name: "get action",
			args: args{
				routeName: "getAction",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"actionID", actionID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/actions/%s/", uid, fl.String(), actionID),
			wantErr: false,
		},
		{
			name: "publish feed item",
			args: args{
				routeName: "publishFeedItem",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/items/", uid, fl.String()),
			wantErr: false,
		},
		{
			name: "publish nudge",
			args: args{
				routeName: "publishNudge",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/nudges/", uid, fl.String()),
			wantErr: false,
		},
		{
			name: "publish action",
			args: args{
				routeName: "publishAction",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/actions/", uid, fl.String()),
			wantErr: false,
		},
		{
			name: "post message",
			args: args{
				routeName: "postMessage",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%s/messages/", uid, fl.String(), itemID),
			wantErr: false,
		},
		{
			name: "post event",
			args: args{
				routeName: "postEvent",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/events/", uid, fl.String()),
			wantErr: false,
		},
		{
			name: "delete feed item",
			args: args{
				routeName: "deleteFeedItem",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/items/%s/", uid, fl.String(), itemID),
			wantErr: false,
		},
		{
			name: "delete nudge",
			args: args{
				routeName: "deleteNudge",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/nudges/%s/", uid, fl.String(), nudgeID),
			wantErr: false,
		},
		{
			name: "delete action",
			args: args{
				routeName: "deleteAction",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"actionID", actionID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/actions/%s/", uid, fl.String(), actionID),
			wantErr: false,
		},
		{
			name: "delete message",
			args: args{
				routeName: "deleteMessage",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"messageID", messageID,
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%s/messages/%s/", uid, fl.String(), itemID, messageID),
			wantErr: false,
		},
		{
			name: "resolve feed item",
			args: args{
				routeName: "resolveFeedItem",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/items/%s/resolve/", uid, fl.String(), itemID),
			wantErr: false,
		},
		{
			name: "unresolve feed item",
			args: args{
				routeName: "unresolveFeedItem",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/items/%s/unresolve/", uid, fl.String(), itemID),
			wantErr: false,
		},
		{
			name: "pin feed item",
			args: args{
				routeName: "pinFeedItem",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/items/%s/pin/", uid, fl.String(), itemID),
			wantErr: false,
		},
		{
			name: "unpin feed item",
			args: args{
				routeName: "unpinFeedItem",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/items/%s/unpin/", uid, fl.String(), itemID),
			wantErr: false,
		},
		{
			name: "hide feed item",
			args: args{
				routeName: "hideFeedItem",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/items/%s/hide/", uid, fl.String(), itemID),
			wantErr: false,
		},
		{
			name: "show feed item",
			args: args{
				routeName: "showFeedItem",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/items/%s/show/", uid, fl.String(), itemID),
			wantErr: false,
		},
		{
			name: "resolve nudge",
			args: args{
				routeName: "resolveNudge",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/nudges/%s/resolve/", uid, fl.String(), nudgeID),
			wantErr: false,
		},
		{
			name: "unresolve nudge",
			args: args{
				routeName: "unresolveNudge",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/nudges/%s/unresolve/", uid, fl.String(), nudgeID),
			wantErr: false,
		},
		{
			name: "show nudge",
			args: args{
				routeName: "showNudge",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/nudges/%s/show/", uid, fl.String(), nudgeID),
			wantErr: false,
		},
		{
			name: "hide nudge",
			args: args{
				routeName: "hideNudge",
				params: []string{
					"uid", uid,
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/nudges/%s/hide/", uid, fl.String(), nudgeID),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := router.Get(tt.args.routeName).URL(tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("route error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
			assert.NotZero(t, url)
			assert.Equal(t, tt.wantURL, url.String())
		})
	}
}

func TestGetFeed(t *testing.T) {
	uid := xid.New().String()
	consumer := feed.FlavourConsumer
	client := http.DefaultClient

	filterParams := feed.FilterParams{
		Labels: []string{"a", "label", "and", "another"},
	}
	filterParamsJSONBytes, err := json.Marshal(filterParams)
	assert.Nil(t, err)
	assert.NotNil(t, filterParamsJSONBytes)
	if err != nil {
		return
	}

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name                   string
		args                   args
		wantStatus             int
		wantNewFeedInitialized bool
		wantErr                bool
	}{
		{
			name: "successful fetch of a consumer feed",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/?persistent=BOTH",
					baseURL,
					uid,
					consumer,
				),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantNewFeedInitialized: true,
			wantStatus:             http.StatusOK,
			wantErr:                false,
		},
		{
			name: "fetch with a status filter",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/?persistent=BOTH&status=PENDING",
					baseURL,
					uid,
					consumer,
				),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fetch with a visibility filter",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/?persistent=BOTH&status=PENDING&visibility=SHOW",
					baseURL,
					uid,
					consumer,
				),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fetch with an expired filter",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/?persistent=BOTH&status=PENDING&visibility=SHOW&expired=FALSE",
					baseURL,
					uid,
					consumer,
				),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fetch with an expired filter",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/?persistent=BOTH&status=PENDING&visibility=SHOW&expired=FALSE&filterParams=%s",
					baseURL,
					uid,
					consumer,
					string(filterParamsJSONBytes),
				),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range getDefaultHeaders(t, baseURL) {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				// early exit
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantNewFeedInitialized {
				returnedFeed := &feed.Feed{}
				err = json.Unmarshal(data, returnedFeed)
				if err != nil {
					t.Errorf("can't unmarshal feed from response JSON: %w", err)
					return
				}

				if len(returnedFeed.Actions) < 1 {
					t.Error("the returned feed has no actions")
				}

				if len(returnedFeed.Nudges) < 1 {
					t.Errorf("the returned feed has no nudges")
				}

				if len(returnedFeed.Items) < 1 {
					t.Errorf("the returned feed has no items")
				}
			}
		})
	}
}

func TestGetFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid feed item retrieval",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGetNudge(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/",
					baseURL,
					uid,
					fl.String(),
					testNudge.ID,
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGetAction(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testAction := getTestAction()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testAction,
		baseURL,
		"publishAction",
	)
	if err != nil {
		t.Errorf("can't post action: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid action",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/actions/%s/",
					baseURL,
					uid,
					fl.String(),
					testAction.ID,
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent action",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/action/%s/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPublishFeedItem(t *testing.T) {
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	headers := getDefaultHeaders(t, baseURL)
	testItem := getTestItem()

	bs, err := json.Marshal(testItem)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid feed item publish",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/",
					baseURL,
					uid,
					fl.String(),
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/",
					baseURL,
					uid,
					fl.String(),
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestDeleteFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid delete",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent element delete - safe to repeat over and over",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestDeleteNudge(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid delete",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/",
					baseURL,
					uid,
					fl.String(),
					testNudge.ID,
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent element delete - safe to repeat over and over",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestDeleteAction(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testAction := getTestAction()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testAction,
		baseURL,
		"publishAction",
	)
	if err != nil {
		t.Errorf("can't post test action: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid delete",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/actions/%s/",
					baseURL,
					uid,
					fl.String(),
					testAction.ID,
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent element delete - safe to repeat over and over",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/actions/%s/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPostMessage(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}

	msg := getTestMessage()
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Errorf("can't marshal message to JSON: %s", err)
		return
	}
	payload := bytes.NewBuffer(msgBytes)

	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid message post",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%s/messages/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil message",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%s/messages/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}

	msg := getTestMessage()
	err = postMessage(
		ctx,
		t,
		uid,
		fl,
		&msg,
		baseURL,
		testItem.ID,
	)
	if err != nil {
		t.Errorf("can't post message: %s", err)
		return
	}

	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid delete",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%s/messages/%s/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
					msg.ID,
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent element delete - safe to repeat over and over",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%s/messages/%s/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
					ksuid.New().String(),
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestProcessEvent(t *testing.T) {
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	headers := getDefaultHeaders(t, baseURL)
	event := getTestEvent()

	bs, err := json.Marshal(event)
	if err != nil {
		t.Errorf("unable to marshal event to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid event publish",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/events/",
					baseURL,
					uid,
					fl.String(),
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil event",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/events/",
					baseURL,
					uid,
					fl.String(),
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPublishNudge(t *testing.T) {
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	headers := getDefaultHeaders(t, baseURL)
	nudge := testNudge()

	bs, err := json.Marshal(nudge)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid nudge publish",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/",
					baseURL,
					uid,
					fl.String(),
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/",
					baseURL,
					uid,
					fl.String(),
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestResolveNudge(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "resolve valid nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/resolve/",
					baseURL,
					uid,
					fl.String(),
					testNudge.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to resolve non existent nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/resolve/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestUnresolveNudge(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "resolve valid nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/unresolve/",
					baseURL,
					uid,
					fl.String(),
					testNudge.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to resolve non existent nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/unresolve/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestShowNudge(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "show valid nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/show/",
					baseURL,
					uid,
					fl.String(),
					testNudge.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to show non existent nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/show/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestHideNudge(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "hide valid nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/hide/",
					baseURL,
					uid,
					fl.String(),
					testNudge.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to hide non existent nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/nudges/%s/hide/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPublishAction(t *testing.T) {
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	headers := getDefaultHeaders(t, baseURL)
	action := getTestAction()

	bs, err := json.Marshal(action)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid action publish",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/actions/",
					baseURL,
					uid,
					fl.String(),
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil action",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/actions/",
					baseURL,
					uid,
					fl.String(),
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestResolveFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "resolve valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/resolve/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to resolve non existent feed uten",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/resolve/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestUnresolveFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "unresolve valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/unresolve/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to unresolve non existent feed uten",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/unresolve/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPinFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "pin valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/pin/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to pin non existent feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/pin/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestUnpinFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "unpin valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/unpin/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to unpin non existent feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/unpin/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestHideFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "hide valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/hide/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to hide non existent feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/hide/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestShowFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := feed.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "show valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/show/",
					baseURL,
					uid,
					fl.String(),
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to show non existent feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/items/%s/show/",
					baseURL,
					uid,
					fl.String(),
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func getInterserviceBearerTokenHeader(t *testing.T, rootDomain string) string {
	isc := getInterserviceClient(t, rootDomain)
	authToken, err := isc.CreateAuthToken()
	assert.Nil(t, err)
	assert.NotZero(t, authToken)
	bearerHeader := fmt.Sprintf("Bearer %s", authToken)
	return bearerHeader
}

func getDefaultHeaders(t *testing.T, rootDomain string) map[string]string {
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": getInterserviceBearerTokenHeader(t, rootDomain),
	}
}

func getGraphQLHeaders(t *testing.T) map[string]string {
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": getBearerTokenHeader(t),
	}
}

func getBearerTokenHeader(t *testing.T) string {
	ctx := context.Background()
	user, err := base.GetOrCreateFirebaseUser(ctx, base.TestUserEmail)
	if err != nil {
		t.Errorf("can't get or create firebase user: %s", err)
		return ""
	}

	if user == nil {
		t.Errorf("nil firebase user")
		return ""
	}

	customToken, err := base.CreateFirebaseCustomToken(ctx, user.UID)
	if err != nil {
		t.Errorf("can't create custom token: %s", err)
		return ""
	}

	if customToken == "" {
		t.Errorf("blank custom token: %s", err)
		return ""
	}

	idTokens, err := base.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		t.Errorf("can't authenticate custom token: %s", err)
		return ""
	}
	if idTokens == nil {
		t.Errorf("nil idTokens")
		return ""
	}

	return fmt.Sprintf("Bearer %s", idTokens.IDToken)
}

func getInterserviceClient(t *testing.T, rootDomain string) *base.InterServiceClient {
	service := base.ISCService{
		Name:       "feed",
		RootDomain: rootDomain,
	}
	isc, err := base.NewInterserviceClient(service)
	assert.Nil(t, err)
	assert.NotNil(t, isc)
	return isc
}

func randomPort() int {
	rand.Seed(time.Now().Unix())
	min := 32768
	max := 60999
	port := rand.Intn(max-min+1) + min
	return port
}

func startTestServer(ctx context.Context) (*http.Server, string, error) {
	// prepare the server
	port := randomPort()
	srv := graph.PrepareServer(ctx, port)
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if srv == nil {
		return nil, "", fmt.Errorf("nil test server")
	}

	// set up the TCP listener
	// this is done early so that we are sure we can connect to the port in
	// the tests; backlogs will be sent to the listener
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, "", fmt.Errorf("unable to listen on port %d: %w", port, err)
	}
	if l == nil {
		return nil, "", fmt.Errorf("nil test server listener")
	}
	log.Printf("LISTENING on port %d", port)

	// start serving
	go func() {
		err := srv.Serve(l)
		if err != nil {
			log.Printf("serve error: %s", err)
		}
	}()

	// the cleanup of this server (deferred shutdown) needs to occur in the
	// acceptance test that will use this
	return srv, baseURL, nil
}

func postElement(
	ctx context.Context,
	t *testing.T,
	uid string,
	fl feed.Flavour,
	el feed.Element,
	baseURL string,
	routeName string,
) error {
	router, err := graph.Router(ctx)
	if err != nil {
		t.Errorf("can't set up router: %s", err)
		return err
	}

	params := []string{
		"uid", uid,
		"flavour", fl.String(),
	}

	route := router.Get(routeName)
	if route == nil {
		return fmt.Errorf(
			"there's no registered route with the name `%s`", routeName)
	}
	path, err := router.Get(routeName).URL(params...)
	if err != nil {
		t.Errorf("can't get URL: %s", err)
		return err
	}
	url := fmt.Sprintf("%s%s", baseURL, path.String())

	data, err := json.Marshal(el)
	if err != nil {
		t.Errorf("can't marshal nudge to JSON: %s", err)
		return err
	}
	payload := bytes.NewBuffer(data)
	r, err := http.NewRequest(
		http.MethodPost,
		url,
		payload,
	)
	if err != nil {
		t.Errorf("error when creating request to post `%v` to %s: %s", payload, url, err)
		return err
	}
	if r == nil {
		t.Errorf("nil request when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	headers := getDefaultHeaders(t, baseURL)
	for k, v := range headers {
		r.Header.Add(k, v)
	}

	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}
	resp, err := client.Do(r)
	if resp == nil {
		t.Errorf("nil response: %s", err)
		return err
	}

	data, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	if err != nil {
		t.Errorf("error when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	assert.NotNil(t, data)
	if data == nil {
		t.Errorf("nil response when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("error status code `%d` and data `%s`", resp.StatusCode, data)
		return fmt.Errorf("error status code `%d` and data `%s`", resp.StatusCode, data)
	}

	return nil
}
func postMessage(
	ctx context.Context,
	t *testing.T,
	uid string,
	fl feed.Flavour,
	el feed.Element,
	baseURL string,
	itemID string,
) error {
	router, err := graph.Router(ctx)
	if err != nil {
		t.Errorf("can't set up router: %s", err)
		return err
	}

	params := []string{
		"uid", uid,
		"flavour", fl.String(),
		"itemID", itemID,
	}

	path, err := router.Get("postMessage").URL(params...)
	if err != nil {
		t.Errorf("can't get URL: %s", err)
		return err
	}
	url := fmt.Sprintf("%s%s", baseURL, path.String())

	data, err := json.Marshal(el)
	if err != nil {
		t.Errorf("can't marshal nudge to JSON: %s", err)
		return err
	}
	payload := bytes.NewBuffer(data)
	r, err := http.NewRequest(
		http.MethodPost,
		url,
		payload,
	)
	if err != nil {
		t.Errorf("error when creating request to post `%v` to %s: %s", payload, url, err)
		return err
	}
	if r == nil {
		t.Errorf("nil request when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	headers := getDefaultHeaders(t, baseURL)
	for k, v := range headers {
		r.Header.Add(k, v)
	}

	client := http.DefaultClient
	resp, err := client.Do(r)
	if resp == nil {
		t.Errorf("nil response: %s", err)
		return err
	}

	data, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	if err != nil {
		t.Errorf("error when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	assert.NotNil(t, data)
	if data == nil {
		t.Errorf("nil response when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("error status code `%d` and data `%s`", resp.StatusCode, data)
		return fmt.Errorf("error status code `%d` and data `%s`", resp.StatusCode, data)
	}

	return nil
}

func getTestItem() feed.Item {
	return feed.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Icon:           feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       feed.TextTypePlain,
		Links: []feed.Link{
			feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
			feed.GetYoutubeVideoLink(feed.SampleVideoURL, "title", "description", feed.LogoURL),
		},
		Actions: []feed.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
				ActionType:     feed.ActionTypeSecondary,
				Handling:       feed.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
				ActionType:     feed.ActionTypePrimary,
				Handling:       feed.HandlingInline,
			},
		},
		Conversations: []feed.Message{
			{
				ID:             "msg-2",
				Text:           "hii ni reply",
				ReplyTo:        "msg-1",
				PostedByName:   ksuid.New().String(),
				PostedByUID:    ksuid.New().String(),
				Timestamp:      time.Now(),
				SequenceNumber: int(time.Now().Unix()),
			},
		},
		Users: []string{
			"user-1",
			"user-2",
		},
		Groups: []string{
			"group-1",
			"group-2",
		},
		NotificationChannels: []feed.Channel{
			feed.ChannelFcm,
			feed.ChannelEmail,
			feed.ChannelSms,
			feed.ChannelWhatsapp,
		},
	}
}

func testNudge() *feed.Nudge {
	return &feed.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Expiry:         time.Now().Add(time.Hour * 24),
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []feed.Link{
			feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
		},
		Text: ksuid.New().String(),
		Actions: []feed.Action{
			getTestAction(),
		},
		Users: []string{
			ksuid.New().String(),
		},
		Groups: []string{
			ksuid.New().String(),
		},
		NotificationChannels: []feed.Channel{
			feed.ChannelEmail,
			feed.ChannelFcm,
			feed.ChannelSms,
			feed.ChannelWhatsapp,
		},
	}
}

func getTestSequenceNumber() int {
	return rand.Intn(intMax)
}

func getTestEvent() feed.Event {
	return feed.Event{
		ID:   ksuid.New().String(),
		Name: "TEST_EVENT",
		Context: feed.Context{
			UserID:         ksuid.New().String(),
			Flavour:        feed.FlavourConsumer,
			OrganizationID: ksuid.New().String(),
			LocationID:     ksuid.New().String(),
			Timestamp:      time.Now(),
		},
	}
}

func getTestAction() feed.Action {
	return feed.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Name:           "TEST_ACTION",
		Icon:           feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
		ActionType:     feed.ActionTypePrimary,
		Handling:       feed.HandlingFullPage,
	}
}

func getTestMessage() feed.Message {
	return feed.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func mapToJSONReader(m map[string]interface{}) (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
	}

	buf := bytes.NewBuffer(bs)
	return buf, nil
}

func TestGoogleCloudPubSubHandler(t *testing.T) {
	ctx := context.Background()
	b64 := base64.StdEncoding.EncodeToString([]byte(ksuid.New().String()))
	testPush := messaging.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: messaging.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      []byte(b64),
			Attributes: map[string]string{
				"topicID": feed.FeedRetrievalTopic,
			},
		},
	}
	testPushJSON, err := json.Marshal(testPush)
	if err != nil {
		t.Errorf("can't marshal JSON: %s", err)
		return
	}

	idTokenHTTPClient, err := idtoken.NewClient(
		ctx,
		messaging.DefaultPubsubTokenAudience,
	)
	if err != nil {
		t.Errorf("can't initialize idToken HTTP client: %s", err)
		return
	}

	pubsubURL := fmt.Sprintf("%s%s", baseURL, messaging.PubSubHandlerPath)
	req, err := http.NewRequest(
		http.MethodPost,
		pubsubURL,
		bytes.NewBuffer(testPushJSON),
	)
	if err != nil {
		t.Errorf("can't initialize request: %s", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	type args struct {
		r      *http.Request
		client *http.Client
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid pubsub format payload with valid auth",
			args: args{
				r:      req,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "no auth header",
			args: args{
				r:      req,
				client: http.DefaultClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.args.client.Do(tt.args.r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			respBs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unable to read response body: %s", err)
				return
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("wanted status %d, got %d and resp %s",
					tt.wantStatus, resp.StatusCode, string(respBs))
				return
			}

			if !tt.wantErr {
				decoded := map[string]string{}
				err = json.Unmarshal(respBs, &decoded)
				if err != nil {
					t.Errorf("can't decode response to map: %s", err)
					return
				}
				if decoded["status"] != "success" {
					t.Errorf("did not get success status")
					return
				}
			}
		})
	}
}
