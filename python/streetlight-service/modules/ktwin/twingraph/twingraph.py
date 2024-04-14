import os
import json
import requests
from ..common import TwinGraph, TwinInstanceReference, TwinInstanceGraph

# Twin Graph methods

def get_twin_graph(twin_instance: str) -> dict:
    if os.getenv("ENV") == "local":
        ktwin_graph = json.loads(os.getenv("KTWIN_GRAPH"))
    else:
        ktwin_graph_url = os.getenv("KTWIN_GRAPH_URL")
        response = requests.get(ktwin_graph_url + "/" + twin_instance)

        if response.status_code == 404:
            return dict()

        if response.status_code != 200:
            raise Exception("Error while calling service status_code: " + str(response.status_code))

        ktwin_graph = response.json()
    
    return ktwin_graph

def load_twin_graph_by_instance(twin_instances: list[str]) -> TwinGraph:
    ktwin_graph_list = list()
    for twin_instance in twin_instances:
        ktwin_graph = get_twin_graph(twin_instance=twin_instance)
        if "twinInstances" in ktwin_graph:
            ktwin_graph_list.append(ktwin_graph)

    if len(ktwin_graph_list) == 0:
        return dict()
        
    twin_instances_graph: dict[str, TwinInstanceGraph] = dict()
    for ktwin_graph in ktwin_graph_list:
        for twin_instance_graph in ktwin_graph["twinInstances"]:
            relationship_list: list[TwinInstanceReference] = list()
            if "relationships" in twin_instance_graph:
                for twin_relationship in twin_instance_graph["relationships"]:
                    relationship = TwinInstanceReference(name=twin_relationship["name"], interface=twin_relationship["interface"], instance=twin_relationship["instance"])
                    relationship_list.append(relationship)
            twin_instances_graph[twin_instance_graph["name"]] = TwinInstanceGraph(interface=twin_instance_graph["interface"], name=twin_instance_graph["name"], relationships=relationship_list)

    ktwin_graph = TwinGraph(twin_instances_graph=twin_instances_graph)
    write_twin_graph(ktwin_graph=ktwin_graph)
    return ktwin_graph

def load_twin_graph() -> TwinGraph:
    if os.getenv("ENV") == "local":
        ktwin_graph = json.loads(os.getenv("KTWIN_GRAPH"))
    else:
        ktwin_graph_url = os.getenv("KTWIN_GRAPH_URL")
        response = requests.get(ktwin_graph_url)

        if response.status_code != 200:
            raise Exception("Error while calling service status_code: " + str(response.status_code))

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

    ktwin_graph = TwinGraph(twin_instances_graph=twin_instances_graph)
    write_twin_graph(ktwin_graph=ktwin_graph)
    return ktwin_graph

def write_twin_graph(ktwin_graph: TwinGraph):
    # Write Twin Graph in the local system
    f = open("ktwin_graph.json", "w")
    f.write(ktwin_graph.toJSON())
    f.close()

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
def get_twin_graph_by_relation(target_twin_interface: str, source_twin_instance: str, twin_graph: TwinGraph) -> TwinInstanceReference:
    if source_twin_instance not in twin_graph.twin_instances_graph:
        raise Exception("Twin Source not available in TwinGraph: " + source_twin_instance)

    source_twin_graph = twin_graph.twin_instances_graph[source_twin_instance]

    for relationship in source_twin_graph.relationships:
        if relationship.interface == target_twin_interface:
            return relationship

    return None