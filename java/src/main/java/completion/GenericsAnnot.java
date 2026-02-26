package completion;

import java.util.List;
import java.util.Map;

/**
 * Generics type-parameter completion and annotation attribute completion.
 * Covers: TC-74 (generics), TC-75 (Map second type), TC-76 (auto-import in generics), TC-80.
 * Based on EX-JV-12.
 */
public class GenericsAnnot {

    @interface Ann {
        String value();
        boolean enabled();
    }

    // <caret> TC-80: Delete annotation args below, place caret inside @Ann(...);
    //   expect value=, enabled=
    @Ann(value = "demo", enabled = true)
    void annotated() {
    }

    // TC-74/TC-75/TC-76: generics type-parameter completion inside angle brackets
    void generics(
            // <caret> TC-74: Delete 'String' below inside List<>, invoke completion;
            //   expect String, Integer, User, etc.
            List<String> a,
            // <caret> TC-75: Delete 'Integer' below inside Map<String, ...>, invoke completion;
            //   expect type suggestions for second generic parameter
            Map<String, Integer> b,
            // <caret> TC-76: Delete 'java.time.' below and invoke completion on LocalDate inside generic;
            //   expect auto-import for java.time.LocalDate after accept
            Map<String, java.time.LocalDate> c
    ) {
        System.out.println(a.size() + b.size() + c.size());
    }
}
