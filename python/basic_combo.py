"""
Basic completion combo: member, smart, expected-type, and path completion.

Covers: TC-7, TC-26, TC-34, TC-37, TC-60 (Python variants).
Source: EX-PY-1.
"""
from dataclasses import dataclass
from pathlib import Path


@dataclass
class User:
    name: str
    age: int


def build_user(name: str, age: int) -> User:
    return User(name, age)


def consume_user(u: User) -> None:
    pass


user_name = "Ann"
user_age = 21
user = User(user_name, user_age)

user.  # <caret> TC-26: member completion after dot — expect name, age

u1: User =   # <caret> TC-34: smart-ish completion — expect build_user, User(...)

consume_user()  # <caret> TC-37: expected type completion — place caret inside parens, expect user

Path("")  # <caret> TC-60: path completion — place caret inside quotes, expect filesystem paths
