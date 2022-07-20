package rcalendar.resource;

import org.junit.jupiter.api.Test;

import java.time.LocalDateTime;

import static java.time.Month.*;
import static org.assertj.core.api.Assertions.assertThat;
import static rcalendar.resource.CalendarResource.timeIn;

class TimeInTest {
    final LocalDateTime start = LocalDateTime.of(2022, JANUARY, 21, 0,0,0);
    final LocalDateTime end = LocalDateTime.of(2022, JANUARY, 22, 0,0,0);

    @Test
    void 昨日の予定() {
        assertThat(timeIn(
                LocalDateTime.of(2022, JANUARY, 20, 0, 0, 0),
                LocalDateTime.of(2022, JANUARY, 20, 0, 30, 0),
                start, end))
                .isFalse();
    }

    @Test
    void ゆく年くる年() {
        assertThat(timeIn(
                LocalDateTime.of(2022, JANUARY, 20, 23, 50, 0),
                LocalDateTime.of(2022, JANUARY, 21, 0, 10, 0),
                start, end))
                .isTrue();
    }

    @Test
    void 当日の予定() {
        assertThat(timeIn(
                LocalDateTime.of(2022, JANUARY, 21, 0, 50, 0),
                LocalDateTime.of(2022, JANUARY, 21, 12, 44, 56),
                start, end))
                .isTrue();
    }

    @Test
    void 深夜作業() {
        assertThat(timeIn(
                LocalDateTime.of(2022, JANUARY, 21, 23, 59, 59),
                LocalDateTime.of(2022, JANUARY, 22, 3, 30, 0),
                start, end))
                .isTrue();
    }

    @Test
    void ぶっ続け() {
        assertThat(timeIn(
                LocalDateTime.of(2022, JANUARY, 20, 12, 30, 0),
                LocalDateTime.of(2022, JANUARY, 22, 12, 30, 0),
                start, end))
                .isTrue();
    }

    @Test
    void 明日の予定() {
        assertThat(timeIn(
                LocalDateTime.of(2022, JANUARY, 22, 12, 30, 0),
                LocalDateTime.of(2022, JANUARY, 22, 13, 30, 0),
                start, end))
                .isFalse();
    }

    @Test
    void 左端範囲外() {
        assertThat(timeIn(
                LocalDateTime.of(2022, JANUARY, 20, 23, 59, 58),
                LocalDateTime.of(2022, JANUARY, 20, 23, 59, 59),
                start, end))
                .isFalse();
    }

    @Test
    void 左端範囲内() {
        assertThat(timeIn(
                LocalDateTime.of(2022, JANUARY, 20, 23, 59, 59),
                LocalDateTime.of(2022, JANUARY, 21, 0, 0, 0),
                start, end))
                .isTrue();
    }

    @Test
    void 右端範囲内() {
        assertThat(timeIn(
                LocalDateTime.of(2022, JANUARY, 21, 23, 59, 59),
                LocalDateTime.of(2022, JANUARY, 22, 0, 0, 0),
                start, end))
                .isTrue();
    }

    @Test
    void 右端範囲外() {
        assertThat(timeIn(
                LocalDateTime.of(2022, JANUARY, 22, 0, 0, 0),
                LocalDateTime.of(2022, JANUARY, 22, 0, 0, 1),
                start, end))
                .isFalse();
    }
}