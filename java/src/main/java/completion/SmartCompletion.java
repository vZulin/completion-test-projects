package completion;

import completion.model.User;

/**
 * Smart completion scenarios — return type and expected argument type.
 * Covers: TC-36 (smart return), TC-37 (expected type in args).
 * Based on EX-JV-6, EX-JV-7.
 */
public class SmartCompletion {

    static User build() {
        return new User();
    }

    static User createUser() {
        return // <caret> TC-36: smart completion for return — expect build(), new User()
                null;
    }

    static void consume(User u) {
    }

    public static void main(String[] args) {
        User u = new User();
        consume( // <caret> TC-37: expected type in args — expect u, build(), new User()
        );
    }
}
