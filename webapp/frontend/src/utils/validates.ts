export class ValidationError extends Error {}

type PartiallyNullable<T, K extends keyof T> = Omit<T, K> & {
  [P in keyof Pick<T, K>]: T[P] | null;
};

type PartiallyFormDataEntryValue<T, K extends keyof T> = Omit<T, K> & {
  [P in keyof Pick<T, K>]: FormDataEntryValue;
};

export interface ScheduleFormData {
  readonly title: string;
  readonly description: string;
  readonly attendeeSet: Set<string>;
  readonly startAt: Date;
  readonly endAt: Date;
  readonly meetingRoom: string;
}

export type MaybeScheduleFormData = PartiallyNullable<
  PartiallyFormDataEntryValue<ScheduleFormData, "title" | "description">,
  "title" | "startAt" | "endAt"
>;

export function validateScheduleData(
  data: MaybeScheduleFormData
): asserts data is ScheduleFormData {
  if (typeof data.title !== "string" || typeof data.description !== "string") {
    throw new ValidationError("入力内容が不正です");
  }

  if (data.title === "") {
    throw new ValidationError("タイトルは必須です");
  }

  if (data.startAt === null || data.endAt === null) {
    throw new ValidationError("日付は必須です");
  }
}
