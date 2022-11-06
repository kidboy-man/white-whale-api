from http import HTTPStatus
from flask import Blueprint, request, jsonify
from werkzeug.security import check_password_hash, generate_password_hash
from flask_jwt_extended import (
    jwt_required,
    create_access_token,
    create_refresh_token,
    get_jwt_identity,
    get_jwt,
)
from flasgger import swag_from
from src.database import User, db
from src.utils import generate_random_string

auth = Blueprint("auth", __name__, url_prefix="/v1/auth")


@auth.post("/register")
@swag_from("./docs/auth/register.yaml")
def register():
    phone = request.json["phone"]
    role = request.json["role"]

    if len(phone) < 3:
        return jsonify({"error": "phone number is too short"}), HTTPStatus.BAD_REQUEST

    if not phone.isnumeric() or " " in phone:
        return (
            jsonify({"error": "phone should be numeric, also no spaces"}),
            HTTPStatus.BAD_REQUEST,
        )

    if phone[:2] != "08":
        return (
            jsonify({"error": "invalid phone number format, must be in 08XXXX"}),
            HTTPStatus.BAD_REQUEST,
        )

    if role not in ["user", "admin"]:
        return (
            jsonify({"error": "role must be one of admin or user"}),
            HTTPStatus.BAD_REQUEST,
        )

    if User.query.filter_by(phone=phone).first() is not None:
        return jsonify({"error": "phone is taken"}), HTTPStatus.CONFLICT

    generated_password = generate_random_string(4)
    pwd_hash = generate_password_hash(generated_password)
    user = User(phone=phone, hashed_password=pwd_hash, role=role)
    db.session.add(user)
    db.session.commit()

    return (
        jsonify(
            {
                "message": "User created",
                "user": {"phone": phone, "password": generated_password},
            }
        ),
        HTTPStatus.CREATED,
    )


@auth.post("/login")
@swag_from("./docs/auth/login.yaml")
def login():
    phone = request.json.get("phone", "")
    password = request.json.get("password", "")

    user = User.query.filter_by(phone=phone).first()

    if user:
        is_pass_correct = check_password_hash(user.hashed_password, password)

        if is_pass_correct:
            refresh = create_refresh_token(identity=user.id)
            additional_claims = {"role": user.role}
            access = create_access_token(
                identity=f"{user.id}", additional_claims=additional_claims
            )

            return (
                jsonify(
                    {
                        "user": {
                            "refresh": refresh,
                            "access": access,
                            "phone": user.phone,
                            "role": user.role,
                        }
                    }
                ),
                HTTPStatus.OK,
            )

    return jsonify({"error": "Wrong credentials"}), HTTPStatus.UNAUTHORIZED


@auth.get("/me")
@jwt_required()
def me():
    claims = get_jwt()
    user_id = get_jwt_identity()
    user = User.query.filter_by(id=user_id).first()
    return (
        jsonify({"phone": user.phone, "role": user.role, "claims": claims}),
        HTTPStatus.OK,
    )


@auth.get("/token/refresh")
@jwt_required(refresh=True)
def refresh_users_token():
    identity = get_jwt_identity()
    access = create_access_token(identity=identity)

    return jsonify({"access": access}), HTTPStatus.OK
