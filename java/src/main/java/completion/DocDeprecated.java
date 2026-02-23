package completion;

/**
 * Deprecated markers, annotation completion, and Javadoc tag completion.
 * Covers: TC-66 (deprecated), TC-70 (annotation), TC-74/TC-75 (Javadoc tags).
 * Based on EX-JV-11, EX-JV-2, EX-JV-13.
 */
public class DocDeprecated {

    // --- TC-66: deprecated marker in completion list ---

    /** @deprecated Use newMethod instead. */
    @Deprecated
    static void oldMethod() {
    }

    static void newMethod() {
    }

    // --- TC-70: annotation completion ---
    // Place caret after @Depr and invoke completion — expect @Deprecated
    @Depr // <caret> TC-70: annotation completion — expect @Deprecated
    static void annotatedMethod() {
    }

    // --- TC-74, TC-75: Javadoc tag completion ---

    /**
     * @ // <caret> TC-74: Javadoc tag completion — expect @param, @return, @throws, etc.
     * TC-75: verify @param a appears for the parameter below
     */
    int compute(int a) {
        return a;
    }

    public static void main(String[] args) {
        old // <caret> TC-66: deprecated style — oldMethod should appear with strikethrough
        ;
    }
}
