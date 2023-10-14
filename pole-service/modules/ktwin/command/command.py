import requests
from ..common import EVENT_TYPE_COMMAND_EXECUTED, build_cloud_event, get_broker_url, KTwinCommandEvent
from cloudevents.http import to_structured, from_http

# Command: is the name of the command that will be invoked in the target interface
# CommandPayload: is the command payload that is going to be sent to the target instance
# RelationshipName: the relationship name in the graph that will receive the command. TwinInstance and Interface are populated based the relationship information of the KTWIN Graph
# Source: the source Twin Instance that is generating the event.
def invoke_command(command: str, commandPayload: dict, relationshipName: str, twin_instance: str):
    ce_type = EVENT_TYPE_COMMAND_EXECUTED.format(twin_interface)
    ce_source = twin_instance
    cloud_event = build_cloud_event(ce_type, ce_source, commandPayload)
    headers, body = to_structured(cloud_event)
    response = requests.post(get_broker_url(), headers=headers, data=body)

    if response.status_code != 202:
        raise Exception("Error when pushing to event broker", response)

def handle_command(request: requests.Request, twin_interface: str, command: str, callback):
    ktwin_command_event = handle_request(request)
    if ktwin_command_event.twin_interface == twin_interface and ktwin_command_event.command == command:
        callback(ktwin_command_event)

def handle_request(request) -> KTwinCommandEvent:
    cloud_event = from_http(request.headers, request.get_data())
    return KTwinCommandEvent(cloud_event)