package feed

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/markbates/pkger"
	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/graph/library"
)

const (
	// DefaultLabel is the label used for welcome content
	DefaultLabel = "WELCOME"

	// DefaultIconPath is the path to the default Be.Well logo
	DefaultIconPath = StaticBase + "/bewell_logo.png"

	// StaticBase is the default path at which static assets are hosted
	StaticBase = "https://assets.healthcloud.co.ke"

	defaultSequenceNumber = 1
	defaultPostedByUID    = "hOcaUv8dqqgmWYf9HEhjdudgf0b2"
	futureHours           = 878400 // hours in a century of leap years...

	getConsultationActionName     = "GET_CONSULTATION"
	getMedicineActionName         = "GET_MEDICINE"
	getTestActionName             = "GET_TEST"
	getInsuranceActionName        = "GET_INSURANCE"
	getCoachingActionName         = "GET_COACHING"
	addPatientActionName          = "ADD_PATIENT"
	searchPatientActionName       = "SEARCH_PATIENT"
	addInsuranceActionName        = "ADD_INSURANCE"
	addNHIFActionName             = "ADD_NHIF"
	partnerAccountSetupActionName = "PARTNER_ACCOUNT_SETUP"
	completeProfileActionName     = "COMPLETE_PROFILE"
	hideItemActionName            = "HIDE_ITEM"
	pinItemActionName             = "PIN_ITEM"
	resolveItemActionName         = "RESOLVE_ITEM"
	helpActionName                = "GET_HELP"
	verifyEmailActionName         = "VERIFY_EMAIL"

	defaultOrg        = "default-org-id-please-change"
	defaultLocation   = "default-location-id-please-change"
	defaultContentDir = "/static/"
	defaultAuthor     = "Be.Well Team"
)

// embed default content assets (e.g images and documents) in the binary
var _ = pkger.Dir(defaultContentDir)

type actionGenerator func(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Action, error)

type nudgeGenerator func(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Nudge, error)

type itemGenerator func(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Item, error)

// SetDefaultActions ensures that a feed has default actions
func SetDefaultActions(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) ([]base.Action, error) {
	actions := []base.Action{}

	switch flavour {
	case base.FlavourConsumer:
		consumerActions, err := defaultConsumerActions(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default consumer actions: %w", err)
		}
		actions = consumerActions
	case base.FlavourPro:
		proActions, err := defaultProActions(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default pro actions: %w", err)
		}
		actions = proActions
	}

	return actions, nil
}

// SetDefaultNudges ensures that a feed has default nudges
func SetDefaultNudges(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) ([]base.Nudge, error) {
	var nudges []base.Nudge

	switch flavour {
	case base.FlavourConsumer:
		consumerNudges, err := defaultConsumerNudges(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default consumer nudges: %w", err)
		}
		nudges = consumerNudges
	case base.FlavourPro:
		proNudges, err := defaultProNudges(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default pro nudges: %w", err)
		}
		nudges = proNudges
	}

	return nudges, nil
}

// SetDefaultItems ensures that a feed has default feed items
func SetDefaultItems(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) ([]base.Item, error) {
	var items []base.Item

	switch flavour {
	case base.FlavourConsumer:
		consumerItems, err := defaultConsumerItems(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default consumer items: %w", err)
		}
		items = consumerItems
	case base.FlavourPro:
		proItems, err := defaultProItems(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default pro items: %w", err)
		}
		items = proItems
	}

	// fetch CMS items from the CMS feed tag
	cmsItems := feedItemsFromCMSFeedTag(ctx)
	items = append(items, cmsItems...)

	return items, nil
}

func defaultConsumerNudges(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) ([]base.Nudge, error) {
	var nudges []base.Nudge
	fns := []nudgeGenerator{
		addInsuranceNudge,
		addNHIFNudge,
		completeProfileNudge,
		verifyEmailNudge,
	}
	for _, fn := range fns {
		nudge, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating nudge: %w", err)
		}
		nudges = append(nudges, *nudge)
	}
	return nudges, nil
}

func defaultProNudges(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) ([]base.Nudge, error) {
	var nudges []base.Nudge
	fns := []nudgeGenerator{
		partnerAccountSetupNudge,
		completeProfileNudge,
		verifyEmailNudge,
	}
	for _, fn := range fns {
		nudge, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating nudge: %w", err)
		}
		nudges = append(nudges, *nudge)
	}
	return nudges, nil
}

