import os
import sys
import logging
import datetime
from dotenv import load_dotenv
from flask import Flask, request
import modules.ktwin.event as kevent
import modules.ktwin.eventstore as keventstore

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

    kevent.handle_event(request, 'ngsi-ld-city-device', handle_device_event)

    # Return 204 - No-content
    return "", 204

# Handle event device
# Logic: 
# - In case the battery level is below 15% reduce the frequency to 1 time every 60 min
# - In case the battery level is above 15% keep the frequency to 1 time every 15 min
def handle_device_event(event: kevent.KTwinEvent):
    HIGH_FREQUENCY = 15 # 15min
    LOW_FREQUENCY = 60 # 60min
    BATTERY_THRESHOLD = 15 # percentage of battery available

    device_event_data = event.cloud_event.data
    device_event_data["dateObserved"] = datetime.datetime.now().isoformat()

    if "batteryLevel" in device_event_data:
        if device_event_data["batteryLevel"] < BATTERY_THRESHOLD:
            # Propagate event to real device to measure in low frequency
            device_event_data["measurementFrequency"] = LOW_FREQUENCY
            kevent.send_to_real_twin(twin_interface=event.twin_interface, twin_instance=event.twin_instance, data=device_event_data)
        elif device_event_data["batteryLevel"] > BATTERY_THRESHOLD and "measurementFrequency" in device_event_data and device_event_data["measurementFrequency"] == LOW_FREQUENCY:
            # Propagate event to real device to measure in high frequency
            device_event_data["measurementFrequency"] = HIGH_FREQUENCY
            kevent.send_to_real_twin(twin_interface=event.twin_interface, twin_instance=event.twin_instance, data=device_event_data)

        event.cloud_event.data = device_event_data
        keventstore.update_twin_event(event)

    else:
        app.logger.info("Battery level was not provided")

if __name__ == "__main__":
    app.logger.info("Starting up server...")
    app.run(host='0.0.0.0', port=8080)