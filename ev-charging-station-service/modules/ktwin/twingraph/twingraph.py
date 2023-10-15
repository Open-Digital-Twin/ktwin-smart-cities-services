import os
import json
from ..common import TwinGraph, TwinReference

# Twin Graph methods

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

# Get the Graph relationship by name and instance
def get_relationship_from_graph(twin_instance: str, relationship_name: str, twin_graph: dict[str, TwinGraph]) -> TwinReference:
    graph_node = twin_graph[twin_instance]
    
    for relationship in graph_node.relationships:
        if relationship.name == relationship_name:
            return relationship
        
    return None

# Get Twin Graph Node by twin instance and interface
def get_twin_graph_by_relationship(relationship_twin_instance: str, relationship_twin_interface: str, twin_graph: dict[str, TwinGraph]) -> TwinReference:
    for twin_instance in twin_graph:
        graph_node = twin_graph[twin_instance]
        for relationship in graph_node.relationships:
            if relationship.twin_instance == relationship_twin_instance and relationship.twin_interface == relationship_twin_interface:
                return graph_node
        
    return None