func defaultConsumerActions(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) ([]base.Action, error) {
	var actions []base.Action
	fns := []actionGenerator{
		defaultHelpAction,
		defaultCoachingAction,
		defaultGetInsuranceAction,
		defaultGetTestAction,
		defaultBuyMedicineAction,
		defaultSeeDoctorAction,
	}
	for _, fn := range fns {
		action, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating action: %w", err)
		}
		actions = append(actions, *action)
	}
	return actions, nil
}

func defaultProActions(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) ([]base.Action, error) {
	var actions []base.Action
	fns := []actionGenerator{
		defaultAddPatientAction,
		defaultHelpAction,
		defaultSearchPatientAction,
	}
	for _, fn := range fns {
		action, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating action: %w", err)
		}
		actions = append(actions, *action)
	}
	return actions, nil
}

func defaultSeeDoctorAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getConsultationActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		StaticBase+"/actions/svg/see_doctor.svg",
		"See Doctor",
		"See a doctor",
		repository,
	)
}

func defaultBuyMedicineAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getMedicineActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		StaticBase+"/actions/svg/medicine.svg",
		"Get Medicine",
		"Get medicines",
		repository,
	)
}

func defaultGetTestAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getTestActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		StaticBase+"/actions/svg/get_tested.svg",
		"Get tests",
		"Get diagnostic tests",
		repository,
	)
}

func defaultGetInsuranceAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getInsuranceActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		StaticBase+"/actions/svg/buy_cover.svg",
		"Buy Cover",
		"Buy medical insurance",
		repository,
	)
}

func defaultCoachingAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		getCoachingActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		StaticBase+"/actions/svg/fitness.svg",
		"Coaching",
		"Get Health Coaching",
		repository,
	)
}

func defaultHelpAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		true,
		flavour,
		helpActionName,
		base.ActionTypeFloating,
		base.HandlingFullPage,
		StaticBase+"/actions/svg/help.svg",
		"Help",
		"Get Help",
		repository,
	)
}

func defaultSearchPatientAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		searchPatientActionName,
		base.ActionTypeSecondary,
		base.HandlingFullPage,
		StaticBase+"/actions/svg/search_user.svg",
		"Search a patient",
		"Search for a patient",
		repository,
	)
}

func defaultAddPatientAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Action, error) {
	return createGlobalAction(
		ctx,
		uid,
		false,
		flavour,
		addPatientActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		StaticBase+"/actions/svg/add_user.svg",
		"Register patient",
		"Register a patient",
		repository,
	)
}

func addInsuranceNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Nudge, error) {
	title := "Add Insurance"
	text := "Link your existing medical cover"
	imgURL := StaticBase + "/nudges/add_insurance.png"
	addInsuranceAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		addInsuranceActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", addInsuranceActionName, err)
	}
	actions := []base.Action{
		*addInsuranceAction,
	}
	return createNudge(
		ctx,
		uid,
		flavour,
		title,
		text,
		imgURL,
		title,
		text,
		actions,
		repository,
	)
}

func addNHIFNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Nudge, error) {
	title := "Add NHIF"
	text := "Link your NHIF cover"
	imgURL := StaticBase + "/nudges/add_insurance.png"
	addNHIFAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		addNHIFActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", addNHIFActionName, err)
	}
	actions := []base.Action{
		*addNHIFAction,
	}
	return createNudge(
		ctx,
		uid,
		flavour,
		title,
		text,
		imgURL,
		title,
		text,
		actions,
		repository,
	)
}

func partnerAccountSetupNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Nudge, error) {
	title := "Setup your partner account"
	text := "Create a partner account to begin transacting on Be.Well"
	imgURL := StaticBase + "/nudges/complete_profile.png"
	partnerAccountSetupAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		partnerAccountSetupActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", partnerAccountSetupActionName, err)
	}
	actions := []base.Action{
		*partnerAccountSetupAction,
	}
	return createNudge(
		ctx,
		uid,
		flavour,
		title,
		text,
		imgURL,
		title,
		text,
		actions,
		repository,
	)
}

func completeProfileNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Nudge, error) {
	title := "Complete your profile"
	text := "Fill in your Be.Well profile to unlock more rewards"
	imgURL := StaticBase + "/nudges/complete_profile.png"
	completeProfileAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		completeProfileActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", completeProfileActionName, err)
	}
	actions := []base.Action{
		*completeProfileAction,
	}
	return createNudge(
		ctx,
		uid,
		flavour,
		title,
		text,
		imgURL,
		title,
		text,
		actions,
		repository,
	)
}

func verifyEmailNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Nudge, error) {
	title := "Email Verification"
	text := "Please add and verify your email address"
	imgURL := StaticBase + "/nudges/verify_email.png"
	verifyEmailAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		verifyEmailActionName,
		base.ActionTypePrimary,
		base.HandlingFullPage,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create %s action: %w", verifyEmailActionName, err)
	}
	actions := []base.Action{
		*verifyEmailAction,
	}
	return createNudge(
		ctx,
		uid,
		flavour,
		title,
		text,
		imgURL,
		title,
		text,
		actions,
		repository,
	)
}

func createNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	title string,
	text string,
	imageURL string,
	imageTitle string,
	imageDescription string,
	actions []base.Action,
	repository Repository,
) (*base.Nudge, error) {
	future := time.Now().Add(time.Hour * futureHours)
	nudge := &base.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Visibility:     base.VisibilityShow,
		Status:         base.StatusPending,
		Expiry:         future,
		Title:          title,
		Text:           text,
		Links: []base.Link{
			base.GetPNGImageLink(
				imageURL, imageTitle, imageDescription, imageURL),
		},
		Actions:              actions,
		Groups:               []string{},
		Users:                []string{uid},
		NotificationChannels: []base.Channel{},
	}
	_, err := nudge.ValidateAndMarshal()
	if err != nil {
		return nil, fmt.Errorf("nudge validation error: %w", err)
	}

	nudge, err = repository.SaveNudge(ctx, uid, flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to save nudge: %w", err)
	}
	return nudge, nil
}

func createGlobalAction(
	ctx context.Context,
	uid string,
	allowAnonymous bool,
	flavour base.Flavour,
	name string,
	actionType base.ActionType,
	handling base.Handling,
	iconLink string,
	iconTitle string,
	iconDescription string,
	repository Repository,
) (*base.Action, error) {
	action := &base.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Name:           name,
		Icon: base.GetSVGImageLink(
			iconLink, iconTitle, iconDescription, iconLink),
		ActionType:     actionType,
		Handling:       handling,
		AllowAnonymous: allowAnonymous,
	}
	_, err := action.ValidateAndMarshal()
	if err != nil {
		return nil, fmt.Errorf("action validation error: %w", err)
	}

	action, err = repository.SaveAction(ctx, uid, flavour, action)
	if err != nil {
		return nil, fmt.Errorf("unable to save action: %w", err)
	}
	return action, nil
}

func createLocalAction(
	ctx context.Context,
	uid string,
	allowAnonymous bool,
	flavour base.Flavour,
	name string,
	actionType base.ActionType,
	handling base.Handling,
	repository Repository,
) (*base.Action, error) {
	action := &base.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Name:           name,
		Icon: base.GetPNGImageLink(
			StaticBase+"/1px.png",
			"Blank Image",
			"Default Blank Image",
			StaticBase+"/1px.png",
		),
		ActionType:     actionType,
		Handling:       handling,
		AllowAnonymous: allowAnonymous,
	}
	_, err := action.ValidateAndMarshal()
	if err != nil {
		return nil, fmt.Errorf("action validation error: %w", err)
	}
	// not saved...intentionally
	// it will save embedded in a nudge or feed item

	return action, nil
}

func createFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
	author string,
	tagline string,
	label string,
	iconImageURL string,
	iconTitle string,
	iconDescription string,
	summary string,
	text string,
	links []base.Link,
	actions []base.Action,
	conversations []base.Message,
	persistent bool,
	repository Repository,
) (*base.Item, error) {
	future := time.Now().Add(time.Hour * futureHours)
	item := &base.Item{
		ID:             itemID,
		SequenceNumber: defaultSequenceNumber,
		Expiry:         future,
		Persistent:     persistent,
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Icon: base.GetPNGImageLink(
			iconImageURL, iconTitle, iconDescription, iconImageURL),
		Author:               author,
		Tagline:              tagline,
		Label:                label,
		Timestamp:            time.Now(),
		Summary:              summary,
		Text:                 text,
		TextType:             base.TextTypeMarkdown,
		Links:                links,
		Actions:              actions,
		Conversations:        conversations,
		Groups:               []string{},
		Users:                []string{uid},
		NotificationChannels: []base.Channel{},
	}
	_, err := item.ValidateAndMarshal()
	if err != nil {
		return nil, fmt.Errorf("item validation error: %w", err)
	}
	item, err = repository.SaveFeedItem(ctx, uid, flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to save item: %w", err)
	}
	return item, nil
}

func defaultConsumerItems(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) ([]base.Item, error) {
	var items []base.Item
	fns := []itemGenerator{
		ultimateComposite,
		simpleConsumerWelcome,
	}
	for _, fn := range fns {
		item, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating item: %w", err)
		}
		items = append(items, *item)
	}
	return items, nil
}

