import requests
from ..common import EVENT_TYPE_COMMAND_EXECUTED, build_cloud_event, get_broker_url, KTwinCommandEvent
from ..twingraph import get_relationship_from_graph, load_twin_graph
from cloudevents.http import to_structured, from_http

# command: the name of the command that will be invoked in the target interface.
# command_payload: is the command payload that is going to be sent to the target interface.
# relationship_name: the relationship name in the graph that will receive the command.
# source: the source twin instance that is generating the event.
def execute_command(command: str, command_payload: dict, relationship_name: str, twin_instance: str):
    twin_graph = load_twin_graph()
    relationship = get_relationship_from_graph(twin_instance=twin_instance, relationship_name=relationship_name, twin_graph=twin_graph)
    if relationship is None:
        raise ValueError("Relationship not exists")
    ce_type = EVENT_TYPE_COMMAND_EXECUTED.format(relationship.twin_interface + "." + command)
    ce_source = twin_instance
    cloud_event = build_cloud_event(ce_type=ce_type, ce_source=ce_source, data=command_payload)
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