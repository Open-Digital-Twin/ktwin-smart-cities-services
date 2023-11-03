import os
import math
import sys
import logging
from dotenv import load_dotenv
from flask import Flask, request
import modules.ktwin.event as kevent
import modules.ktwin.eventstore as keventstore
import modules.ktwin.twingraph as ktwingraph
import modules.ktwin.command as kcommand
import modules.airqualityindex as aqi

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

    kevent.handle_event(request, 'ngsi-ld-city-airqualityobserved', handle_air_quality_observed_event)
    kevent.handle_event(request, 'ngsi-ld-city-weatherobserved', handle_weather_observed_event)
    
    # Return 204 - No-content
    return "", 204

def handle_air_quality_observed_event(event: kevent.KTwinEvent):
    air_quality_observed = event.cloud_event.data
    # air_quality_observed["CO2AqiLevel"] = air_quality_level(air_quality_observed["CO2Density"])
    # air_quality_observed["NOAqiLevel"] = air_quality_level(air_quality_observed["NODensity"])
    # air_quality_observed["C6H6AqiLevel"] = air_quality_level(air_quality_observed["C6H6Density"])
    # air_quality_observed["CDAqiLevel"] = air_quality_level(air_quality_observed["CDDensity"])
    # air_quality_observed["PBAqiLevel"] = air_quality_level(air_quality_observed["PBDensity"])
    # air_quality_observed["SH2AqiLevel"] = air_quality_level(air_quality_observed["SH2Density"])

    air_quality_observed["COAqiLevel"] = aqi.COAirQualityIndex(concentration=air_quality_observed["CODensity"]).get_air_quality_category()
    air_quality_observed["PM10AqiLevel"] = aqi.PM10AirQualityIndex(concentration=air_quality_observed["PM10Density"]).get_air_quality_category()
    air_quality_observed["PM25AqiLevel"] = aqi.PM25AirQualityIndex(concentration=air_quality_observed["PM25Density"]).get_air_quality_category()
    air_quality_observed["SO2AqiLevel"] = aqi.PM25AirQualityIndex(concentration=air_quality_observed["SO2Density"]).get_air_quality_category()
    air_quality_observed["O3AqiLevel"] = aqi.O3AirQualityIndex(concentration=air_quality_observed["O3Density"]).get_air_quality_category()

    event.cloud_event.data = air_quality_observed

    print("event")
    print(event)

    print("event.cloud_event")
    print(event.cloud_event)

    print("event.cloud_event.data")
    print(event.cloud_event.data)

    keventstore.update_twin_event(event)

    all_levels = list()
    all_levels.append(air_quality_observed["COAqiLevel"])
    all_levels.append(air_quality_observed["PM10AqiLevel"])
    all_levels.append(air_quality_observed["PM25AqiLevel"])
    all_levels.append(air_quality_observed["SO2AqiLevel"])
    all_levels.append(air_quality_observed["O3AqiLevel"])

    # The largest or "dominant" AQI value is reported for the location and propagated to the neighborhood.

    payload = dict()
    if aqi.AQICategory.HAZARDOUS in all_levels:
        payload["aqiLevel"] = aqi.AQICategory.HAZARDOUS
    elif aqi.AQICategory.VERY_UNHEALTHY in all_levels:
        payload["aqiLevel"] = aqi.AQICategory.VERY_UNHEALTHY
    elif aqi.AQICategory.UNHEALTHY in all_levels:
        payload["aqiLevel"] = aqi.AQICategory.UNHEALTHY
    elif aqi.AQICategory.UNHEALTHY_FOR_SENSITIVE_GROUPS in all_levels:
        payload["aqiLevel"] = aqi.AQICategory.UNHEALTHY_FOR_SENSITIVE_GROUPS
    elif aqi.AQICategory.MODERATE in all_levels:
        payload["aqiLevel"] = aqi.AQICategory.MODERATE
    else:
        payload["aqiLevel"] = aqi.AQICategory.GOOD

    try:
        kcommand.execute_command(command_payload=payload, command="updateairqualityindex", relationship_name="neighborhood", twin_instance=event.cloud_event["source"], twin_graph=ktwin_graph)
    except Exception as error:
        app.logger.error(f"Error to execute command updateairqualityindex in relation neighborhood in TwinInstance {event.twin_instance}")
        app.logger.error(error)


def handle_weather_observed_event(event: kevent.KTwinEvent):
    app.logger.info(f"Processing {event.twin_instance} event")

    latest_event = keventstore.get_latest_twin_event(event.twin_interface, event.twin_instance)
    if latest_event is None:
        latest_event = event

    weather_observed = event.cloud_event.data
    weather_observed["pressureTendency"] = calculate_pressure_tendency(latest_event, event)
    weather_observed["FeelsLikeTemperature"] = calculate_feel_like_temperature(weather_observed["temperature"], weather_observed["WindSpeed"])
    weather_observed["dewpoint"] = calculate_dewpoint(weather_observed["temperature"], weather_observed["relativeHumidity"])
    event.cloud_event.data = weather_observed

    keventstore.update_twin_event(event)

def calculate_pressure_tendency(latest_event: kevent.KTwinEvent, current_event: kevent.KTwinEvent):
    latest_cloud_event = latest_event.cloud_event.data
    current_cloud_event = current_event.cloud_event.data

    if latest_cloud_event["atmosphericPressure"] is not None and current_cloud_event["atmosphericPressure"] is not None:
        difference = current_cloud_event["atmosphericPressure"] - latest_cloud_event["atmosphericPressure"]
        if abs(difference) < 0.1:
            return "steady"
        if difference < 0:
            return "falling"
        return "raising"
    else:
        "steady"

def calculate_feel_like_temperature(temperature: float, wind_speed: float):
    return 33 + (10 * math.sqrt(wind_speed) + 10.45 - wind_speed) * (temperature - 33)/22

def calculate_dewpoint(temperature: float, relative_humidity: float):
    return temperature - ((100-relative_humidity/5))

if __name__ == "__main__":
    app.logger.info("Starting up server...")
    app.run(host='0.0.0.0', port=8080)