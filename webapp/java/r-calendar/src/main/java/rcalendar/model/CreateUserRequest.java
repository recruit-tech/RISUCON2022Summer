package rcalendar.model;

import am.ik.yavi.builder.ValidatorBuilder;
import am.ik.yavi.core.ConstraintViolations;
import am.ik.yavi.core.Validator;
import com.fasterxml.jackson.annotation.JsonProperty;

public record CreateUserRequest(
        @JsonProperty("email") String email,
        @JsonProperty("name") String name,
        @JsonProperty("password") String password
) {
    private static final Validator<CreateUserRequest> validator = ValidatorBuilder.of(CreateUserRequest.class)
            .constraint(CreateUserRequest::email, "email", c -> c.pattern("^[a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$")
                    .message("メールアドレスの形式が不正です"))
            .constraint(CreateUserRequest::password, "password", c -> c.pattern("^[a-zA-Z0-9.!@#$%^&-]{8,64}$")
                    .message("パスワードは数字、アルファベット、記号(!@#$%^&-)から8~64文字以内で指定してください"))
            .failFast(true)
            .build();

    public ConstraintViolations validate() {
        return validator.validate(this);
    }
}
