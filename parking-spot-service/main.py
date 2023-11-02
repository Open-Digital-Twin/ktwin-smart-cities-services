import os
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

ktwin_graph = ktwingraph.load_twin_graph()

@app.route("/", methods=["POST"])
def home():
    event = kevent.handle_request(request)

    app.logger.info(
        f"Event TwinInstance: {event.twin_instance} - Event TwinInterface: {event.twin_interface}"
    )

    kevent.handle_event(request, 'ngsi-ld-city-parkingspot', handle_parkingspot_event)

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

        if "category" not in current_parkingspot_event:
            latest_parkingspot_event = keventstore.get_latest_twin_event(twin_instance=event.twin_instance, twin_interface=event.twin_interface)

            if latest_parkingspot_event is None or "category" not in latest_parkingspot_event.cloud_event.data:
                # If not provided, the default is offStreet
                # TODO: set the category in the moment creation
                parkingspot_category = "offStreet"
            else:
                parkingspot_category = latest_parkingspot_event["category"]

        else:
            parkingspot_category = current_parkingspot_event["category"]

        keventstore.update_twin_event(event)

        if parkingspot_category == "offStreet":
            if parkingspot_status == "occupied":
                command_payload = dict()
                command_payload["vehicleEntranceCount"] = 1
                try:
                    kcommand.execute_command(command="updateVehicleCount", command_payload=command_payload, relationship_name="refOffStreetParking", twin_instance_source=event.twin_instance, twin_graph=ktwin_graph)
                except Exception as error:
                    app.logger.error(f"Error to execute command updateVehicleCount in relationship refOffStreetParking in TwinInstance {event.twin_instance}")
                    app.logger.error(error)

            if parkingspot_status == "free":
                command_payload = dict()
                command_payload["vehicleExitCount"] = 1
                try:
                    kcommand.execute_command(command="updateVehicleCount", command_payload=command_payload, relationship_name="refOffStreetParking", twin_instance_source=event.twin_instance, twin_graph=ktwin_graph)
                except Exception as error:
                    app.logger.error(f"Error to execute command updateVehicleCount in relationship refOffStreetParking in TwinInstance {event.twin_instance}")
                    app.logger.error(error)

if __name__ == "__main__":
    app.logger.info("Starting up server...")
    app.run(host='0.0.0.0', port=8080)
