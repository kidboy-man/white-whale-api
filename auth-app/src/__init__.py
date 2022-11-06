from flask.json import jsonify
from flask import Flask
import os
from http import HTTPStatus
from src.auth import auth
from src.database import db
from flask_jwt_extended import JWTManager
from flasgger import Swagger, swag_from
from src.config.swagger import template, swagger_config


def create_app(test_config=None):

    app = Flask(__name__, instance_relative_config=True)

    if test_config is None:
        app.config.from_mapping(
            SECRET_KEY=os.environ.get("SECRET_KEY"),
            SQLALCHEMY_DATABASE_URI=os.environ.get("SQLALCHEMY_DB_URI"),
            SQLALCHEMY_TRACK_MODIFICATIONS=False,
            JWT_SECRET_KEY=os.environ.get("JWT_SECRET_KEY"),
            SWAGGER={"title": "Auth-App API", "uiversion": 3},
        )
    else:
        app.config.from_mapping(test_config)

    db.app = app
    db.init_app(app)
    db.create_all()

    JWTManager(app)
    app.register_blueprint(auth)

    Swagger(app, config=swagger_config, template=template)

    @app.errorhandler(HTTPStatus.NOT_FOUND)
    def handle_404(e):
        return jsonify({"error": "Not found"}), HTTPStatus.NOT_FOUND

    @app.errorhandler(HTTPStatus.INTERNAL_SERVER_ERROR)
    def handle_500(e):
        return (
            jsonify({"error": "Something went wrong"}),
            HTTPStatus.INTERNAL_SERVER_ERROR,
        )

    return app
