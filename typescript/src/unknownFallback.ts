// Test cases covered: TC-38 (TS variant)
// Source: EX-TS-6

declare const dyn: any;

dyn. // <caret> TC-38: unknown type fallback — smart completion should not crash; results may be limited
