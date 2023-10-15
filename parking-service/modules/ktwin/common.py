import os
from cloudevents.http import CloudEvent

EVENT_TYPE_REAL_GENERATED = "ktwin.real.%s"
EVENT_TYPE_VIRTUAL_GENERATED = "ktwin.virtual.%s"
EVENT_TYPE_COMMAND_EXECUTED = "ktwin.command.%s"

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
        self.twin_instance = cloud_event["source"]

# Twin Graph Components
class TwinReference:
    def __init__(self, name: str, twin_interface: str , twin_instance: str) -> None:
        self.name = name
        self.twin_interface = twin_interface
        self.twin_instance = twin_instance

class TwinGraph:
    def __init__(self, relationships: list[TwinReference]) -> None:
        self.relationships = relationships
