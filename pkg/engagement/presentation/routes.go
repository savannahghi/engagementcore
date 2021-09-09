package presentation

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/twilio"
	"github.com/savannahghi/engagementcore/pkg/engagement/presentation/rest"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func newPresentationHandlers() rest.PresentationHandlers {
	infrastructure := infrastructure.NewInteractor()
	usecases := usecases.NewUsecasesInteractor(infrastructure)
	return rest.NewPresentationHandlers(infrastructure, usecases)
}

// SharedUnauthenticatedRoutes return REST routes shared by open/closed engagement services
func SharedUnauthenticatedRoutes(ctx context.Context, r *mux.Router) {
	h := newPresentationHandlers()

	r.Use(otelmux.Middleware(serverutils.MetricsCollectorService("engagement")))
	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error
	r.Use(serverutils.RequestDebugMiddleware())

	// Add Middleware that records the metrics for our HTTP routes
	r.Use(serverutils.CustomHTTPRequestMetricsMiddleware())

	// Unauthenticated routes
	r.Path("/ide").HandlerFunc(playground.Handler("GraphQL IDE", "/graphql"))
	r.Path("/health").HandlerFunc(HealthStatusCheck)

	r.Path(pubsubtools.PubSubHandlerPath).Methods(
		http.MethodPost).HandlerFunc(h.GoogleCloudPubSubHandler)

	// Expose a bulk SMS sending endpoint
	r.Path("/send_sms").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(h.SendToMany())

	// Callbacks
	r.Path("/twilio_notification").
		Methods(http.MethodPost).
		HandlerFunc(h.GetNotificationHandler())
	r.Path("/twilio_incoming_message").
		Methods(http.MethodPost).
		HandlerFunc(h.GetIncomingMessageHandler())
	r.Path("/twilio_fallback").
		Methods(http.MethodPost).
		HandlerFunc(h.GetFallbackHandler())
	r.Path(twilio.TwilioCallbackPath).
		Methods(http.MethodPost).
		HandlerFunc(h.GetTwilioVideoCallbackFunc())

	r.Path("/upload").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(h.Upload())

}

