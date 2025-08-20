package publisher

import (
	"fmt"
	"os"
)

func Select() (Publisher, string) {
	switch os.Getenv("PUBLISHER") {
	case "pubsub":
		return PubSubPublisher{}, "pubsub"
	case "devpush":
		return DevPushPublisher{}, "devpush"
	case "nop", "":
		return NopPublisher{}, "nop"
	default:
		return NopPublisher{}, fmt.Sprintf("unknown:%s", os.Getenv("PUBLISHER"))
	}
}

