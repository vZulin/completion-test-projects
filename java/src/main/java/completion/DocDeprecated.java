package completion;

/**
 * Deprecated markers, annotation completion, and Javadoc tag completion.
 * Covers: TC-73 (deprecated), TC-77 (annotation), TC-81/TC-82 (Javadoc tags).
 * Based on EX-JV-11, EX-JV-2, EX-JV-13.
 */
public class DocDeprecated {

    // --- TC-73: deprecated marker in completion list ---

    /** @deprecated Use newMethod instead. */
    @Deprecated
    static void oldMethod() {
    }

    static void newMethod() {
    }

    // --- TC-57, TC-77: annotation completion ---
    // <caret> TC-57/TC-77: In annotation below, delete 'ecated' so token becomes '@Depr',
    //   then invoke completion — expect @Deprecated
    @Deprecated
    static void annotatedMethod() {
    }

    // --- TC-81, TC-82: Javadoc tag completion ---

    /**
     * <caret> TC-81/TC-82: In tags below, delete 'param a' so only '@' remains,
     * then invoke completion — expect @param, @return, @throws, etc.
     * Verify @param a appears for the method parameter below.
     *
     * @param a input value
     * @return input value
     */
    int compute(int a) {
        return a;
    }

    @SuppressWarnings("deprecation")
    public static void main(String[] args) {
        // <caret> TC-73: Delete 'Method()' below so only 'old' remains,
        //   then invoke completion — oldMethod should appear with deprecated styling
        oldMethod();
    }
}
