"""
Basic completion combo: member, smart, expected-type, and path completion.

Covers: TC-1, TC-2, TC-7, TC-8, TC-11, TC-19, TC-22, TC-28, TC-31, TC-47 (Python variants).
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

# <caret> TC-1/TC-7/TC-22: Delete 'name' below so caret is after 'user.',
#   then invoke completion — expect name, age and other members.
tc26 = user.name

# <caret> TC-2/TC-28: Delete expression below after '=', invoke smart-ish completion;
#   expect build_user(...), User(...)
u1: User = build_user(user_name, user_age)

# <caret> TC-8: Delete args inside build_user(...), place caret after '(' or after ',',
#   and invoke completion for argument suggestions.
tc_arg = build_user(user_name, user_age)

# <caret> TC-31: Delete 'user' below inside consume_user(...), invoke completion;
#   expect variable 'user' of type User
consume_user(user)

# <caret> TC-11/TC-47: Place caret inside quotes below and invoke completion;
#   expect filesystem path suggestions
tc60 = Path("")


def keyword_return_completion(value: int) -> int:
    # <caret> TC-19: Inside function body type 'ret' and invoke completion; expect 'return'.
    return value


_ = (tc26, u1, tc_arg, tc60, keyword_return_completion(1))
