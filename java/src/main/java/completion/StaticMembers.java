package completion;

/**
 * Static members completion — only static members should appear.
 * Covers: TC-31.
 * Based on EX-JV-5.
 */
public class StaticMembers {

    public static void main(String[] args) {
        Math. // <caret> TC-31: static members only — expect abs, max, PI; no instance methods
        ;
    }
}
