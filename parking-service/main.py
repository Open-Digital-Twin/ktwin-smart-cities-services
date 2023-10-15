import os
import datetime
import sys
import logging
from dotenv import load_dotenv
from flask import Flask, request
import modules.ktwin.event as kevent
import modules.ktwin.eventstore as keventstore
import modules.ktwin.command as kcommand
import modules.ktwin.twingraph as ktwingraph

if os.getenv("ENV") == "local":
    load_dotenv('local.env')

app = Flask(__name__)

handler = logging.StreamHandler(sys.stdout)
handler.setFormatter(logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s'))
app.logger.addHandler(handler)
app.logger.setLevel(logging.INFO)

@app.route("/", methods=["POST"])
def home():
    event = kevent.handle_request(request)

    app.logger.info(
        f"Event TwinInstance: {event.twin_instance} - Event TwinInterface: {event.twin_interface}"
    )

    kevent.handle_event(request, 'ngsi-ld-city-parkingspot', handle_parkingspot_event)
    kcommand.handle_command(request=request, twin_interface='ngsi-ld-city-offstreetparking', command='updateVehicleCount', callback=handle_update_vehicle_count_command)

    # Return 204 - No-content
    return "", 204

# Set the Parking Spot to Occupied or free
# Generate Event to Off Street Parking to update the number of occupied and free slots.
def handle_parkingspot_event(event: kevent.KTwinEvent):
    current_parkingspot_event = event.cloud_event.data

    if "status" not in current_parkingspot_event:
        app.logger.info(f"Event {event.cloud_event} has no status attribute value")
    else:
        parkingspot_status = current_parkingspot_event["status"]
        parkingspot_category = current_parkingspot_event["category"]
        keventstore.update_twin_event(event)

        if parkingspot_category == "offStreet":
            event_data = dict()
            event_data["vehicleEntranceCount"]
            if parkingspot_status == "occupied":
                command_payload = dict()
                command_payload["vehicleEntranceCount"] = 1
                command_payload["vehicleExitCount"] = -1
                kcommand.execute_command(command="updateVehicleCount", command_payload=command_payload, relationship_name="refOffStreetParking", twin_instance_source=event.twin_instance)

            if parkingspot_status == "free":
                command_payload = dict()
                command_payload["vehicleEntranceCount"] = -1
                command_payload["vehicleExitCount"] = 1
                kcommand.execute_command(command="updateVehicleCount", command_payload=command_payload, relationship_name="refOffStreetParking", twin_instance_source=event.twin_instance)

def handle_update_vehicle_count_command(command_event: kcommand.KTwinCommandEvent, target_twin_instance: ktwingraph.TwinInstanceReference):
    current_offstreetparking_data = command_event.cloud_event.data
    latest_offstreetparking_event = keventstore.get_latest_twin_event(twin_interface=target_twin_instance.twin_interface, twin_instance=target_twin_instance.twin_instance)
    latest_offstreetparking_data = None

    if latest_offstreetparking_event is None:
        latest_offstreetparking_data = current_offstreetparking_data
    else:
        latest_offstreetparking_data = latest_offstreetparking_event.cloud_event.data

    if "vehicleEntranceCount" not in current_offstreetparking_data:
        app.logger.info(f"Event {command_event.cloud_event} has no status attribute value")
    else:
        latest_vehicleEntranceCount = 0
        if "vehicleEntranceCount" in latest_offstreetparking_data:
            latest_vehicleEntranceCount = latest_offstreetparking_data["vehicleEntranceCount"]
        current_offstreetparking_data["vehicleEntranceCount"] = current_offstreetparking_data["vehicleEntranceCount"] + latest_vehicleEntranceCount
        command_event.cloud_event.data = current_offstreetparking_data
    
    current_offstreetparking_data = command_event.cloud_event.data
    if "vehicleExitCount" not in current_offstreetparking_data:
        app.logger.info(f"Event {command_event.cloud_event} has no status attribute value")
    else:
        latest_vehicleExitCount = 0
        if "vehicleExitCount" in latest_offstreetparking_data:
            latest_vehicleExitCount = latest_offstreetparking_data["vehicleExitCount"]
        current_offstreetparking_data["vehicleExitCount"] = current_offstreetparking_data["vehicleExitCount"] + latest_vehicleExitCount
        command_event.cloud_event.data = current_offstreetparking_data

    keventstore.update_twin_event(command_event)

if __name__ == "__main__":
    app.logger.info("Starting up server...")
    app.run(host='0.0.0.0', port=8080)
