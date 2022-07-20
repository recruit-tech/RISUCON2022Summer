const assert = require("node:assert");
const test = require("node:test");
const { timeIn } = require("./time_in");

test("`timeIn` function", async (t) => {
  const targetStartAt = new Date(Date.UTC(2022, 0, 21, 0, 0, 0, 0));
  const targetEndAt = new Date(Date.UTC(2022, 0, 22, 0, 0, 0, 0));

  await t.test("昨日の予定", () => {
    const start = new Date(Date.UTC(2022, 0, 20, 0, 0, 0, 0));
    const end = new Date(Date.UTC(2022, 0, 20, 0, 30, 0, 0));

    const actual = timeIn(start, end, targetStartAt, targetEndAt);
    const expect = false;
    assert.strictEqual(actual, expect);
  });

  await t.test("ゆく年くる年", () => {
    const start = new Date(Date.UTC(2022, 0, 20, 23, 50, 0, 0));
    const end = new Date(Date.UTC(2022, 0, 21, 0, 10, 0, 0));

    const actual = timeIn(start, end, targetStartAt, targetEndAt);
    const expect = true;
    assert.strictEqual(actual, expect);
  });

  await t.test("当日の予定", () => {
    const start = new Date(Date.UTC(2022, 0, 21, 0, 0, 0, 0));
    const end = new Date(Date.UTC(2022, 0, 21, 12, 44, 56, 0));

    const actual = timeIn(start, end, targetStartAt, targetEndAt);
    const expect = true;
    assert.strictEqual(actual, expect);
  });

  await t.test("深夜作業", () => {
    const start = new Date(Date.UTC(2022, 0, 21, 23, 59, 59, 0));
    const end = new Date(Date.UTC(2022, 0, 22, 3, 30, 0, 0));

    const actual = timeIn(start, end, targetStartAt, targetEndAt);
    const expect = true;
    assert.strictEqual(actual, expect);
  });

  await t.test("ぶっ続け", () => {
    const start = new Date(Date.UTC(2022, 0, 20, 12, 30, 0, 0));
    const end = new Date(Date.UTC(2022, 0, 22, 12, 30, 0, 0));

    const actual = timeIn(start, end, targetStartAt, targetEndAt);
    const expect = true;
    assert.strictEqual(actual, expect);
  });

  await t.test("明日の予定", () => {
    const start = new Date(Date.UTC(2022, 0, 22, 12, 30, 0, 0));
    const end = new Date(Date.UTC(2022, 0, 22, 13, 30, 0, 0));

    const actual = timeIn(start, end, targetStartAt, targetEndAt);
    const expect = false;
    assert.strictEqual(actual, expect);
  });

  await t.test("左端範囲外", () => {
    const start = new Date(Date.UTC(2022, 0, 20, 23, 59, 58, 0));
    const end = new Date(Date.UTC(2022, 0, 20, 23, 59, 59, 0));

    const actual = timeIn(start, end, targetStartAt, targetEndAt);
    const expect = false;
    assert.strictEqual(actual, expect);
  });

  await t.test("左端範囲内", () => {
    const start = new Date(Date.UTC(2022, 0, 20, 23, 59, 59, 0));
    const end = new Date(Date.UTC(2022, 0, 21, 0, 0, 0, 0));

    const actual = timeIn(start, end, targetStartAt, targetEndAt);
    const expect = true;
    assert.strictEqual(actual, expect);
  });

  await t.test("右端範囲内", () => {
    const start = new Date(Date.UTC(2022, 0, 21, 23, 59, 59, 0));
    const end = new Date(Date.UTC(2022, 0, 22, 0, 0, 0, 0));

    const actual = timeIn(start, end, targetStartAt, targetEndAt);
    const expect = true;
    assert.strictEqual(actual, expect);
  });

  await t.test("右端範囲外", () => {
    const start = new Date(Date.UTC(2022, 0, 22, 0, 0, 0, 0));
    const end = new Date(Date.UTC(2022, 0, 22, 0, 0, 1, 0));

    const actual = timeIn(start, end, targetStartAt, targetEndAt);
    const expect = false;
    assert.strictEqual(actual, expect);
  });
});
