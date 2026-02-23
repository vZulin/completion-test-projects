package completion;

import java.util.List;
import java.util.Map;

/**
 * Generics type-parameter completion and annotation attribute completion.
 * Covers: TC-68 (generics), TC-69 (Map generics), TC-73 (annotation params).
 * Based on EX-JV-12.
 */
public class GenericsAnnot {

    @interface Ann {
        String value();
        boolean enabled();
    }

    // TC-73: annotation attribute completion — expect value=, enabled=
    @Ann( // <caret> TC-73: annotation params — expect value=, enabled=
    )
    void annotated() {
    }

    // TC-68: generics type-parameter completion inside angle brackets
    void generics(
            List< // <caret> TC-68: generic type param — expect String, Integer, User, etc.
                    > a,
            Map<String, // <caret> TC-69: second generic type param — expect type suggestions
                    > b
    ) {
    }
}
