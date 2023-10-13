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

    handle_event(request, 'ngsi-ld-city-ev-charging-station', handle_ev_charging_station_event)

    # Return 204 - No-content
    return "", 204

def handle_ev_charging_station_event(event: KTwinEvent):
    current_ev_charing_event_event = event.cloud_event.data
    

if __name__ == "__main__":
    app.logger.info("Starting up server...")
    app.run(host='0.0.0.0', port=8080)