func defaultProItems(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) ([]base.Item, error) {
	var items []base.Item
	fns := []itemGenerator{
		ultimateComposite,
		simpleProWelcome,
	}
	for _, fn := range fns {
		item, err := fn(ctx, uid, flavour, repository)
		if err != nil {
			return nil, fmt.Errorf("error when generating item: %w", err)
		}
		items = append(items, *item)
	}
	return items, nil
}

func simpleConsumerWelcome(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Item, error) {
	persistent := true // at least one persistent message in welcome data
	tagline := "Welcome to Be.Well"
	summary := "What is Be.Well?"
	text := "Be.Well is a virtual and physical healthcare community. Our goal is to make it easy for you to access affordable high-quality healthcare - whether online or in person."
	links := getFeedWelcomeVideos()
	actions, err := defaultActions(ctx, uid, flavour, repository)
	if err != nil {
		return nil, fmt.Errorf("can't initialize default actions: %w", err)
	}

	itemID := ksuid.New().String()
	conversations, err := getConsumerWelcomeThread(ctx, uid, flavour, itemID, repository)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize welcome message thread: %w", err)
	}

	return createFeedItem(
		ctx,
		uid,
		flavour,
		itemID,
		defaultAuthor,
		tagline,
		DefaultLabel,
		DefaultIconPath,
		"Feed Item Icon",
		"Feed Item Icon",
		summary,
		text,
		links,
		actions,
		conversations,
		persistent,
		repository,
	)
}

func simpleProWelcome(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Item, error) {
	persistent := true // at least one persistent message in welcome data
	tagline := "Welcome to Be.Well"
	summary := "What is Be.Well?"
	text := "Be.Well is a virtual and physical healthcare community. Our goal is to make it easy for you to provide affordable high-quality healthcare - whether online or in person."
	links := getFeedWelcomeVideos()
	actions, err := defaultActions(ctx, uid, flavour, repository)
	if err != nil {
		return nil, fmt.Errorf("can't initialize default actions: %w", err)
	}

	itemID := ksuid.New().String()
	conversations, err := getProWelcomeThread(ctx, uid, flavour, itemID, repository)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize welcome message thread: %w", err)
	}

	return createFeedItem(
		ctx,
		uid,
		flavour,
		itemID,
		defaultAuthor,
		tagline,
		DefaultLabel,
		DefaultIconPath,
		"Feed Item Icon",
		"Feed Item Icon",
		summary,
		text,
		links,
		actions,
		conversations,
		persistent,
		repository,
	)
}

