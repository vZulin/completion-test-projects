// Test cases covered: TC-67 (TS variant)
// Source: EX-TS-9

type Box<T> = { v: T };

let b: Box< // <caret> TC-67: generics — completion inside angle brackets should suggest type arguments
>
