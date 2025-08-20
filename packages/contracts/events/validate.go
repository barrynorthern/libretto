package events

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	eventsv1 "github.com/barrynorthern/libretto/gen/go/libretto/events/v1"
	"github.com/google/uuid"
)

var semverRe = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

// ValidateEnvelope performs lightweight semantic validation on a protobuf Envelope.
// It checks required fields, basic formats (semver, UUID), and occurred_at presence.
func ValidateEnvelope(env *eventsv1.Envelope) error {
	if env == nil {
		return errors.New("envelope is nil")
	}
	if strings.TrimSpace(env.GetEventName()) == "" {
		return errors.New("eventName is required")
	}
	if !semverRe.MatchString(env.GetEventVersion()) {
		return fmt.Errorf("eventVersion must be semver (x.y.z): %q", env.GetEventVersion())
	}
	if err := mustUUID("eventId", env.GetEventId()); err != nil { return err }
	if err := mustUUID("correlationId", env.GetCorrelationId()); err != nil { return err }
	if err := mustUUID("causationId", env.GetCausationId()); err != nil { return err }
	if strings.TrimSpace(env.GetIdempotencyKey()) == "" {
		return errors.New("idempotencyKey is required")
	}
	if strings.TrimSpace(env.GetProducer()) == "" {
		return errors.New("producer is required")
	}
	if strings.TrimSpace(env.GetTenantId()) == "" {
		return errors.New("tenantId is required")
	}
	if env.GetOccurredAt() == nil || env.GetOccurredAt().AsTime().IsZero() {
		return errors.New("occurredAt is required")
	}
	return nil
}

func mustUUID(field, v string) error {
	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("%s is required", field)
	}
	if _, err := uuid.Parse(v); err != nil {
		return fmt.Errorf("%s must be a UUID: %v", field, err)
	}
	return nil
}

