// Test cases covered: TC-40 (TS variant)
// Source: EX-TS-7

function doWork(): void {}

// <caret> TC-40: Delete 'ork()' below so only 'doW' remains,
//   then press '(' to accept completion; should produce doWork()
doWork();

export {};
