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

    try:
        kcommand.handle_command(request=request, twin_interface='s4city-city-neighbourhood', command='updateairqualityindex', twin_graph=ktwin_graph, callback=handle_update_air_quality_index)
    except Exception as error:
        app.logger.error(f"Error to handle command updateairqualityindex in TwinInstance {event.twin_instance}")
        app.logger.error(error)

    # Return 204 - No-content
    return "", 204

def handle_update_air_quality_index(command_event: kcommand.KTwinCommandEvent, target_twin_instance: ktwingraph.TwinInstanceReference):
    app.logger.info(command_event)
    app.logger.info(target_twin_instance)

if __name__ == "__main__":
    app.logger.info("Starting up server...")
    app.run(host='0.0.0.0', port=8080)