func ultimateComposite(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) (*base.Item, error) {
	// here's what Be.Well can do for you... / help you do for your patients
	persistent := false
	tagline := "This is Be.Well..."
	summary := "This is Be.Well..."
	text := "This is Be.Well..."
	links := []base.Link{
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner01.png",
			"As within, so without",
			"As within, so without. Care begins with the self.",
			StaticBase+"/items/images/thumbs/bewell_banner01.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner02.png",
			"Have you made your bed today?",
			"Have you made your bed today? Your morning routine can get you started on a high note or not.",
			StaticBase+"/items/images/thumbs/bewell_banner02.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner03.png",
			"Smiles are contagious",
			"Smiles are contagious, pass it on and heal a heart",
			StaticBase+"/items/images/thumbs/bewell_banner03.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner04.png",
			"Tossing and turning?",
			"Too much screen time tampers with sleep",
			StaticBase+"/items/images/thumbs/bewell_banner04.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner05.png",
			"Wellness review",
			"You should have a wellness review daily",
			StaticBase+"/items/images/thumbs/bewell_banner05.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner06.png",
			"Wellness review",
			"You should have a wellness review daily",
			StaticBase+"/items/images/thumbs/bewell_banner06.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner07.png",
			"Healthcare Simplified",
			"HealthCare Simplified. Get Be.Well",
			StaticBase+"/items/images/thumbs/bewell_banner07.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner08.png",
			"Hot + Hot",
			"Did you know that Hot + Hot=Cool?",
			StaticBase+"/items/images/thumbs/bewell_banner08.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner09.png",
			"Does chewing gum improve focus?",
			"Does chewing gum improve focus? Yes or No. Tell us in the comments section.",
			StaticBase+"/items/images/thumbs/bewell_banner09.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner10.png",
			"Are you well?",
			"We take care of you with convenient teleconsultation.",
			StaticBase+"/items/images/thumbs/bewell_banner10.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner11.png",
			"Hands that heal",
			"Get Be.Well",
			StaticBase+"/items/images/thumbs/bewell_banner11.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner12.png",
			"Mindfulness matters",
			"Use Be.Well to consult",
			StaticBase+"/items/images/thumbs/bewell_banner12.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner13.png",
			"Be confident",
			"Be confident, you can consult in confidence. Use Be.Well.",
			StaticBase+"/items/images/thumbs/bewell_banner13.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner14.png",
			"Wellnes begins in the womb",
			"Wellness begins in the womb, we are here to walk with you",
			StaticBase+"/items/images/thumbs/bewell_banner14.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner15.png",
			"Self healing is...",
			"Self healing is...what techniques do you use to slow down and rejuvenate",
			StaticBase+"/items/images/thumbs/bewell_banner15.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner16.png",
			"Live. Love. Laugh. and Be.Well",
			"Live. Love. Laugh and Be.Well. Download Be.Well",
			StaticBase+"/items/images/thumbs/bewell_banner16.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner17.png",
			"Mindfulness matters",
			"Mindfulness matters. Use Be.Well to consult",
			StaticBase+"/items/images/thumbs/bewell_banner17.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner18.png",
			"Be.Well can transform you",
			"Be.Well can transform you. Get Be.Well",
			StaticBase+"/items/images/thumbs/bewell_banner18.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner19.png",
			"Need to refill your meds?",
			"Need to refill your meds? We will deliver",
			StaticBase+"/items/images/thumbs/tbewell_banner19.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner20.png",
			"Be optimistic, live longer",
			"Use Be.Well to help you manage ailments",
			StaticBase+"/items/images/thumbs/bewell_banner20.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner21.png",
			"Are you too anxious?",
			"Are you too anxious? Consult on Be.Well",
			StaticBase+"/items/images/thumbs/bewell_banner21.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner22.png",
			"Would you rather be taking tablets or syrup?",
			"Share your reasons in the comments section",
			StaticBase+"/items/images/thumbs/bewell_banner22.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner23.png",
			"Does coffee cure depression in women?",
			"How does coffee affect you?",
			StaticBase+"/items/images/thumbs/bewell_banner23.png",
		),
		base.GetPNGImageLink(
			StaticBase+"/items/images/bewell_banner24.png",
			"Is your allergy triggered by stress?",
			"Get Be.Well for more insights",
			StaticBase+"/items/images/thumbs/bewell_banner24.png",
		),

		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_25.pdf",
			"Does sunlight enable weight loss?",
			"Get Be.Well for meaningful insights",
			StaticBase+"/items/documents/thumbs/bewell_banner_25.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_26.pdf",
			"Anti-depressants and your libido.",
			"Consult your doctor on Be.Well",
			StaticBase+"/items/documents/thumbs/bewell_banner_26.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_27.pdf",
			"Need to refill your meds?",
			"We will deliver",
			StaticBase+"/items/documents/thumbs/bewell_banner_27.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_28.pdf",
			"Sexual wellness tips.",
			"Get Be.Well for the best tips",
			StaticBase+"/items/documents/thumbs/bewell_banner_28.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_29.pdf",
			"How do you get your kids to take meds?",
			"Do you negotiate or introduce Kiboko the motivator?",
			StaticBase+"/items/documents/thumbs/bewell_banner_29.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_30.pdf",
			"What are some of the meds that were part of you growing up?",
			"Get Be.Well",
			StaticBase+"/items/documents/thumbs/bewell_banner_30.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_31.pdf",
			"Mind your Mind",
			"Use Be.Well",
			StaticBase+"/items/documents/thumbs/bewell_banner_31.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_32.pdf",
			"Did you know that financial wellness affects your wellbeing as a whole?",
			"What else? Get Be.Well",
			StaticBase+"/items/documents/thumbs/bewell_banner_32.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_33.pdf",
			"How does technology preserve your health?",
			"Get Be.Well for meaningful insights",
			StaticBase+"/items/documents/thumbs/bewell_banner_33.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_34.pdf",
			"Your life can be transformed by an App",
			"Get Be.Well",
			StaticBase+"/items/documents/thumbs/bewell_banner_34.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_35.pdf",
			"It's not enough to just have a life transforming App",
			"You have to USE the App",
			StaticBase+"/items/documents/thumbs/bewell_banner_35.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_36.pdf",
			"Create your wellness habits",
			"Create your wellness habits and add your accountability partners on Be.Well",
			StaticBase+"/items/documents/thumbs/bewell_banner_36.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_37.pdf",
			"Patients don't need patience",
			"Dial your doctor, Now!",
			StaticBase+"/items/documents/thumbs/bewell_banner_37.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_38.pdf",
			"Simplifying and streamlining Healthcare",
			"Get Be.Well today!",
			StaticBase+"/items/documents/thumbs/bewell_banner_38.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_39.pdf",
			"Monitor your mental state",
			"Get Be.Well today!",
			StaticBase+"/items/documents/thumbs/bewell_banner_39.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_40.pdf",
			"Modern Medicine vs Mystic Medicine",
			"Modern Medicine vs Mystic Medicine. What works for you?",
			StaticBase+"/items/documents/thumbs/bewell_banner_40.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_41.pdf",
			"To Be.Well, your finances have to be well",
			"To Be.Well, your finances have to be well",
			StaticBase+"/items/documents/thumbs/bewell_banner_41.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_42.pdf",
			"Finance and Romance",
			"do they go hand in hand?",
			StaticBase+"/items/documents/thumbs/bewell_banner_42.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_43.pdf",
			"Let's talk about sex",
			"Let's talk about sex, the Be.Well way",
			StaticBase+"/items/documents/thumbs/bewell_banner_43.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_44.pdf",
			"Consult about sexual health in confidence",
			"Consult about sexual health in confidence",
			StaticBase+"/items/documents/thumbs/bewell_banner_44.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_45.pdf",
			"Is technology replacing your parenting?",
			"Is technology replacing your parenting?",
			StaticBase+"/items/documents/thumbs/bewell_banner_45.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_46.pdf",
			"Bottled water cures diseases",
			"Bottled water cures diseases. True or False?",
			StaticBase+"/items/documents/thumbs/bewell_banner_46.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_47.pdf",
			"Why fight for these gadgets?",
			"Why fight for these gadgets? Could it be a symptom or addiction?",
			StaticBase+"/items/documents/thumbs/bewell_banner_47.png",
		),
		base.GetPDFDocumentLink(
			StaticBase+"/items/documents/bewell_banner_48.pdf",
			"Put your money where your mind is",
			"Put your money where your mind is",
			StaticBase+"/items/documents/thumbs/bewell_banner_48.png",
		),
	}
	actions, err := defaultActions(ctx, uid, flavour, repository)
	if err != nil {
		return nil, fmt.Errorf("can't initialize default actions: %w", err)
	}

	conversations := []base.Message{}
	itemID := ksuid.New().String()
	return createFeedItem(
		ctx,
		uid,
		flavour,
		itemID,
		defaultAuthor,
		tagline,
		DefaultLabel,
		DefaultIconPath,
		"Feed Item Icon",
		"Feed Item Icon",
		summary,
		text,
		links,
		actions,
		conversations,
		persistent,
		repository,
	)
}

