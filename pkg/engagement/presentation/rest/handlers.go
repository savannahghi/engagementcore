package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"

	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/fcm"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/otp"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases"

	"net/http"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"

	errorcode "github.com/savannahghi/errorcodeutil"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/exceptions"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/helpers"
)

const (
	// StaticDir is the directory that contains schemata, default images etc
	StaticDir = "github.com/savannahghi/engagementcore:/static/"

	mbBytes              = 1048576
	serverTimeoutSeconds = 120
	bewellURL            = "https://bewell.co.ke/data_deletion.html/?id="
)

var errNotFound = fmt.Errorf("not found")

// PresentationHandlers represents all the REST API logic
type PresentationHandlers interface {
	GoogleCloudPubSubHandler(w http.ResponseWriter, r *http.Request)
	GetFeed() http.HandlerFunc

	GetFeedItem() http.HandlerFunc

	GetNudge() http.HandlerFunc

	GetAction() http.HandlerFunc

	PublishFeedItem() http.HandlerFunc

	DeleteFeedItem() http.HandlerFunc

	ResolveFeedItem() http.HandlerFunc

	PinFeedItem() http.HandlerFunc

	UnpinFeedItem() http.HandlerFunc

	HideFeedItem() http.HandlerFunc

	ShowFeedItem() http.HandlerFunc

	UnresolveFeedItem() http.HandlerFunc

	PublishNudge() http.HandlerFunc

	ResolveNudge() http.HandlerFunc

	ResolveDefaultNudge() http.HandlerFunc

	UnresolveNudge() http.HandlerFunc

	HideNudge() http.HandlerFunc

	ShowNudge() http.HandlerFunc

	DeleteNudge() http.HandlerFunc

	PublishAction() http.HandlerFunc

	DeleteAction() http.HandlerFunc

	PostMessage() http.HandlerFunc

	DeleteMessage() http.HandlerFunc

	ProcessEvent() http.HandlerFunc

	Upload() http.HandlerFunc

	FindUpload() http.HandlerFunc

	SendEmail() http.HandlerFunc

	SendToMany() http.HandlerFunc

	GetNotificationHandler() http.HandlerFunc

	GetIncomingMessageHandler() http.HandlerFunc

	GetFallbackHandler() http.HandlerFunc

	PhoneNumberVerificationCodeHandler() http.HandlerFunc

	SendOTPHandler() http.HandlerFunc

	SendRetryOTPHandler() http.HandlerFunc

	VerifyRetryOTPHandler() http.HandlerFunc

	VerifyRetryEmailOTPHandler() http.HandlerFunc

	SendNotificationHandler() http.HandlerFunc

	UpdateMailgunDeliveryStatus() http.HandlerFunc

	DataDeletionRequestCallback() http.HandlerFunc

	GetTwilioVideoCallbackFunc() http.HandlerFunc

	SendTemporaryPIN() http.HandlerFunc

	SendEmailOTP() http.HandlerFunc
}

// PresentationHandlersImpl represents the usecase implementation object
type PresentationHandlersImpl struct {
	usecases       usecases.Interactor
	infrastructure infrastructure.Interactor
}

// NewPresentationHandlers initializes a new rest handlers usecase
func NewPresentationHandlers(infrastructure infrastructure.Interactor, usecases usecases.Interactor) PresentationHandlers {
	return &PresentationHandlersImpl{infrastructure: infrastructure, usecases: usecases}
}

