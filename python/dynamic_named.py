"""
Dynamic/unknown fallback and named-argument completion.

Covers: TC-32, TC-66 (Python variants).
Source: EX-PY-4, EX-PY-5.
"""


# --- TC-32: dynamic/unknown fallback ---
# Place caret after 'dyn.' and invoke completion; should not crash,
# may show limited members from object.
dyn = object()
# <caret> TC-32: Delete '__class__.__name__' below so caret is after 'dyn.',
#   then invoke completion; should not crash, may show limited members.
dynamic_class_name = dyn.__class__.__name__


# --- TC-66: named argument completion ---
# Place caret inside f() parens and invoke completion; expect name=, age=.
def greet(name: str, age: int) -> None:
    pass


# <caret> TC-66: Delete named args below, place caret inside greet(...),
#   invoke completion; expect name=, age=.
greet(name="Ann", age=21)

_ = dynamic_class_name
