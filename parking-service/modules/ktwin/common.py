import os
import json
from cloudevents.http import CloudEvent

EVENT_TYPE_REAL_GENERATED = "ktwin.real.{0}"
EVENT_TYPE_VIRTUAL_GENERATED = "ktwin.virtual.{0}"
EVENT_TYPE_COMMAND_EXECUTED = "ktwin.command.{0}"

def get_event_store_url():
    return os.getenv("KTWIN_EVENT_STORE")

def get_broker_url():
    return os.getenv("KTWIN_BROKER")

# KTWIN Events
class KTwinEvent:
    def __init__(self, cloud_event: CloudEvent):
        self.cloud_event = cloud_event
        self.twin_interface = None
        ce_type_split = cloud_event["type"].split(".")
        if len(ce_type_split) > 2:
            self.twin_interface = ce_type_split[2]
        self.twin_instance = cloud_event["source"]

    def set_event(self, twin_interface: str, twin_instance: str, data):
        self.twin_interface = twin_interface
        self.twin_instance = twin_instance
        ce_type = "ktwin." + self.twin_interface
        self.cloud_event = build_cloud_event(ce_type=ce_type, ce_source=twin_instance, data=data)

def build_cloud_event(ce_type, ce_source, data):
    attributes = {
        "type" : ce_type,
        "source" : ce_source
    }
    return CloudEvent(attributes, data)

# KTWIN Command Events
class KTwinCommandEvent:
    def __init__(self, cloud_event: CloudEvent):
        self.cloud_event = cloud_event
        self.twin_interface = None
        self.command = None
        ce_type_split = cloud_event["type"].split(".")
        if len(ce_type_split) > 2:
            self.twin_interface = ce_type_split[2]
        if len(ce_type_split) > 3:
            self.command = ce_type_split[3]
        self.twin_instance_source = cloud_event["source"]

# Twin Graph Components

class TwinInstanceReference:
    def __init__(self, name: str, interface: str, instance: str):
        self.name = name
        self.interface = interface
        self.instance = instance

    def toJSON(self):
        return json.dumps(self, default=lambda o: o.__dict__, sort_keys=True, indent=4)


class TwinInstanceGraph:
    def __init__(self, name: str, interface: str, relationships: list[TwinInstanceReference]):
        self.name = name
        self.interface = interface
        self.relationships = relationships

    def toJSON(self):
        return json.dumps(self, default=lambda o: o.__dict__, sort_keys=True, indent=4)


class TwinGraph:
    def __init__(self, twin_instances_graph: dict[str, TwinInstanceGraph]):
        self.twin_instances_graph = twin_instances_graph

    def toJSON(self):
        return json.dumps(self, default=lambda o: o.__dict__, sort_keys=True, indent=4)