// GoogleCloudPubSubHandler receives push messages from Google Cloud Pub-Sub
func (p PresentationHandlersImpl) GoogleCloudPubSubHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	m, err := pubsubtools.VerifyPubSubJWTAndDecodePayload(w, r)
	if err != nil {
		serverutils.WriteJSONResponse(
			w,
			errorcode.ErrorMap(err),
			http.StatusBadRequest,
		)
		return
	}

	topicID, err := pubsubtools.GetPubSubTopic(m)
	if err != nil {
		serverutils.WriteJSONResponse(
			w,
			errorcode.ErrorMap(err),
			http.StatusBadRequest,
		)
		return
	}

	// get the UID frrom the payload
	var envelope dto.NotificationEnvelope
	err = json.Unmarshal(m.Message.Data, &envelope)
	if err != nil {
		serverutils.WriteJSONResponse(
			w,
			errorcode.ErrorMap(err),
			http.StatusBadRequest,
		)
		return
	}
	ctx = addUIDToContext(ctx, envelope.UID)

	switch topicID {
	case helpers.AddPubSubNamespace(common.ItemPublishTopic):
		err = p.usecases.HandleItemPublish(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemDeleteTopic):
		err = p.usecases.HandleItemDelete(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemResolveTopic):
		err = p.usecases.HandleItemResolve(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemUnresolveTopic):
		err = p.usecases.HandleItemUnresolve(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemHideTopic):
		err = p.usecases.HandleItemHide(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemShowTopic):
		err = p.usecases.HandleItemShow(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemPinTopic):
		err = p.usecases.HandleItemPin(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ItemUnpinTopic):
		err = p.usecases.HandleItemUnpin(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgePublishTopic):
		err = p.usecases.HandleNudgePublish(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeDeleteTopic):
		err = p.usecases.HandleNudgeDelete(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeResolveTopic):
		err = p.usecases.HandleNudgeResolve(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeUnresolveTopic):
		err = p.usecases.HandleNudgeUnresolve(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeHideTopic):
		err = p.usecases.HandleNudgeHide(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.NudgeShowTopic):
		err = p.usecases.HandleNudgeShow(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ActionPublishTopic):
		err = p.usecases.HandleActionPublish(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.ActionDeleteTopic):
		err = p.usecases.HandleActionDelete(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.MessagePostTopic):
		err = p.usecases.HandleMessagePost(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.MessageDeleteTopic):
		err = p.usecases.HandleMessageDelete(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.IncomingEventTopic):
		err = p.usecases.HandleIncomingEvent(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.FcmPublishTopic):
		err = p.usecases.HandleSendNotification(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	case helpers.AddPubSubNamespace(common.SentEmailTopic):
		err = p.usecases.SendNotificationEmail(ctx, m)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
	default:
		// the topic should be anticipated/handled here
		errMsg := fmt.Sprintf(
			"pub sub handler error: unknown topic `%s`",
			topicID,
		)
		log.Print(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	resp := map[string]string{"status": "success"}
	marshalledSuccessMsg, err := json.Marshal(resp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	_, _ = w.Write(marshalledSuccessMsg)
}

// GetFeed retrieves and serves a feed
func (p PresentationHandlersImpl) GetFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		uid, flavour, anonymous, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		persistent, err := getRequiredBooleanFilterQueryParam(r, "persistent")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		status, err := getOptionalStatusQueryParam(r, "status")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		visibility, err := getOptionalVisibilityQueryParam(r, "visibility")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		expired, err := getOptionalBooleanFilterQueryParam(r, "expired")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		playMP4, err := getplayMP4QueryParam(r, "playMP4")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		filterParams, err := getOptionalFilterParamsQueryParam(
			r,
			"filterParams",
		)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		feed, err := p.usecases.GetFeed(
			addUIDToContext(ctx, *uid),
			uid,
			anonymous,
			*flavour,
			playMP4,
			persistent,
			status,
			visibility,
			expired,
			filterParams,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}
		bs, err := json.Marshal(feed)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, bs)
	}
}

// GetFeedItem retrieves a single feed item
func (p PresentationHandlersImpl) GetFeedItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		itemID, err := getStringVar(r, "itemID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		item, err := p.usecases.GetFeedItem(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			itemID,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		if item == nil {
			respondWithError(w, http.StatusNotFound, errNotFound)
		}

		marshalled, err := item.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// GetNudge retrieves a single nudge
func (p PresentationHandlersImpl) GetNudge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		nudgeID, err := getStringVar(r, "nudgeID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		ctx = addUIDToContext(ctx, *uid)
		nudge, err := p.usecases.GetNudge(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			nudgeID,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		if nudge == nil {
			respondWithError(w, http.StatusNotFound, errNotFound)
		}

		marshalled, err := nudge.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// GetAction retrieves a single action
func (p PresentationHandlersImpl) GetAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		actionID, err := getStringVar(r, "actionID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		action, err := p.usecases.GetAction(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			actionID,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		if action == nil {
			respondWithError(w, http.StatusNotFound, errNotFound)
		}

		marshalled, err := action.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

func readBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, mbBytes))
	if err != nil {
		return nil, fmt.Errorf("can't read request body: %w", err)
	}
	return body, nil
}

