import werkzeug


class InvalidRequestError(werkzeug.exceptions.BadRequest):
    """
    This exception is raised when we receive an invalid request.
    """
    pass


class InternalServerError(werkzeug.exceptions.InternalServerError):
    """
    This exception is raised internal server error occurred.
    """
    pass