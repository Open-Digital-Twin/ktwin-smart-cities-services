import os
import datetime
import sys
import logging
from enum import Enum
from dotenv import load_dotenv
from flask import Flask, request
import modules.ktwin.event as kevent
import modules.ktwin.eventstore as keventstore
import modules.ktwin.twingraph as ktwingraph
import modules.ktwin.command as kcommand

if os.getenv("ENV") == "local":
    load_dotenv('local.env')

app = Flask(__name__)

handler = logging.StreamHandler(sys.stdout)
handler.setFormatter(logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s'))
app.logger.addHandler(handler)
app.logger.setLevel(logging.INFO)

ktwin_graph = ktwingraph.load_twin_graph()

class AQICategory(Enum):
    GOOD="GOOD"
    MODERATE="MODERATE"
    UNHEALTHY_FOR_SENSITIVE_GROUPS="UNHEALTHY_FOR_SENSITIVE_GROUPS"
    UNHEALTHY="UNHEALTHY"
    VERY_UNHEALTHY="VERY_UNHEALTHY"
    HAZARDOUS="HAZARDOUS"

def convert_aqi_category_int(aqi_category_value: str):
    if aqi_category_value == AQICategory.GOOD.value:
        return 1
    elif aqi_category_value == AQICategory.MODERATE.value:
        return 2
    elif aqi_category_value == AQICategory.UNHEALTHY_FOR_SENSITIVE_GROUPS.value:
        return 3
    elif aqi_category_value == AQICategory.UNHEALTHY.value:
        return 4
    elif aqi_category_value == AQICategory.VERY_UNHEALTHY.value:
        return 5
    elif aqi_category_value == AQICategory.HAZARDOUS.value:
        return 6

@app.route("/", methods=["POST"])
def home():
    event = kevent.handle_request(request)

    app.logger.info(
        f"Event TwinInstance: {event.twin_instance} - Event TwinInterface: {event.twin_interface}"
    )

    try:
        kcommand.handle_command(request=request, twin_interface='s4city-city-neighbourhood', command='updateairqualityindex', twin_graph=ktwin_graph, callback=handle_update_air_quality_index)
    except Exception as error:
        app.logger.error(f"Error to handle command updateairqualityindex in TwinInstance {event.twin_instance}")
        app.logger.error(error)

    # Return 204 - No-content
    return "", 204

# Save the worst level in the last 24h
def handle_update_air_quality_index(command_event: kcommand.KTwinCommandEvent, target_twin_instance: ktwingraph.TwinInstanceReference):

    latest_neighborhood_event = keventstore.get_latest_twin_event(twin_interface=target_twin_instance.interface, twin_instance=target_twin_instance.instance)
    latest_neighborhood_data = None

    if latest_neighborhood_event is None:
        event_data = dict()
        event_data["aqiLevel"] = AQICategory.GOOD.value
        event_data["dateObserved"] = datetime.datetime.now().isoformat()
        latest_neighborhood_event = kevent.KTwinEvent()
        latest_neighborhood_event.set_event(twin_interface=command_event.twin_interface, twin_instance=target_twin_instance.instance, data=event_data)

    latest_neighborhood_data = latest_neighborhood_event.cloud_event.data

    datetime_now = datetime.datetime.now()

    if "aqiLevel" not in latest_neighborhood_data:
        app.logger.info(f"Event {command_event.cloud_event} has no aqiLevel attribute value")
    else:
        command_data = command_event.cloud_event.data

        if "aqiLevel" not in command_data:
            app.logger.info(f"Command {command_event.cloud_event} has no aqiLevel attribute value")
        else:
            latest_air_index = convert_aqi_category_int(latest_neighborhood_data["aqiLevel"])
            current_air_index = convert_aqi_category_int(latest_neighborhood_data["aqiLevel"])

            if current_air_index >= latest_air_index:
                latest_neighborhood_data["aqiLevel"] = command_data["aqiLevel"]
            
            has_expired = True
            if "dateObserved" in latest_neighborhood_data:
                has_expired = has_time_expired(datetime_now=datetime_now, date_last_switching=datetime.datetime.fromisoformat(latest_neighborhood_data["dateObserved"]))

            if has_expired:
                latest_neighborhood_data["aqiLevel"] = command_data["aqiLevel"]

            latest_neighborhood_event.cloud_event.data = latest_neighborhood_data

            keventstore.update_twin_event(latest_neighborhood_event)

    app.logger.info(command_event)
    app.logger.info(target_twin_instance)

# In case of 24h of no change in the state, we update the current state
def has_time_expired(datetime_now: datetime, date_last_switching: datetime):
    if (datetime_now - date_last_switching).seconds > 24*60*60:
        return True
    return 

if __name__ == "__main__":
    app.logger.info("Starting up server...")
    app.run(host='0.0.0.0', port=8080)