// PublishFeedItem posts a feed item
func (p PresentationHandlersImpl) PublishFeedItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		item := &feedlib.Item{}
		err = item.ValidateAndUnmarshal(data)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		publishedItem, err := p.usecases.PublishFeedItem(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			item,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := publishedItem.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// DeleteFeedItem removes a feed item
func (p PresentationHandlersImpl) DeleteFeedItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		itemID, err := getStringVar(r, "itemID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = p.usecases.DeleteFeedItem(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			itemID,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// ResolveFeedItem marks a feed item as done
func (p PresentationHandlersImpl) ResolveFeedItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		patchItem(ctx, p.usecases.ResolveFeedItem, w, r)
	}
}

// PinFeedItem marks a feed item as done
func (p PresentationHandlersImpl) PinFeedItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		patchItem(ctx, p.usecases.PinFeedItem, w, r)
	}
}

// UnpinFeedItem marks a feed item as done
func (p PresentationHandlersImpl) UnpinFeedItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		patchItem(ctx, p.usecases.UnpinFeedItem, w, r)
	}
}

// HideFeedItem marks a feed item as done
func (p PresentationHandlersImpl) HideFeedItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		patchItem(ctx, p.usecases.HideFeedItem, w, r)
	}
}

// ShowFeedItem marks a feed item as done
func (p PresentationHandlersImpl) ShowFeedItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		patchItem(ctx, p.usecases.ShowFeedItem, w, r)
	}
}

// UnresolveFeedItem marks a feed item as not resolved
func (p PresentationHandlersImpl) UnresolveFeedItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		patchItem(ctx, p.usecases.UnresolveFeedItem, w, r)
	}
}

// PublishNudge posts a new nudge
func (p PresentationHandlersImpl) PublishNudge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		nudge := &feedlib.Nudge{}
		err = nudge.ValidateAndUnmarshal(data)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		publishedNudge, err := p.usecases.PublishNudge(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			nudge,
		)
		if err != nil {
			if strings.Contains(
				err.Error(),
				"found an existing nudge with same title",
			) {
				respondWithError(w, http.StatusConflict, err)
				return
			}
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := publishedNudge.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// ResolveNudge marks a nudge as resolved
func (p PresentationHandlersImpl) ResolveNudge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		patchNudge(ctx, p.usecases.ResolveNudge, w, r)
	}
}

// ResolveDefaultNudge marks a default nudges as resolved
func (p PresentationHandlersImpl) ResolveDefaultNudge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		title, err := getStringVar(r, "title")

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		nudge, err := p.usecases.GetDefaultNudgeByTitle(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			title,
		)
		if err != nil {
			if errors.Is(err, exceptions.ErrNilNudge) {
				respondWithError(w, http.StatusNotFound, err)
				return
			}
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		if nudge.Status == feedlib.StatusDone {
			respondWithJSON(w, http.StatusOK, marshalled)
		}

		_, err = p.usecases.ResolveNudge(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			nudge.ID,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// UnresolveNudge marks a nudge as not resolved
func (p PresentationHandlersImpl) UnresolveNudge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		patchNudge(ctx, p.usecases.UnresolveNudge, w, r)
	}
}

// HideNudge marks a nudge as not resolved
func (p PresentationHandlersImpl) HideNudge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		patchNudge(ctx, p.usecases.HideNudge, w, r)
	}
}

// ShowNudge marks a nudge as not resolved
func (p PresentationHandlersImpl) ShowNudge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		patchNudge(ctx, p.usecases.ShowNudge, w, r)
	}
}

