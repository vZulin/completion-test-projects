package completion;

import completion.model.User;

import java.util.List;

import static completion.model.User.buildUser;
import static completion.model.User.consumeUser;

/**
 * Basic combination of completion scenarios.
 * Covers Java-side variants for:
 * TC-1, TC-2, TC-7, TC-8, TC-9, TC-16, TC-19, TC-22, TC-23,
 * TC-24, TC-28, TC-34, TC-41, TC-42, TC-63, TC-71, TC-127,
 * TC-132, TC-136, TC-137.
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
        //   - choose getName from completion and verify method parentheses insertion
        //   - while popup is open, use Ctrl+Q for QuickDoc
        String memberName = user.getName();

        // <caret> TC-2/TC-28: Delete 'user' below after '=',
        //   invoke smart completion; expect user/new User(...)/buildUser(...)
        User u1 = user;

        // <caret> TC-8/TC-9/TC-16/TC-24/TC-34/TC-63:
        //   - after '(' and after ',' verify auto-popup/parameter completion
        //   - verify second argument suggestions (Int), and parameter info hint
        User u2 = buildUser(userName, userAge);

        // <caret> TC-41: Delete 'u2' below inside consumeUser(...), invoke completion;
        //   selecting consumeUser(...) should keep caret in function-call context
        consumeUser(u2);

        // <caret> TC-42: Delete constructor call below after 'new ';
        //   invoke completion, choose User(...), verify caret lands in expected constructor position.
        User created = new User(userName, userAge);

        List<String> list = new java.util.ArrayList<>();
        // <caret> TC-24: Delete '\"Ann\"' below inside add(...), invoke completion;
        //   expect String variables/values
        list.add("Ann");

        System.out.println(keywordReturnScenario() + created.getAge());
    }

    static void classNameCommitByDotScenario() {
        // <caret> TC-127: Delete 'r.class' below so the expression token is 'Use'.
        //   Invoke completion, select User, then press '.' to accept it and open static member completion.
        Class<User> userClass = User.class;

        System.out.println(userClass.getName());
    }

    static void consume(List<String> value) {
    }

    static void expectedTypeArgumentConstructorCompletion() {
        // <caret> TC-132: Remove 'java.util.' from the constructor call below and delete 't<>())'
        //   so the argument ends with 'new ArrayLis'. Accept ArrayList and verify constructor insertion
        //   plus import java.util.ArrayList.
        consume(new java.util.ArrayList<>());
    }

    static void fileConstructorParameterTemplateCompletion() {
        // <caret> TC-136: Remove 'java.io.' from the constructor call below and delete 'e(".");'
        //   so the expression ends with 'new Fil'. Accept File and verify constructor parentheses
        //   plus parameter info/template.
        java.io.File file = new java.io.File(".");

        System.out.println(file.getPath());
    }

    static void anonymousClassConstructorCompletion() {
        // <caret> TC-137: Remove 'java.util.' from the constructor call below, then delete from
        //   'rator<>() {' through the matching '};' so the right-hand side is 'new Compa'.
        //   Accept Comparator<Integer> and verify the anonymous class body is generated.
        java.util.Comparator<Integer> comparator = new java.util.Comparator<>() {
            @Override
            public int compare(Integer left, Integer right) {
                return Integer.compare(left, right);
            }
        };

        System.out.println(comparator.compare(1, 2));
    }

    static int keywordReturnScenario() {
        // <caret> TC-19: In method body type 'ret' and invoke completion; expect 'return'.
        return 1;
    }
}
