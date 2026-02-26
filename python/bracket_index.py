"""
Bracket-index completion scenarios.

Covers: TC-101 (Python variant).
"""

users = [
    {"name": "Ann", "age": 21},
    {"name": "Bob", "age": 22},
]
idx = 0

# <caret> TC-101-A: Delete 'idx' below inside users[...], invoke completion;
#   expect numeric index variables from scope (e.g., idx)
tc106a = users[idx]

# <caret> TC-101-B: Delete '"name"' below inside users[idx][...], type 'na', invoke completion;
#   completion may suggest known dict keys in current context.
tc106b = users[idx]["name"]

_ = (tc106a, tc106b)
