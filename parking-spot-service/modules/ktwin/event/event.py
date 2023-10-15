import os
import requests
from ..common import EVENT_TYPE_VIRTUAL_GENERATED, EVENT_TYPE_REAL_GENERATED, get_broker_url, KTwinEvent, build_cloud_event
from cloudevents.http import to_structured, from_http


# Event Handling
# Handle Real-Virtual and Virtual-Real Twin communications

def send_to_real_twin(twin_interface, twin_instance, data):
    ce_type = EVENT_TYPE_VIRTUAL_GENERATED.format(twin_interface)
    ce_source = twin_instance
    cloud_event = build_cloud_event(ce_type, ce_source, data)
    headers, body = to_structured(cloud_event)

    if os.getenv("ENV") != "local":
        response = requests.post(get_broker_url(), headers=headers, data=body)

        if response.status_code != 202:
            raise Exception("Error when pushing to event broker", response)

def push_to_virtual_twin(twin_interface, twin_instance, data):
    ce_type = EVENT_TYPE_REAL_GENERATED.format(twin_interface)
    ce_source = twin_instance
    cloud_event = build_cloud_event(ce_type, ce_source, data)
    headers, body = to_structured(cloud_event)

    if os.getenv("ENV") != "local":
        response = requests.post(get_broker_url(), headers=headers, data=body)

        if response.status_code != 202:
            raise Exception("Error when pushing to event broker", response)

def handle_event(request: requests.Request, twin_interface: str, callback):
    ktwin_event = handle_request(request)
    if ktwin_event.twin_interface == twin_interface:
        callback(ktwin_event)

def handle_request(request) -> KTwinEvent:
    cloud_event = from_http(request.headers, request.get_data())
    return KTwinEvent(cloud_event)