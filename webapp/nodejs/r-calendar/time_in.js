/**
 * @arg {Date} start
 * @arg {Date} end
 * @arg {Date} targetStartAt
 * @arg {Date} targetEndAt
 * @return {boolean}
 */
const timeIn = (start, end, targetStartAt, targetEndAt) =>
  !(
    start.getTime() >= targetEndAt.getTime() ||
    end.getTime() < targetStartAt.getTime()
  );

module.exports = {
  timeIn,
};
