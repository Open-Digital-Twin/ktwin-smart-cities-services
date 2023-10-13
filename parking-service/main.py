import os
import datetime
import sys
import logging
from dotenv import load_dotenv
from flask import Flask, request
from modules.ktwin import handle_request, handle_event, KTwinEvent, Twin, get_latest_twin_event, update_twin_event, get_parent_twins, push_to_virtual_twin

if os.getenv("ENV") == "local":
    load_dotenv('local.env')

app = Flask(__name__)

handler = logging.StreamHandler(sys.stdout)
handler.setFormatter(logging.Formatter(
    '%(asctime)s - %(name)s - %(levelname)s - %(message)s'))
app.logger.addHandler(handler)
app.logger.setLevel(logging.INFO)

@app.route("/", methods=["POST"])
def home():
    event = handle_request(request)

    app.logger.info(
        f"Event TwinInstance: {event.twin_instance} - Event TwinInterface: {event.twin_interface}"
    )

    handle_event(request, 'ngsi-ld-city-parkingspot', handle_parkingspot_event)
    handle_event(request, 'ngsi-ld-city-offstreetparking', handle_offstreetparking_event)

    # Return 204 - No-content
    return "", 204

def handle_parkingspot_event(event: KTwinEvent):
    current_parkingspot_event = event.cloud_event.data

    if "status" not in current_parkingspot_event:
        app.logger.info(f"Event {event.cloud_event} has no status attribute value")
    else:
        parkingspot_status = current_parkingspot_event["status"]
        parkingspot_category = current_parkingspot_event["category"]
        update_twin_event(event)

        if parkingspot_category == "offStreet":
            event_data = dict()
            event_data["vehicleEntranceCount"]
            if parkingspot_status == "occupied":
                event_data["vehicleEntranceCount"] = 1
                event_data["vehicleExitCount"] = -1
                app.logger.info(f"Generate event to decrement number of available slots")

            if parkingspot_status == "free":
                event_data["vehicleEntranceCount"] = -1
                event_data["vehicleExitCount"] = 1
                app.logger.info(f"Generate event to increment number of available slots")

            parent_twins = get_parent_twins()
            if (parent_twins) > 0:
                parent_twins
            else:
                app.logger.info(f"No available parent entities")


def handle_offstreetparking_event(event: KTwinEvent):
    current_offstreetparking_event = event.cloud_event.data
    latest_offstreetparking_event = get_latest_twin_event(event.twin_interface, event.twin_instance)
    if latest_offstreetparking_event is None:
        latest_offstreetparking_event = current_offstreetparking_event

    if "vehicleEntranceCount" not in current_offstreetparking_event:
        app.logger.info(f"Event {event.cloud_event} has no status attribute value")
    else:
        latest_vehicleEntranceCount = 0
        if "vehicleEntranceCount" in latest_offstreetparking_event:
            latest_vehicleEntranceCount = latest_offstreetparking_event["vehicleEntranceCount"]
        current_offstreetparking_event["vehicleEntranceCount"] = current_offstreetparking_event["vehicleEntranceCount"] + latest_vehicleEntranceCount
        event.cloud_event = current_offstreetparking_event
    
    current_offstreetparking_event = event.cloud_event.data
    if "vehicleExitCount" not in current_offstreetparking_event:
        app.logger.info(f"Event {event.cloud_event} has no status attribute value")
    else:
        latest_vehicleExitCount = 0
        if "vehicleExitCount" in latest_offstreetparking_event:
            latest_vehicleExitCount = latest_offstreetparking_event["vehicleExitCount"]
        current_offstreetparking_event["vehicleExitCount"] = current_offstreetparking_event["vehicleExitCount"] + latest_vehicleExitCount
        event.cloud_event.data = current_offstreetparking_event

    update_twin_event(event)

if __name__ == "__main__":
    app.logger.info("Starting up server...")
    app.run(host='0.0.0.0', port=8080)