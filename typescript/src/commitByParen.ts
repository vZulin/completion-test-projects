// Test cases covered: TC-50 (TS variant)
// Source: EX-TS-7

function doWork(): void {}

doW // <caret> TC-50: commit by '(' — accept completion with '(' should produce doWork()
