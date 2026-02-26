package completion;

import completion.model.User;

import java.util.ArrayList;
import java.util.List;

import static completion.model.User.buildUser;
import static completion.model.User.consumeUser;

/**
 * Basic combination of completion scenarios.
 * Covers: TC-24 (member completion), TC-42 (smart completion).
 * Based on EX-JV-1.
 */
public class BasicCombo {

    public static void main(String[] args) {
        String userName = "Ann";
        int userAge = 21;
        User user = new User(userName, userAge);

        // <caret> TC-24: Delete 'getName()' below, place caret after "user." and invoke completion;
        //   expect getName(), getAge(), toString(), etc.
        String tc28 = user.getName();

        // <caret> TC-42: Delete 'user' below after '=', invoke smart completion;
        //   expect user, new User(...), buildUser(...)
        User u1 = user;

        // <caret> TC-24: Delete args below inside buildUser(...), invoke completion;
        //   expect String for first arg, Int for second arg
        User u2 = buildUser(userName, userAge);

        // <caret> TC-42: Delete 'u2' below inside consumeUser(...), invoke completion;
        //   expect user, u1, u2 by expected type User
        consumeUser(u2);

        List<String> list = new ArrayList<>();
        // <caret> TC-24: Delete '\"Ann\"' below inside add(...), invoke completion;
        //   expect String variables/values
        list.add("Ann");
    }
}