// SharedAuthenticatedISCRoutes return ISC REST routes shared by open/closed engagement services
func SharedAuthenticatedISCRoutes(ctx context.Context, r *mux.Router) {
	h := newPresentationHandlers()

	// Interservice Authenticated routes
	feedISC := r.PathPrefix("/feed/{uid}/{flavour}/{isAnonymous}/").Subrouter()
	feedISC.Use(interserviceclient.InterServiceAuthenticationMiddleware())

	// retrieval
	feedISC.Methods(
		http.MethodGet,
	).Path("/").HandlerFunc(
		h.GetFeed(),
	).Name("getFeed")

	feedISC.Methods(
		http.MethodGet,
	).Path("/items/{itemID}/").HandlerFunc(
		h.GetFeedItem(),
	).Name("getFeedItem")

	feedISC.Methods(
		http.MethodGet,
	).Path("/nudges/{nudgeID}/").HandlerFunc(
		h.GetNudge(),
	).Name("getNudge")

	feedISC.Methods(
		http.MethodGet,
	).Path("/actions/{actionID}/").HandlerFunc(
		h.GetAction(),
	).Name("getAction")

	// creation
	feedISC.Methods(
		http.MethodPost,
	).Path("/items/").HandlerFunc(
		h.PublishFeedItem(),
	).Name("publishFeedItem")

	feedISC.Methods(
		http.MethodPost,
	).Path("/nudges/").HandlerFunc(
		h.PublishNudge(),
	).Name("publishNudge")

	feedISC.Methods(
		http.MethodPost,
	).Path("/actions/").HandlerFunc(
		h.PublishAction(),
	).Name("publishAction")

	feedISC.Methods(
		http.MethodPost,
	).Path("/{itemID}/messages/").HandlerFunc(
		h.PostMessage(),
	).Name("postMessage")

	feedISC.Methods(
		http.MethodPost,
	).Path("/events/").HandlerFunc(
		h.ProcessEvent(),
	).Name("postEvent")

	// deleting
	feedISC.Methods(
		http.MethodDelete,
	).Path("/items/{itemID}/").HandlerFunc(
		h.DeleteFeedItem(),
	).Name("deleteFeedItem")

	feedISC.Methods(
		http.MethodDelete,
	).Path("/nudges/{nudgeID}/").HandlerFunc(
		h.DeleteNudge(),
	).Name("deleteNudge")

	feedISC.Methods(
		http.MethodDelete,
	).Path("/actions/{actionID}/").HandlerFunc(
		h.DeleteAction(),
	).Name("deleteAction")

	feedISC.Methods(
		http.MethodDelete,
	).Path("/{itemID}/messages/{messageID}/").HandlerFunc(
		h.DeleteMessage(),
	).Name("deleteMessage")

	// modifying (patching)
	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/resolve/").HandlerFunc(
		h.ResolveFeedItem(),
	).Name("resolveFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/unresolve/").HandlerFunc(
		h.UnresolveFeedItem(),
	).Name("unresolveFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/pin/").HandlerFunc(
		h.PinFeedItem(),
	).Name("pinFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/unpin/").HandlerFunc(
		h.UnpinFeedItem(),
	).Name("unpinFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/hide/").HandlerFunc(
		h.HideFeedItem(),
	).Name("hideFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/items/{itemID}/show/").HandlerFunc(
		h.ShowFeedItem(),
	).Name("showFeedItem")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/resolve/").HandlerFunc(
		h.ResolveNudge(),
	).Name("resolveNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/defaultnudges/{title}/resolve/").HandlerFunc(
		h.ResolveDefaultNudge(),
	).Name("resolveDefaultNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/unresolve/").HandlerFunc(
		h.UnresolveNudge(),
	).Name("unresolveNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/show/").HandlerFunc(
		h.ShowNudge(),
	).Name("showNudge")

	feedISC.Methods(
		http.MethodPatch,
	).Path("/nudges/{nudgeID}/hide/").HandlerFunc(
		h.HideNudge(),
	).Name("hideNudge")

	isc := r.PathPrefix("/internal/").Subrouter()
	isc.Use(interserviceclient.InterServiceAuthenticationMiddleware())

	isc.Methods(
		http.MethodGet,
	).Path("/upload/{uploadID}/").HandlerFunc(
		h.FindUpload(),
	).Name("getUpload")

	isc.Methods(
		http.MethodPost,
	).Path("/upload/").HandlerFunc(
		h.Upload(),
	).Name("upload")

	isc.Methods(
		http.MethodPost,
	).Path("/send_email").HandlerFunc(
		h.SendEmail(),
	).Name("sendEmail")

	isc.Methods(
		http.MethodPost,
	).Path("/mailgun_delivery_webhook").HandlerFunc(
		h.UpdateMailgunDeliveryStatus(),
	).Name("mailgun_delivery_webhook")

	isc.Methods(
		http.MethodPost,
	).Path("/send_sms").HandlerFunc(
		h.SendToMany(),
	).Name("sendToMany")

	isc.Path("/verify_phonenumber").Methods(http.MethodPost).HandlerFunc(
		h.PhoneNumberVerificationCodeHandler(),
	)

	isc.Path("/send_otp/").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(h.SendOTPHandler())

	isc.Path("/send_retry_otp/").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(h.SendRetryOTPHandler())

	isc.Path("/verify_otp/").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(h.VerifyRetryOTPHandler())

	isc.Path("/verify_email_otp/").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(h.VerifyRetryEmailOTPHandler())

	isc.Path("/send_notification").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(h.SendNotificationHandler())

	isc.Path("/send_temporary_pin").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(h.SendTemporaryPIN())
}

// AuthenticatedGraphQLRoute inits an authenticated GraphQL route
func AuthenticatedGraphQLRoute(ctx context.Context, r *mux.Router) {
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		log.Fatal(err)
	}

	infrastructure := infrastructure.NewInteractor()
	usecases := usecases.NewUsecasesInteractor(infrastructure)

	authR := r.Path("/graphql").Subrouter()
	authR.Use(firebasetools.AuthenticationMiddleware(firebaseApp))
	authR.Methods(
		http.MethodPost,
		http.MethodGet,
	).HandlerFunc(GQLHandler(ctx, usecases, infrastructure))
}
