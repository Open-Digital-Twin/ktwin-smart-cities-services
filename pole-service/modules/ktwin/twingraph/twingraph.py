import os
import json
import requests
from ..common import EVENT_TYPE_VIRTUAL_GENERATED, build_cloud_event, get_broker_url
from cloudevents.http import to_structured


# Twin Graph methods
class TwinReference:
    def __init__(self, name: str, twin_interface: str , twin_instance: str) -> None:
        self.name = name
        self.twin_interface = twin_interface
        self.twin_instance = twin_instance

class TwinGraph:
    def __init__(self, relationships: list[TwinReference]) -> None:
        self.relationships = relationships


def load_twin_graph() -> dict[TwinGraph]:
    ktwin_graph = os.getenv("KTWIN_GRAPH")
    if ktwin_graph is None:
        return ktwin_graph
    graph_json = json.loads(ktwin_graph)

    twin_graph = dict()
    for twin_instance_key in graph_json:
        twin_instance_graph = graph_json[twin_instance_key]

        relationship_list = list()
        for twin_relationship in twin_instance_graph["relationships"]:
            reference = TwinReference(name=twin_relationship["name"], twin_interface=["interface"], twin_instance=twin_relationship["instance"])
            relationship_list.append(reference)
        twin_graph[twin_instance_key] = TwinGraph(relationships=relationship_list)

    return twin_graph

def push_event_to_relationship(twin_instance: str, relationship_name: str, twin_graph: dict[str, TwinGraph], payload_data: dict) -> TwinReference:
    relationship = get_relationship_reference(twin_instance=twin_instance, relationship_name=relationship_name, twin_graph=twin_graph)

    if relationship != None:
        ce_type = EVENT_TYPE_VIRTUAL_GENERATED.format(relationship.twin_interface)
        ce_source = twin_instance
        cloud_event = build_cloud_event(ce_type, ce_source, payload_data)
        headers, body = to_structured(cloud_event)
        response = requests.post(get_broker_url(), headers=headers, data=body)

        if response.status_code != 202:
            raise Exception("Error when pushing to event broker", response)

def get_relationship_reference(twin_instance: str, relationship_name: str, twin_graph: dict[str, TwinGraph]) -> TwinReference:
    graph_node = twin_graph[twin_instance]
    
    for relationship in graph_node.relationships:
        if relationship.name == relationship_name:
            return relationship
        
    return None