package completion;

import completion.model.User;

import java.util.List;

/**
 * Smart completion scenarios — return type and expected argument type.
 * Covers: TC-30 (smart return), TC-31 (expected type in args),
 * TC-133 (constructor completion from return expected type).
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

    static List<String> createNames() {
        // <caret> TC-133: Remove 'java.util.' from the constructor call below and delete 't<>();'
        //   so the return expression ends with 'new ArrayLis'. Accept ArrayList and verify constructor
        //   insertion plus import java.util.ArrayList.
        return new java.util.ArrayList<>();
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
