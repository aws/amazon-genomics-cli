import connexion

from rest_api import encoder
from paste.translogger import TransLogger
import logging

def create_app(config=None):
    app = connexion.App(__name__, specification_dir="openapi/")
    app.app.json_encoder = encoder.JSONEncoder
    app.add_api("openapi.yaml",
                arguments={"title": "Workflow Execution Service"},
                pythonic_params=True)

    logger = logging.getLogger('waitress')
    logger.setLevel(logging.INFO)
    return TransLogger(app)
