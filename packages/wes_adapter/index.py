import serverless_wsgi
from rest_api import create_app

app = create_app()

# stops wsgi from sending Json as a Base64 string
serverless_wsgi.TEXT_MIME_TYPES.append("application/problem+json")

def handler(event, context):
    return serverless_wsgi.handle_request(app, event, context)
