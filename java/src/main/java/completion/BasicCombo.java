package completion;

import completion.model.User;

import java.util.ArrayList;
import java.util.List;

import static completion.model.User.buildUser;
import static completion.model.User.consumeUser;

/**
 * Basic combination of completion scenarios.
 * Covers: TC-28 (member completion), TC-53 (smart completion).
 * Based on EX-JV-1.
 */
public class BasicCombo {

    public static void main(String[] args) {
        String userName = "Ann";
        int userAge = 21;
        User user = new User();
        User user2 = new User()

        user.getAge() // <caret> TC-28: member completion — expect getName, getAge, name, age
        ;

        User u1 = // <caret> TC-53: smart completion — expect user, new User(), buildUser(...)
                null;

        User u2 = buildUser( // <caret> TC-28: args completion — String name, then int age
        );

        consumeUser( // <caret> TC-53: expected type — expect user, u1, u2
        );

        List<String> list = new ArrayList<>();
        list.add( // <caret> TC-28: argument completion for String parameter
        );
    }
}
