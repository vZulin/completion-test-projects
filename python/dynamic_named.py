"""
Dynamic/unknown fallback and named-argument completion.

Covers: TC-38, TC-45 (Python variants).
Source: EX-PY-4, EX-PY-5.
"""


# --- TC-38: dynamic/unknown fallback ---
# Place caret after 'dyn.' and invoke completion; should not crash,
# may show limited members from object.
dyn = object()
dyn.  # <caret> TC-38: dynamic fallback — should not crash, may show limited members


# --- TC-45: named argument completion ---
# Place caret inside f() parens and invoke completion; expect name=, age=.
def greet(name: str, age: int) -> None:
    pass


greet()  # <caret> TC-45: named args — place caret inside parens, expect name=, age=
