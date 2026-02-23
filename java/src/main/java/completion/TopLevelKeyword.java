package completion;

/**
 * Top-level keyword completion test.
 * Covers: TC-21 (keyword completion at top level).
 * Based on EX-JV-4.
 */
// cla // <caret> TC-21: keyword class — expect 'class' keyword suggestion

public class TopLevelKeyword {
    // Placeholder class so the file compiles.
    // To test TC-21, place caret after "cla" in a blank line at the top level.
}
