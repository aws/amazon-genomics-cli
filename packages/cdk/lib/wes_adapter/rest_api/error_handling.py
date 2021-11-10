from rest_api.models import ErrorResponse


def handle_invalid_request(error):
    return __generate_error_response(error, 400)


def handle_internal_error(error):
    return __generate_error_response(error, 500)


def __generate_error_response(error, status_code):
    return ErrorResponse(msg=str(error), status_code=status_code)
