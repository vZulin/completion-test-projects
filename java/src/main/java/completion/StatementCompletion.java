package completion;

/**
 * Statement-level completion scenarios: if-else, switch-case-default.
 *
 * Covers: TC-92 (if-else completion), TC-93 (switch/case/default).
 */
public class StatementCompletion {

    enum Status { ACTIVE, INACTIVE, BANNED }

    /** TC-92: if-else statement completion. */
    static String checkAge(int age) {
        // <caret> TC-92-A: type 'if' here, invoke completion — expect 'if' statement
        if (age > 18) {
            return "adult";
        }
        // <caret> TC-92-B: place caret here after '}', type 'el' — expect 'else' / 'else if'
        else {
            return "minor";
        }
    }

    /** TC-93: switch/case/default completion. */
    static String describeStatus(Status status) {
        // <caret> TC-93-A: type 'sw' inside method body, invoke completion — expect 'switch' keyword
        // <caret> TC-93-B: place caret inside switch block after last 'break;' — expect remaining cases / 'default'
        switch (status) {
            case ACTIVE:
                return "active user";
            case INACTIVE:
                return "inactive user";
            case BANNED:
                return "banned user";
            default:
                return "unknown";
        }
    }

    /** TC-93 continued: switch with int — case and default completion. */
    static String classifyAge(int age) {
        // <caret> TC-93-C: inside switch block, invoke completion — expect 'case', 'default'
        switch (age / 10) {
            case 0:
                return "child";
            case 1:
                return "teen";
            default:
                return "adult";
        }
    }
}