func getMessage(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
	text string,
	replyTo *base.Message,
	postedByName string,
	repository Repository,
) (*base.Message, error) {
	msg := &base.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: defaultSequenceNumber,
		Text:           text,
		PostedByUID:    defaultPostedByUID,
		PostedByName:   postedByName,
		Timestamp:      time.Now(),
	}
	if replyTo != nil {
		msg.ReplyTo = replyTo.ID
	}

	savedMsg, err := repository.PostMessage(ctx, uid, flavour, itemID, msg)
	if err != nil {
		return nil, fmt.Errorf("can't save message for default welcome thread(s): %w", err)
	}

	if savedMsg == nil {
		return nil, fmt.Errorf("nil saved message")
	}

	return savedMsg, nil
}

func getConsumerWelcomeThread(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
	repository Repository,
) ([]base.Message, error) {
	welcome, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"Welcome to Be.Well. We are glad to meet you!",
		nil,
		"Be.Well",
		repository,
	)
	if err != nil {
		return nil, err
	}

	pharmacyReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the medications service. I'll ensure that you get quality and affordable medications, on time. 👋!",
		welcome,
		"Medications Service",
		repository,
	)
	if err != nil {
		return nil, err
	}

	deliveryAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the delivery assistant. I help the medications service get medicines to you on time. 👋!",
		pharmacyReply,
		"Delivery Assistant",
		repository,
	)
	if err != nil {
		return nil, err
	}

	dispensingAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the dispensing assistant. I help your preferred pharmacy prepare your order before you go for it. 👋!",
		pharmacyReply,
		"Dispensing Assistant",
		repository,
	)
	if err != nil {
		return nil, err
	}

	testsReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the tests service. I'll ensure that you get quality and affordable diagnostic tests. 👋!",
		welcome,
		"Tests Service",
		repository,
	)
	if err != nil {
		return nil, err
	}

	consultationsReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the consultations service. I'll ensure that you can get in-person or remote(tele) advice from qualified medical professionals. 👋!",
		welcome,
		"Consultations Service",
		repository,
	)
	if err != nil {
		return nil, err
	}

	teleconsultAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the teleconsultations assistant. I'll ensure that you can reach a qualified medical professional via video or audio conference, whenever you need to. If you have an emergency, I'll help you find the nearest hospital for emergencies. 👋!",
		consultationsReply,
		"Teleconsultations Assistant",
		repository,
	)
	if err != nil {
		return nil, err
	}

	bookingAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the booking assistant. I'll help you book appointments for your care and remind you when it's time. 👋!",
		consultationsReply,
		"Booking Assistant",
		repository,
	)
	if err != nil {
		return nil, err
	}

	coachingReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the coaching service. I'll link you up to *awesome* wellness and fitness coaches. 👋!",
		welcome,
		"Coaching Service",
		repository,
	)
	if err != nil {
		return nil, err
	}

	insuranceReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the insurance service. I'll get you great quotes for medical cover and assist you when you need to use your insurance. 👋!",
		welcome,
		"Coaching Service",
		repository,
	)
	if err != nil {
		return nil, err
	}

	remindersReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the reminders service. I'll help you remember things related to your health. It could be an appointment or when you need to take some medication etc. Try me 👋!",
		welcome,
		"Reminders Service",
		repository,
	)
	if err != nil {
		return nil, err
	}

	return []base.Message{
		*welcome,
		*pharmacyReply,
		*deliveryAssistant,
		*dispensingAssistant,
		*testsReply,
		*consultationsReply,
		*teleconsultAssistant,
		*bookingAssistant,
		*coachingReply,
		*insuranceReply,
		*remindersReply,
	}, nil
}
func getProWelcomeThread(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
	repository Repository,
) ([]base.Message, error) {
	welcome, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"Welcome to Be.Well. We are glad to meet you!",
		nil,
		"Be.Well",
		repository,
	)
	if err != nil {
		return nil, err
	}

	pharmacyReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the medications service. I'll help you deliver quality and affordable medications, on time. 👋!",
		welcome,
		"Medications Service",
		repository,
	)
	if err != nil {
		return nil, err
	}

	deliveryAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the delivery assistant. I help the medications service deliver medicines on time. 👋!",
		pharmacyReply,
		"Delivery Assistant",
		repository,
	)
	if err != nil {
		return nil, err
	}

	dispensingAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the dispensing assistant. I help you prepare your orders. 👋!",
		pharmacyReply,
		"Dispensing Assistant",
		repository,
	)
	if err != nil {
		return nil, err
	}

	testsReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the tests service. I'll help you deliver quality and affordable diagnostic tests. 👋!",
		welcome,
		"Tests Service",
		repository,
	)
	if err != nil {
		return nil, err
	}

	consultationsReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the consultations service. I'll set up in-person and remote consultations for you. 👋!",
		welcome,
		"Consultations Service",
		repository,
	)
	if err != nil {
		return nil, err
	}

	teleconsultAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the teleconsultations assistant. I'll ensure that you can conduct consultations via video or audio conference, whenever you need to. If you have an emergency, I'll help you find the nearest hospital for emergencies. 👋!",
		consultationsReply,
		"Teleconsultations Assistant",
		repository,
	)
	if err != nil {
		return nil, err
	}

	bookingAssistant, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the booking assistant. I'll help you book appointments and remind you when it's time. 👋!",
		consultationsReply,
		"Booking Assistant",
		repository,
	)
	if err != nil {
		return nil, err
	}

	coachingReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the coaching service. I'll help you deliver your *awesome* coaching services to clients. 👋!",
		welcome,
		"Coaching Service",
		repository,
	)
	if err != nil {
		return nil, err
	}

	remindersReply, err := getMessage(
		ctx,
		uid,
		flavour,
		itemID,
		"I'm the reminders service. I'll help you remember things that you need to do. 👋!",
		welcome,
		"Reminders Service",
		repository,
	)
	if err != nil {
		return nil, err
	}

	return []base.Message{
		*welcome,
		*pharmacyReply,
		*deliveryAssistant,
		*dispensingAssistant,
		*testsReply,
		*consultationsReply,
		*teleconsultAssistant,
		*bookingAssistant,
		*coachingReply,
		*remindersReply,
	}, nil
}

