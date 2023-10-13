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

    handle_event(request, 'ngsi-ld-city-streetlight', handle_streetlight_event)

    # Return 204 - No-content
    return "", 204

def handle_streetlight_event(event: KTwinEvent):
    current_streetlight_event = event.cloud_event.data
    
    if "powerState" not in current_streetlight_event:
        app.logger.info(f"Event {event.cloud_event} has no powerState attribute value")
    else:
        current_power_state_value = current_streetlight_event["powerState"]
        latest_streetlight_event = get_latest_twin_event(event.twin_interface, event.twin_instance)
        datetime_now = datetime.datetime.now()

        if latest_streetlight_event is not None:
            latest_power_state_value = latest_streetlight_event.cloud_event.data["powerState"]
            defect = False
            if current_power_state_value == latest_power_state_value:
                if current_power_state_value == "on":
                    date_last_switching_on = latest_streetlight_event.cloud_event.data["dateLastSwitchingOn"]
                    defect = is_with_defect(datetime_now=datetime_now, date_last_switching=date_last_switching_on)
                elif current_power_state_value == "off":
                    date_last_switching_off = latest_streetlight_event.cloud_event.data["dateLastSwitchingOff"]
                    defect = is_with_defect(datetime_now=datetime_now, date_last_switching=date_last_switching_off)

                if defect:
                    event.cloud_event.data["powerState"] = "defectiveLamp"
                else:
                    if current_power_state_value == "on":
                        event.cloud_event.data["dateLastSwitchingOn"] = datetime_now.isoformat()
                    elif current_power_state_value == "off":
                        event.cloud_event.data["dateLastSwitchingOff"] = datetime_now.isoformat()
                    event.cloud_event.data["powerState"] = "ok"
        else:
            if current_power_state_value == "on":
                event.cloud_event.data["dateLastSwitchingOn"] = datetime_now.isoformat()
            elif current_power_state_value == "off":
                event.cloud_event.data["dateLastSwitchingOff"] = datetime_now.isoformat()
            event.cloud_event.data["powerState"] = "ok"

        update_twin_event(event)

# In case of 48h of no change in the state, we consider that lamp with a defect
def is_with_defect(datetime_now: datetime, date_last_switching: datetime):
    if (datetime_now - date_last_switching).seconds > 2*24*60*60:
        return True
    return 

if __name__ == "__main__":
    app.logger.info("Starting up server...")
    app.run(host='0.0.0.0', port=8080)