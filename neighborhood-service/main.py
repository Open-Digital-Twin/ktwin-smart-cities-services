import os
import math
import sys
import logging
from dotenv import load_dotenv
from flask import Flask, request
import modules.ktwin.event as kevent
import modules.ktwin.eventstore as keventstore
import modules.ktwin.command as kcommand

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

    kevent.handle_event(request, 'air-quality-observed', handle_neighborhood)

    # Return 204 - No-content
    return "", 204

def handle_neighborhood(event: kevent.KTwinEvent):
    air_quality_observed = event.cloud_event.data

if __name__ == "__main__":
    app.logger.info("Starting up server...")
    app.run(host='0.0.0.0', port=8080)