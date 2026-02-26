package completion;

import completion.model.User;

import java.util.ArrayList;
import java.util.List;

import static completion.model.User.buildUser;
import static completion.model.User.consumeUser;

/**
 * Basic combination of completion scenarios.
 * Covers Java-side variants for:
 * TC-1, TC-2, TC-7, TC-8, TC-9, TC-16, TC-19, TC-22, TC-23,
 * TC-24, TC-28, TC-34, TC-41, TC-42, TC-63, TC-71.
 * Based on EX-JV-1.
 */
public class BasicCombo {

    public static void main(String[] args) {
        String userName = "Ann";
        int userAge = 21;
        User user = new User(userName, userAge);

        // <caret> TC-1/TC-7/TC-22/TC-23/TC-71:
        //   - place caret after "user." for basic/member completion
        //   - type 'na' to validate prefix filtering to getName
        //   - while popup is open, use Ctrl+Q for QuickDoc
        String tc28 = user.getName();

        // <caret> TC-2/TC-28/TC-42: Delete 'user' below after '=',
        //   invoke smart completion; expect user/new User(...)/buildUser(...)
        User u1 = user;

        // <caret> TC-8/TC-9/TC-16/TC-24/TC-34/TC-63:
        //   - after '(' and after ',' verify auto-popup/parameter completion
        //   - verify second argument suggestions (Int), and parameter info hint
        User u2 = buildUser(userName, userAge);

        // <caret> TC-41/TC-42: Delete 'u2' below inside consumeUser(...), invoke completion;
        //   selecting consumeUser(...) should keep caret in function-call context
        consumeUser(u2);

        List<String> list = new ArrayList<>();
        // <caret> TC-24: Delete '\"Ann\"' below inside add(...), invoke completion;
        //   expect String variables/values
        list.add("Ann");

        System.out.println(keywordReturnScenario());
    }

    static int keywordReturnScenario() {
        // <caret> TC-19: In method body type 'ret' and invoke completion; expect 'return'.
        return 1;
    }
}
