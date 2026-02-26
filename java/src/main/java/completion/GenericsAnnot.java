package completion;

import java.util.List;
import java.util.Map;

/**
 * Generics type-parameter completion and annotation attribute completion.
 * Covers: TC-75 (generics), TC-76 (Map generics), TC-80 (annotation params).
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

    // TC-75/TC-76: generics type-parameter completion inside angle brackets
    void generics(
            // <caret> TC-75: Delete 'String' below inside List<>, invoke completion;
            //   expect String, Integer, User, etc.
            List<String> a,
            // <caret> TC-76: Delete 'Integer' below inside Map<String, ...>, invoke completion;
            //   expect type suggestions for second generic parameter
            Map<String, Integer> b
    ) {
        System.out.println(a.size() + b.size());
    }
}
