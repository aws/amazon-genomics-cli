from functools import wraps
import logging

import inspect


def logged(
    _func=None,
    *,
    level=logging.INFO,
    logger=None,
    log_input=True,
    log_output=True,
    log_exceptions=False,
):
    def logging_decorator(func):
        function_name = func.__name__
        arg_names = inspect.signature(func).parameters

        @wraps(func)
        def wrapper(*args, **kwargs):
            all_args = {**dict(zip(arg_names, args)), **kwargs}
            # Don't log 'self' and also use the objects logger if it has one.
            self = all_args.pop("self", None)
            if logger:
                log = logger
            elif self and self.logger:
                log = self.logger
            else:
                log = logging.getLogger(func.__module__)

            if log_input:
                log.log(level, f"{function_name} called with: {all_args}")
            try:
                result = func(*args, **kwargs)
                if log_output:
                    log.log(level, f"{function_name} returned: {result}")
                return result
            except Exception as e:
                if log_exceptions:
                    log.error(level, f"Caught exception calling {function_name}", e)
                raise e

        return wrapper

    if _func is None:
        return logging_decorator
    else:
        return logging_decorator(_func)