// DeleteNudge permanently deletes a nudge
func (p PresentationHandlersImpl) DeleteNudge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		nudgeID, err := getStringVar(r, "nudgeID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = p.usecases.DeleteNudge(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			nudgeID,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// PublishAction posts a new action to a user's feed
func (p PresentationHandlersImpl) PublishAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		action := &feedlib.Action{}
		err = action.ValidateAndUnmarshal(data)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		publishedAction, err := p.usecases.PublishAction(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			action,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := publishedAction.ValidateAndMarshal()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// DeleteAction permanently removes an action from a user's feed
func (p PresentationHandlersImpl) DeleteAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		actionID, err := getStringVar(r, "actionID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = p.usecases.DeleteAction(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			actionID,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// PostMessage adds a message to a thread
func (p PresentationHandlersImpl) PostMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		itemID, err := getStringVar(r, "itemID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		message := &feedlib.Message{}
		err = message.ValidateAndUnmarshal(data)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		postedMessage, err := p.usecases.PostMessage(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			itemID,
			message,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := json.Marshal(postedMessage)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// DeleteMessage removes a message from a thread
func (p PresentationHandlersImpl) DeleteMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		itemID, err := getStringVar(r, "itemID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		messageID, err := getStringVar(r, "messageID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = p.usecases.DeleteMessage(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			itemID,
			messageID,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// ProcessEvent saves an event
func (p PresentationHandlersImpl) ProcessEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		event := &feedlib.Event{}
		err = event.ValidateAndUnmarshal(data)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		err = p.usecases.ProcessEvent(
			addUIDToContext(ctx, *uid),
			*uid,
			*flavour,
			event,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		resp := map[string]string{"status": "success"}
		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// Upload saves an upload in cloud storage
func (p PresentationHandlersImpl) Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data, err := readBody(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		uploadInput := profileutils.UploadInput{}
		err = json.Unmarshal(data, &uploadInput)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		if uploadInput.Base64data == "" {
			err := fmt.Errorf("blank upload base64 data")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		if uploadInput.Filename == "" {
			err := fmt.Errorf("blank upload filename")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		if uploadInput.Title == "" {
			err := fmt.Errorf("blank upload title")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		upload, err := p.infrastructure.Upload(ctx, uploadInput)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		if upload == nil {
			err := fmt.Errorf("nil upload in response from upload service")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		marshalled, err := json.Marshal(upload)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// FindUpload retrieves an upload by it's ID
func (p PresentationHandlersImpl) FindUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		uploadID, err := getStringVar(r, "uploadID")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		upload, err := p.infrastructure.FindUploadByID(ctx, uploadID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		if upload == nil {
			err := fmt.Errorf("nil upload in response from upload service")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		marshalled, err := json.Marshal(upload)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// SendEmail sends the specified email to the recipient(s) specified in `to`
// and returns the status
func (p PresentationHandlersImpl) SendEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.EMailMessage{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.Subject == "" {
			err := fmt.Errorf("blank email subject")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		if payload.Text == "" {
			err := fmt.Errorf("blank email text")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		if len(payload.To) == 0 {
			err := fmt.Errorf("no destination email supplied")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		resp, _, err := p.infrastructure.SendEmail(
			ctx,
			payload.Subject,
			payload.Subject,
			&payload.Text,
			payload.To...,
		)
		if err != nil {
			err := fmt.Errorf("email not sent: %s", err)
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// SendToMany sends a data message to the specified recipient
func (p PresentationHandlersImpl) SendToMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.SendSMSPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		for _, phoneNo := range payload.To {
			_, err := converterandformatter.NormalizeMSISDN(phoneNo)
			if err != nil {
				err := fmt.Errorf(
					"can't send sms, expected a valid phone number",
				)
				respondWithError(w, http.StatusBadRequest, err)
				return
			}
		}

		if payload.Message == "" {
			err := fmt.Errorf("can't send sms, expected a message")
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		resp, err := p.infrastructure.SendToMany(
			ctx,
			payload.To,
			payload.Message,
		)
		if err != nil {
			err := fmt.Errorf("sms not sent: %s", err)

			isBadReq := strings.Contains(err.Error(), "http error status: 400")

			if isBadReq {
				respondWithError(w, http.StatusBadRequest, err)
				return
			}

			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(w, http.StatusOK, marshalled)
	}
}

// GetNotificationHandler returns a handler that processes an Africa's Talking payment notification
func (p PresentationHandlersImpl) GetNotificationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.Message{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.AccountSID == "" {
			err := fmt.Errorf(
				"twilio notification payload not parsed correctly",
			)
			log.Printf("Twilio callback error: %s", err)
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusInternalServerError,
			)
			return
		}

		// save Twilio response for audit purposes
		err := p.infrastructure.SaveTwilioCallbackResponse(ctx, *payload)
		if err != nil {
			err := fmt.Errorf("twilio notification payload not saved")
			log.Printf("Twilio callback error: %s", err)
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusInternalServerError,
			)
			return
		}
		// TODO Common pathway for saving, returning OK etc

		type okResp struct {
			Status string `json:"status"`
		}
		serverutils.WriteJSONResponse(w, okResp{Status: "ok"}, http.StatusOK)
	}
}

// GetIncomingMessageHandler returns a handler that processes an Africa's Talking payment notification
func (p PresentationHandlersImpl) GetIncomingMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.Message{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.AccountSID == "" {
			err := fmt.Errorf(
				"twilio notification payload not parsed correctly",
			)
			log.Printf("Twilio callback error: %s", err)
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusInternalServerError,
			)
			return
		}

		// save Twilio response for audit purposes
		err := p.infrastructure.SaveTwilioCallbackResponse(ctx, *payload)
		if err != nil {
			err := fmt.Errorf("twilio notification payload not saved")
			log.Printf("Twilio callback error: %s", err)
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(err),
				http.StatusInternalServerError,
			)
			return
		}
		// TODO Common pathway for saving, returning OK etc

		type okResp struct {
			Status string `json:"status"`
		}
		serverutils.WriteJSONResponse(w, okResp{Status: "ok"}, http.StatusOK)
	}
}

// GetFallbackHandler returns a handler that processes an Africa's Talking payment notification
func (p PresentationHandlersImpl) GetFallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO ErrorCode and ErrorURL sent here as params
		// TODO Implement WhatsAPP fallback handler: serverutils.DecodeJSONToTargetStruct(w, r, notificationPayload)
		// errorcode.ErrorMap(fmt.Errorf("unbound mandatory notification payload fields")),
		// serverutils.WriteJSONResponse(w, okResp{Status: "ok"}, http.StatusOK)
	}
}

// PhoneNumberVerificationCodeHandler process ISC request to PhoneNumberVerificationCode
func (p PresentationHandlersImpl) PhoneNumberVerificationCodeHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type PayloadRequest struct {
			To               string `json:"to"`
			Code             string `json:"code"`
			MarketingMessage string `json:"marketingMessage"`
		}

		payloadRequest := &PayloadRequest{}

		serverutils.DecodeJSONToTargetStruct(rw, r, payloadRequest)

		ok, err := p.infrastructure.PhoneNumberVerificationCode(
			ctx,
			payloadRequest.To,
			payloadRequest.Code,
			payloadRequest.MarketingMessage,
		)
		if err != nil {
			errorcode.RespondWithError(rw, http.StatusInternalServerError, err)
			return
		}

		type PayloadResponse struct {
			Status bool `json:"status"`
		}

		response := &PayloadResponse{Status: ok}
		serverutils.WriteJSONResponse(rw, response, http.StatusOK)
	}
}

// SendOTPHandler is an isc api that generates and sends an otp to an msisdn
func (p PresentationHandlersImpl) SendOTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload, err := otp.ValidateSendOTPPayload(w, r)
		if err != nil {
			errorcode.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		code, err := p.infrastructure.GenerateAndSendOTP(
			ctx,
			payload.Msisdn,
			payload.AppID,
		)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(
					fmt.Errorf("unable to generate and send otp: %v", err),
				),
				http.StatusInternalServerError,
			)
			return
		}

		serverutils.WriteJSONResponse(w, code, http.StatusOK)
	}
}

// SendEmailOTP generates and send an OTP to the email supplied
func (p PresentationHandlersImpl) SendEmailOTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload, err := otp.ValidateSendEmailOTPPayload(w, r)
		if err != nil {
			errorcode.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		code, err := p.infrastructure.EmailVerificationOtp(
			ctx,
			&payload.Email,
		)
		if err != nil {
			serverutils.WriteJSONResponse(
				w,
				errorcode.ErrorMap(
					fmt.Errorf("unable to generate and send otp: %v", err),
				),
				http.StatusInternalServerError,
			)
			return
		}

		serverutils.WriteJSONResponse(w, code, http.StatusOK)
	}
}

// SendRetryOTPHandler is an isc api that generates
// fallback OTPs when Africa is talking sms fails
func (p PresentationHandlersImpl) SendRetryOTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload, err := otp.ValidateGenerateRetryOTPPayload(w, r)
		if err != nil {
			errorcode.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		code, err := p.infrastructure.GenerateRetryOTP(
			ctx,
			payload.Msisdn,
			payload.RetryStep,
			payload.AppID,
		)
		if err != nil {
			err := errorcode.ErrorMap(
				fmt.Errorf(
					"unable to generate and send a fallback OTP: %v",
					err,
				),
			)
			serverutils.WriteJSONResponse(
				w,
				err,
				http.StatusInternalServerError,
			)
			return
		}

		serverutils.WriteJSONResponse(w, code, http.StatusOK)
	}
}

// VerifyRetryOTPHandler is an isc api that confirms OTPs earlier sent
func (p PresentationHandlersImpl) VerifyRetryOTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload, err := otp.ValidateVerifyOTPPayload(w, r, false)
		if err != nil {
			errorcode.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		isVerified, err := p.infrastructure.VerifyOtp(
			ctx,
			payload.Msisdn,
			payload.VerificationCode,
		)
		if err != nil {
			errorcode.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		type otpResponse struct {
			IsVerified bool `json:"IsVerified"`
		}

		serverutils.WriteJSONResponse(
			w,
			otpResponse{IsVerified: isVerified},
			http.StatusOK,
		)
	}
}

// VerifyRetryEmailOTPHandler is an isc api that confirms OTPs earlier sent via email.
func (p PresentationHandlersImpl) VerifyRetryEmailOTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload, err := otp.ValidateVerifyOTPPayload(w, r, true)
		if err != nil {
			errorcode.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		isVerified, err := p.infrastructure.VerifyEmailOtp(
			ctx,
			payload.Email,
			payload.VerificationCode,
		)
		if err != nil {
			errorcode.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		type otpResponse struct {
			IsVerified bool `json:"IsVerified"`
		}

		serverutils.WriteJSONResponse(
			w,
			otpResponse{IsVerified: isVerified},
			http.StatusOK,
		)
	}
}

// SendNotificationHandler sends a data message to the specified registration tokens.
func (p PresentationHandlersImpl) SendNotificationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload, payloadErr := fcm.ValidateSendNotificationPayload(w, r)
		if payloadErr != nil {
			errorcode.ReportErr(w, payloadErr, http.StatusBadRequest)
			return
		}

		_, err := p.infrastructure.SendNotification(
			ctx,
			payload.RegistrationTokens,
			payload.Data,
			payload.Notification,
			payload.Android,
			payload.Ios,
			payload.Web,
		)
		if err != nil {
			err := fmt.Errorf("notification not sent: %s", err)

			isBadReq := strings.Contains(err.Error(), "http error status: 400")

			if isBadReq {
				errorcode.ReportErr(w, err, http.StatusBadRequest)
				return
			}

			errorcode.ReportErr(w, err, http.StatusInternalServerError)
			return
		}

		type okResp struct {
			Status string `json:"status"`
		}
		serverutils.WriteJSONResponse(w, okResp{Status: "ok"}, http.StatusOK)
	}
}

// UpdateMailgunDeliveryStatus gets the status of the sent emails and logs them in the database
func (p PresentationHandlersImpl) UpdateMailgunDeliveryStatus() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.MailgunEvent{}
		serverutils.DecodeJSONToTargetStruct(rw, r, payload)

		emailLog, err := p.infrastructure.Repository.UpdateMailgunDeliveryStatus(
			ctx,
			payload,
		)
		if err != nil {
			err := fmt.Errorf("email not sent: %s", err)
			respondWithError(rw, http.StatusInternalServerError, err)
			return
		}

		marshalled, err := json.Marshal(emailLog)
		if err != nil {
			respondWithError(rw, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(rw, http.StatusOK, marshalled)
	}
}

// DataDeletionRequestCallback is a Facebook's data deletion request callback
func (p PresentationHandlersImpl) DataDeletionRequestCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.FormValue("user_id")
		data := map[string]interface{}{
			"url":               bewellURL + userID,
			"confirmation_code": userID,
		}
		resp, err := json.Marshal(data)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(w, http.StatusOK, resp)
	}
}

// GetTwilioVideoCallbackFunc generates a Twilio Video callback handling function.
//
// Twilio sends the data with the "Content-Type" header to “application/x-www-urlencoded”.
func (p PresentationHandlersImpl) GetTwilioVideoCallbackFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Form == nil || len(r.Form) == 0 {
			return
		}
		if err := p.infrastructure.Repository.SaveTwilioVideoCallbackStatus(
			r.Context(),
			dto.CallbackData{Values: r.Form},
		); err != nil {
			log.Printf(
				"Twilio video callback error: unable to save callback response: %v",
				err,
			)
		}
	}
}

// SendTemporaryPIN is a twilio endpoint to send PIN via whatsapp and SMS
func (p PresentationHandlersImpl) SendTemporaryPIN() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		payload := dto.TemporaryPIN{}
		serverutils.DecodeJSONToTargetStruct(rw, r, payload)

		err := p.infrastructure.SendTemporaryPIN(ctx, payload)
		if err != nil {
			respondWithError(rw, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(rw, http.StatusOK, nil)
	}
}
