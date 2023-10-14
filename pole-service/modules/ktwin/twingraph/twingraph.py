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

def get_relationship_from_graph(twin_instance: str, relationship_name: str, twin_graph: dict[str, TwinGraph]) -> TwinReference:
    graph_node = twin_graph[twin_instance]
    
    for relationship in graph_node.relationships:
        if relationship.name == relationship_name:
            return relationship
        
    return None
