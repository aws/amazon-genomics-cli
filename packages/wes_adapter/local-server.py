#!/usr/bin/env python3

from waitress import serve
from rest_api import create_app


def main():
    app = create_app()
    serve(app, host="0.0.0.0", port=80)


if __name__ == '__main__':
    main()
