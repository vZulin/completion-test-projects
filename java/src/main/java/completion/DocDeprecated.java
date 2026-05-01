package completion;

/**
 * Deprecated markers, annotation completion, and Javadoc tag completion.
 * Covers: TC-73 (deprecated), TC-77 (annotation), TC-81/TC-82 (Javadoc tags),
 * TC-124 (Javadoc FQN class reference), TC-125 (Javadoc throws class reference),
 * TC-126 (annotation completion with required value).
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

    /**
     * <caret> TC-124: In the @see tag below, delete 'java.util.' and the final 't'
     * so the token is 'Lis'. Accept List and verify the reference becomes fully qualified
     * without adding a Java import.
     *
     * @see java.util.List
     */
    void javadocFullyQualifiedReference() {
    }

    /**
     * <caret> TC-125: In the @throws tag below, delete 'java.io.' and the final 'on'
     * so the exception token is 'IOExcepti'. Accept IOException and verify the tag remains valid
     * without adding a Java import.
     *
     * @throws java.io.IOException when the read operation cannot complete
     */
    void javadocThrowsReference() throws java.io.IOException {
    }

    // <caret> TC-126: In the annotation below, delete 'gs("unchecked")' so the token is
    //   '@SuppressWarnin'. Accept SuppressWarnings and verify parentheses are inserted with
    //   the caret inside for the required value.
    @SuppressWarnings("unchecked")
    void suppressWarningsCompletion() {
    }

    @SuppressWarnings("deprecation")
    public static void main(String[] args) {
        // <caret> TC-73: Delete 'Method()' below so only 'old' remains,
        //   then invoke completion — oldMethod should appear with deprecated styling
        oldMethod();
    }
}
