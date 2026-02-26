package completion;

/**
 * Performance test — completion with a large number of methods.
 * Covers: TC-90.
 * Based on EX-JV-17.
 */
public class PerfMany {

    void method000() {}
    void method001() {}
    void method002() {}
    void method003() {}
    void method004() {}
    void method005() {}
    void method006() {}
    void method007() {}
    void method008() {}
    void method009() {}
    void method010() {}
    void method011() {}
    void method012() {}
    void method013() {}
    void method014() {}
    void method015() {}
    void method016() {}
    void method017() {}
    void method018() {}
    void method019() {}
    void method020() {}
    void method021() {}
    void method022() {}
    void method023() {}
    void method024() {}
    void method025() {}
    void method026() {}
    void method027() {}
    void method028() {}
    void method029() {}
    void method030() {}
    void method031() {}
    void method032() {}
    void method033() {}
    void method034() {}
    void method035() {}
    void method036() {}
    void method037() {}
    void method038() {}
    void method039() {}
    void method040() {}
    void method041() {}
    void method042() {}
    void method043() {}
    void method044() {}
    void method045() {}
    void method046() {}
    void method047() {}
    void method048() {}
    void method049() {}

    public static void main(String[] args) {
        PerfMany m = new PerfMany();
        // <caret> TC-90: Delete 'hod000()' below so only 'm.met' remains,
        //   then invoke completion — popup should appear quickly with 50+ methods
        m.method000();
    }
}