func getFeedWelcomeVideos() []base.Link {
	return []base.Link{
		base.GetYoutubeVideoLink(
			"https://youtu.be/gcv2Z2AdpjM",
			"Be.Well lead",
			"Introducing Be.Well by Slade 360",
			StaticBase+"/items/videos/thumbs/01_lead.png",
		),
		base.GetYoutubeVideoLink(
			"https://youtu.be/W_daZjDET9Q",
			"Prescription delivery",
			"Get your medications delivered on Be.Well",
			StaticBase+"/items/videos/thumbs/02_prescription.png",
		),
		base.GetYoutubeVideoLink(
			"https://youtu.be/IbtVBXNvpSA",
			"Teleconsults",
			"Consult a doctor on Be.Well",
			StaticBase+"/items/videos/thumbs/03_teleconsult.png",
		),
		base.GetYoutubeVideoLink(
			"https://youtu.be/mKnlXcS3_Z0",
			"Slade 360",
			"Slade 360. HealthCare. Simplified.",
			StaticBase+"/items/videos/thumbs/04_slade.png",
		),
		base.GetYoutubeVideoLink(
			"https://youtu.be/XNDLnPfagLQ",
			"Mental health",
			"Mental health",
			StaticBase+"/items/videos/thumbs/05_mental_health.png",
		),
	}
}

