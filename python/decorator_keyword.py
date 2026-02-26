"""
Decorator and keyword completion.

Covers: TC-79, TC-20 (Python variants).
Source: EX-PY-2, EX-PY-3.
"""
from dataclasses import dataclass


# --- TC-79: decorator completion ---
# <caret> TC-79: In decorator below, delete 'lass' so token becomes '@datac',
#   then invoke completion; expect @dataclass.
@dataclass
class SampleDataclass:
    value: int

# --- TC-57: auto-popup on '@' ---
# <caret> TC-57: Place caret on empty line before class below, type '@',
#   verify decorator/annotation completion popup appears.
class DecoratorAutoPopupTarget:
    pass


# --- TC-20: keyword completion ---
# <caret> TC-20: On a new top-level line, type 'cla' and invoke completion;
#   expect 'class' keyword.
class KeywordPlaceholder:
    pass
