# SendGrid

Event web hook consumer for SendGrid's
[eventwebhook](https://sendgrid.com/docs/for-developers/tracking-events/getting-started-event-webhook/)
events. It accepts all events and stores the last N events in memory,
and serves the recent events for inspection later.

## How to configure

Configure SendGrid to point to the `/eventwebhook` endpoint. To view
recent events visit the `/recent` endpoint.

| envvar             | desc                                                      |
|--------------------|-----------------------------------------------------------|
| PRETTY_PRINT       | Pretty prints log output. JSON is always served condensed |
| RECENT_EVENT_COUNT | Maximum number of events to hold in memory.               |
| PORT               | What port to serve on.                                    |

```
export PRETTY_PRINT=true
export RECENT_EVENT_COUNT=200
export PORT=5555
go run cmd/eventwebhook/main.go
```