func feedItemsFromCMSFeedTag(ctx context.Context) []base.Item {
	libraryService := library.NewService()
	items := []base.Item{}
	feedPosts, err := libraryService.GetFeedContent(ctx)
	if err != nil {
		//  non-fatal, intentionally
		log.Printf("ERROR: unable to fetch welcome feed posts from CMS: %s", err)
	}
	for _, post := range feedPosts {
		if post == nil {
			// non fatal, intentionally
			log.Printf("ERROR: nil CMS post when adding welcome posts to feed")
			continue
		}
		items = append(items, feedItemFromCMSPost(*post))
	}
	return items
}

func feedItemFromCMSPost(post library.GhostCMSPost) base.Item {
	future := time.Now().Add(time.Hour * futureHours)
	return base.Item{
		ID:                   post.UUID,
		SequenceNumber:       int(post.UpdatedAt.Unix()),
		Expiry:               future,
		Persistent:           false,
		Status:               base.StatusPending,
		Visibility:           base.VisibilityShow,
		Icon:                 base.GetPNGImageLink(DefaultIconPath, "Icon", "Feed Item Icon", DefaultIconPath),
		Author:               defaultAuthor,
		Tagline:              post.Slug,
		Label:                DefaultLabel,
		Summary:              TruncateStringWithEllipses(post.Excerpt, 140),
		Timestamp:            post.UpdatedAt,
		Text:                 post.HTML,
		TextType:             base.TextTypeHTML,
		Links:                getLinks(post),
		Actions:              []base.Action{},
		Conversations:        []base.Message{},
		Users:                []string{},
		Groups:               []string{},
		NotificationChannels: []base.Channel{},
	}
}

func getLinks(post library.GhostCMSPost) []base.Link {
	featureImageLink := post.FeatureImage
	if strings.HasSuffix(featureImageLink, ".png") {
		return []base.Link{
			{
				ID:       ksuid.New().String(),
				URL:      featureImageLink,
				LinkType: base.LinkTypePngImage,
			},
		}
	}
	return []base.Link{}
}

// TruncateStringWithEllipses truncates a string at the indicated length and adds trailing ellipses
func TruncateStringWithEllipses(str string, length int) string {
	if length <= 0 {
		return ""
	}

	targetLength := length
	addEllipses := false
	if length >= 140 {
		targetLength = length - 4 // room for ellipses for longer strings
		addEllipses = true
	}

	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= targetLength {
			break
		}
	}
	if addEllipses {
		return truncated + "..."
	}
	return truncated
}

func defaultActions(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	repository Repository,
) ([]base.Action, error) {
	resolveAction, err := createLocalAction(
		ctx,
		uid,
		false,
		flavour,
		resolveItemActionName,
		base.ActionTypePrimary,
		base.HandlingInline,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create resolve action: %w", err)
	}

	pinAction, err := createLocalAction(
		ctx,
		uid,
		true,
		flavour,
		pinItemActionName,
		base.ActionTypePrimary,
		base.HandlingInline,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create pin action: %w", err)
	}

	hideAction, err := createLocalAction(
		ctx,
		uid,
		true,
		flavour,
		hideItemActionName,
		base.ActionTypePrimary,
		base.HandlingInline,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create hide action: %w", err)
	}
	actions := []base.Action{
		*resolveAction,
		*pinAction,
		*hideAction,
	}

	return actions, nil
}
