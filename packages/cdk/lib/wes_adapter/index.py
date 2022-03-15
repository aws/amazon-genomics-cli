import serverless_wsgi
from rest_api import create_app

app = create_app()


def handler(event, context):
    return serverless_wsgi.handle_request(app, event, context)
