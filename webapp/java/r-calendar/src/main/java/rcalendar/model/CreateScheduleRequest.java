package rcalendar.model;

import am.ik.yavi.builder.ValidatorBuilder;
import am.ik.yavi.core.ConstraintViolations;
import am.ik.yavi.core.Validator;
import am.ik.yavi.core.ViolationMessage;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public record CreateScheduleRequest(
        @JsonProperty("attendees") List<String> attendees,
        @JsonProperty("start_at") Long startAt,
        @JsonProperty("end_at") Long endAt,
        @JsonProperty("title") String title,
        @JsonProperty("meeting_room") String meetingRoom,
        @JsonProperty("description") String description
) {
    private static final Validator<CreateScheduleRequest> validator = ValidatorBuilder.of(CreateScheduleRequest.class)
            .constraint(CreateScheduleRequest::title, "title", c->c.notBlank()
                    .message("タイトルを設定してください"))
            .constraint(CreateScheduleRequest::startAt, "start_at", c->c.greaterThan(0L)
                    .message("時間の指定が不正です"))
            .constraint(CreateScheduleRequest::endAt, "end_at", c->c.greaterThan(0L)
                    .message("時間の指定が不正です"))
            .constraintOnTarget(r -> r.startAt <= r.endAt, "", ViolationMessage.of(null, "終了時間は開始時間よりも後に設定してください"))
            .constraint(CreateScheduleRequest::attendees, "attendees", c -> c.notEmpty()
                    .message("参加者を指定してください"))
            .failFast(true)
            .build();


    public ConstraintViolations validate() {
        return validator.validate(this);
    }
}
