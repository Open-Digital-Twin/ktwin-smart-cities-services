import os
import json
import requests
from ..common import TwinGraph, TwinInstanceReference, TwinInstanceGraph

# Twin Graph methods

def load_twin_graph() -> TwinGraph:
    ktwin_graph_url = os.getenv("KTWIN_GRAPH_URL")
    response = requests.get(ktwin_graph_url)

    if response.status_code != 200:
        raise Exception("Error while calling service status_code: " + response.status_code)

    ktwin_graph = response.json()

    if "twinInstances" not in ktwin_graph:
        return dict()

    twin_instances_graph: dict[str, TwinInstanceGraph] = dict()
    for twin_instance_graph in ktwin_graph["twinInstances"]:
        relationship_list: list[TwinInstanceReference] = list()
        if "relationships" in twin_instance_graph:
            for twin_relationship in twin_instance_graph["relationships"]:
                relationship = TwinInstanceReference(name=twin_relationship["name"], interface=twin_relationship["interface"], instance=twin_relationship["instance"])
                relationship_list.append(relationship)
        twin_instances_graph[twin_instance_graph["name"]] = TwinInstanceGraph(interface=twin_instance_graph["interface"], name=twin_instance_graph["name"], relationships=relationship_list)

    return TwinGraph(twin_instances_graph=twin_instances_graph)

# Get the Graph relationship by name and instance
def get_relationship_from_graph(twin_instance: str, relationship_name: str, twin_graph: TwinGraph) -> TwinInstanceReference:
    for instance in twin_graph.twin_instances_graph:
        twin_instance_graph = twin_graph.twin_instances_graph[instance]
        if twin_instance_graph.name == twin_instance:
            for relationship in twin_instance_graph.relationships:
                if relationship.name == relationship_name:
                    return relationship

    return None

# Get Twin Graph Node by twin instance and interface
def get_twin_graph_by_relationship(relationship_twin_instance: str, relationship_twin_interface: str, twin_graph: TwinGraph) -> TwinInstanceReference:
    for twin_instance in twin_graph.twin_instances_graph:
        twin_instance_graph = twin_graph.twin_instances_graph[twin_instance]
        for relationship in twin_instance_graph.relationships:
            if relationship.twin_instance == relationship_twin_instance and relationship.twin_interface == relationship_twin_interface:
                return relationship
        
    return None