"""
Decorator and keyword completion.

Covers: TC-72, TC-21 (Python variants).
Source: EX-PY-2, EX-PY-3.
"""
from dataclasses import dataclass


# --- TC-72: decorator completion ---
# Place caret after '@datac' and invoke completion; expect @dataclass.
@datac  # <caret> TC-72: decorator completion
class SampleDataclass:
    value: int


# --- TC-21: keyword completion ---
# Place caret after 'cla' and invoke completion; expect 'class' keyword.
cla  # <caret> TC-21: keyword completion — expect 'class'
