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

TWIN_INTERFACE_PARKING = 'ngsi-ld-city-offstreetparking'
TWIN_INTERFACE_PARKING_SPOT = 'ngsi-ld-city-offstreetparkingspot'
app.logger.info('Loading Twin Interface {0}, {1}'.format(TWIN_INTERFACE_PARKING, TWIN_INTERFACE_PARKING_SPOT))
ktwin_graph = ktwingraph.load_twin_graph_by_instance(twin_instances=[TWIN_INTERFACE_PARKING, TWIN_INTERFACE_PARKING_SPOT])

@app.route("/", methods=["POST"])
def home():
    event = kevent.handle_request(request)

    app.logger.info(
        f"Event TwinInstance: {event.twin_instance} - Event TwinInterface: {event.twin_interface}"
    )
    
    try:
        kcommand.handle_command(request=request, twin_interface=TWIN_INTERFACE_PARKING, command='updateVehicleCount', twin_graph=ktwin_graph, callback=handle_update_vehicle_count_command)
    except Exception as error:
        app.logger.error(f"Error to handle command updateVehicleCount in TwinInstance {event.twin_instance}")
        app.logger.error(error)

    # Return 204 - No-content
    return "", 204

def handle_update_vehicle_count_command(command_event: kcommand.KTwinCommandEvent, target_twin_instance: ktwingraph.TwinInstanceReference):
    command_offstreetparking_data = command_event.cloud_event.data

    latest_offstreetparking_event = keventstore.get_latest_twin_event(twin_interface=target_twin_instance.interface, twin_instance=target_twin_instance.instance)
    latest_offstreetparking_data = None

    if latest_offstreetparking_event is None:
        event_data = dict()
        event_data["occupiedSpotNumber"] = 0
        latest_offstreetparking_event = kevent.KTwinEvent()
        latest_offstreetparking_event.set_event(twin_interface=command_event.twin_interface, twin_instance=target_twin_instance.instance, data=event_data)

    latest_offstreetparking_data = latest_offstreetparking_event.cloud_event.data

    if "vehicleEntranceCount" not in command_offstreetparking_data:
        app.logger.info(f"Event {command_event.cloud_event} has no vehicleEntranceCount attribute value")
    else:
        latest_offstreetparking_data["occupiedSpotNumber"] = latest_offstreetparking_data["occupiedSpotNumber"] + command_offstreetparking_data["vehicleEntranceCount"]
 
    if "vehicleExitCount" not in command_offstreetparking_data:
        app.logger.info(f"Event {command_event.cloud_event} has no vehicleExitCount attribute value")
    else:
        if "vehicleExitCount" in latest_offstreetparking_data:
            if latest_offstreetparking_data["occupiedSpotNumber"] <= 0:
                latest_offstreetparking_data["occupiedSpotNumber"] = latest_offstreetparking_data["occupiedSpotNumber"] - command_offstreetparking_data["vehicleExitCount"]
            else:
                app.logger.info(f"The number of occupied spot number cannot be negative")

    latest_offstreetparking_event.cloud_event.data = latest_offstreetparking_data

    keventstore.update_twin_event(latest_offstreetparking_event)

if __name__ == "__main__":
    app.logger.info("Starting up server...")
    app.run(host='0.0.0.0', port=8080)
