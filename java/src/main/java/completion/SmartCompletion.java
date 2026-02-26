package completion;

import completion.model.User;

/**
 * Smart completion scenarios — return type and expected argument type.
 * Covers: TC-30 (smart return), TC-31 (expected type in args).
 * Based on EX-JV-6, EX-JV-7.
 */
public class SmartCompletion {

    static User build() {
        return new User();
    }

    static User createUser() {
        // <caret> TC-30: Delete 'build()' below after 'return', invoke smart completion;
        //   expect build(), new User(...)
        return build();
    }

    static void consume(User u) {
    }

    public static void main(String[] args) {
        User u = new User();
        // <caret> TC-31: Delete 'u' below inside consume(...), invoke completion;
        //   expect u, build(), new User(...)
        consume(u);
    }
}
