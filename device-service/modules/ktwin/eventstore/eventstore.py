import requests
from ..common import get_event_store_url, KTwinEvent
from cloudevents.http import from_http, to_binary

# Event Store Methods
# Interact with event store

def get_latest_twin_event(twin_interface, twin_instance):
    url = get_event_store_url() + "/api/v1/twin-events/%s/%s/latest" % (twin_interface, twin_instance)
    response = requests.get(url)

    if response.status_code == 404:
        return None

    cloud_event = from_http(response.headers, response.content)
    return KTwinEvent(cloud_event)

def update_twin_event(ktwin_event: KTwinEvent):
    url = get_event_store_url() + "/api/v1/twin-events"
    headers, body = to_binary(ktwin_event.cloud_event)
    response = requests.post(url, data=body, headers=headers)

    if response.status_code != 202:
        raise Exception("Error while updating twin event", response)

    return response