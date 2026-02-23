"""
Docstring completion inside function documentation.

Covers: TC-77 (Python variant).
Source: EX-PY-6.
"""


# --- TC-77: docstring completion ---
# Place caret on the blank line inside the docstring and invoke completion;
# expect parameter stubs (:param a:), return type, etc.
def compute(a: int) -> int:
    """
      # <caret> TC-77: docstring completion — expect :param a:, :return:, etc.
    """
    return a
