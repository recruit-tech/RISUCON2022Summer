import cookieParser from "cookie-parser";
import express from "express";
import jwt from "jsonwebtoken";
import morgan from "morgan";
import multer from "multer";
import mysql, { type PoolConfig } from "mysql";
import { exec as _exec } from "node:child_process";
import crypto from "node:crypto";
import { mkdtemp, rm, writeFile } from "node:fs/promises";
import path from "node:path";
import { URL } from "node:url";
import { promisify } from "node:util";
import { ulid } from "ulid";
import { timeIn } from "./time_in";

const exec = promisify(_exec);

process.env.TZ = "Etc/Universal";

const PORT = process.env.CALENDAR_PORT ?? 3000;
const ORECOCO_RESERVE_URL =
  process.env.ORECOCO_RESERVE_URL ?? "http://localhost:3003";
const JWT_SECRET_KEY = "secret_key";
const EMAIL_REGEXP = `^[a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$`;
const PASSWORD_REGEXP = `^[a-zA-Z0-9.!@#$%^&-]{8,64}$`;
const ALLOWED_FILE_TYPE_AND_MIME = {
  "image/jpeg": "JPEG image data",
  "image/png": "PNG image data",
  "image/gif": "GIF image data",
  "image/bmp": "PB bitmap",
} as const;
type MimeType = keyof typeof ALLOWED_FILE_TYPE_AND_MIME;
const POOL_CONFIG: PoolConfig = {
  host: process.env.MYSQL_HOST ?? "127.0.0.1",
  port: Number(process.env.MYSQL_PORT ?? 3306),
  user: process.env.MYSQL_USER ?? "r-isucon",
  password: process.env.MYSQL_PASS ?? "r-isucon",
  database: process.env.MYSQL_DBNAME ?? "r-calendar",
};
const session = new Map<string, Symbol>();

let orecocoReserveToken: string;

const app = express();
const upload = multer();
const db = mysql.createPool(POOL_CONFIG);
app.set("db", db);

app.use(morgan("combined"));
app.use(cookieParser());
app.use(express.urlencoded({ extended: true }));
app.use(express.json());

const authRouter = express.Router();
interface AuthLocals {
  readonly jti: string;
  readonly userId: string;
}
authRouter.use((req, res, next) => {
  const jwtString: unknown = req.cookies["jwt"];
  if (typeof jwtString !== "string") {
    res.status(401).json({
      message: "認証されていないユーザーです",
    });
    return;
  }

  try {
    const payload = jwt.verify(jwtString, JWT_SECRET_KEY);
    if (typeof payload === "string") throw new Error();
    res.locals.jti = payload.jti;
    res.locals.userId = payload["user_id"];
  } catch (e) {
    console.error(e);
    res.status(401).json({
      message: "認証されていないユーザーです",
    });
    return;
  }

  next();
});

class ApplicationError extends Error {
  readonly statusCode: number;
  constructor(statusCode: number, message: string) {
    super(message);
    this.statusCode = statusCode;
  }
}

function getConnection(db: mysql.Pool) {
  return new Promise<mysql.PoolConnection>((resolve, reject) => {
    db.getConnection((err, connection) => {
      if (err) {
        reject(err);
      } else {
        resolve(connection);
      }
    });
  });
}

const query =
  (connection: mysql.PoolConnection) =>
  <T = unknown>(query: string, values?: unknown[]) =>
    new Promise<T>((resolve, reject) => {
      connection.query(query, values, (err, result) => {
        if (err) {
          reject(err);
        } else {
          resolve(result);
        }
      });
    });

const beginTransaction = (connection: mysql.PoolConnection) =>
  new Promise((resolve, reject) => {
    connection.beginTransaction((err) => {
      if (err) {
        reject(err);
      } else {
        resolve(undefined);
      }
    });
  });

const commit = (connection: mysql.PoolConnection) =>
  new Promise((resolve, reject) => {
    connection.commit((err) => {
      if (err) {
        reject(err);
      } else {
        resolve(undefined);
      }
    });
  });

const rollback = (connection: mysql.PoolConnection) =>
  new Promise((resolve, reject) => {
    connection.rollback((err) => {
      if (err) {
        reject(err);
      } else {
        resolve(undefined);
      }
    });
  });

const generateSessionID = () => {
  const LETTERS =
    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  const possibleSessionIDs = [];
  for (let i = 0; i < 10000; i++) {
    const code = crypto
      .randomBytes(128)
      .reduce((acc, byte) => acc + LETTERS.at(byte % LETTERS.length), "");
    possibleSessionIDs.push(code);
  }

  return possibleSessionIDs[~~(Math.random() * possibleSessionIDs.length)];
};

interface ErrorResBody {
  readonly message: string;
}

app.post<undefined, { readonly language: "nodejs" } | ErrorResBody>(
  "/initialize",
  async (_, res) => {
    try {
      const db_dir = path.resolve("..", "..", "sql");
      const exec_files = [
        "r-calendar-0_Schema.sql",
        "r-calendar-1_DummyUserData.sql",
      ].map((file) => path.join(db_dir, file));
      for (const exec_file of exec_files) {
        await exec(
          `mysql -h ${POOL_CONFIG.host} -u ${POOL_CONFIG.user} -p${POOL_CONFIG.password} -P ${POOL_CONFIG.port} ${POOL_CONFIG.database} < ${exec_file}`
        );
      }

      const response = await fetch(
        new URL("/initialize", ORECOCO_RESERVE_URL).toString(),
        {
          method: "POST",
        }
      );

      if (response.status !== 200) {
        console.error(
          `orecoco_reserve initialize failed returned http status ${response.status}`
        );
        res.status(500).json({
          message: "initialize failed",
        });
        return;
      }

      const json = (await response.json()) as unknown;
      type OrecocoReserveInitializeResponseBody = {
        readonly token: string;
      };

      function isOrecocoReserveInitializeResponseBody(
        data: unknown
      ): data is OrecocoReserveInitializeResponseBody {
        return (
          typeof data === "object" &&
          data !== null &&
          "token" in data &&
          typeof (
            data as Record<keyof OrecocoReserveInitializeResponseBody, unknown>
          ).token === "string"
        );
      }
      if (!isOrecocoReserveInitializeResponseBody(json)) {
        res.status(500).json({
          message: "サーバー側のエラーです",
        });
        return;
      }
      app.set("orecocoReserveToken", json.token);

      session.clear();

      res.json({
        language: "nodejs",
      });
    } catch (e) {
      console.error(e);
      res.status(500).json({
        message: "サーバー側のエラーです",
      });
    }
  }
);

interface User {
  readonly id: string;
  readonly email: string;
  readonly name: string;
  readonly password: string;
  readonly image_binary?: Buffer | null;
  readonly created_at: string;
}

interface UserResponse extends Pick<User, "id" | "email" | "name"> {
  readonly icon?: string;
}

interface Schedule {
  readonly id: string;
  readonly title: string;
  readonly description: string;
  readonly schedule_attendee: string;
  readonly start_at: Date;
  readonly end_at: Date;
}

interface ScheduleResponse
  extends Pick<Schedule, "id" | "title" | "description"> {
  readonly start_at: number;
  readonly end_at: number;
  readonly attendees: UserResponse[];
  readonly meeting_room: string;
}

app.post<
  undefined,
  ErrorResBody,
  Partial<Pick<User, "name" | "email" | "password">>
>("/user", async (req, res, next) => {
  if (!req.body.email?.match(new RegExp(EMAIL_REGEXP))) {
    res.status(400).json({
      message: "メールアドレスの形式が不正です",
    });
  }

  let connection: mysql.PoolConnection | undefined;
  try {
    connection = await getConnection(db);
    const users = await query(connection)<User[]>(
      "SELECT * FROM user WHERE email = ?",
      [req.body.email]
    );
    if (users.length > 0) {
      res.status(400).json({
        message: "指定されたメールアドレスはすでに利用されています",
      });
      return;
    }
    if (req.body.password === undefined || req.body.password === "") {
      res.status(400).json({
        message: "パスワードを設定してください",
      });
      return;
    }
    if (!req.body.password.match(new RegExp(PASSWORD_REGEXP))) {
      res.status(400).json({
        message:
          "パスワードは数字、アルファベット、記号(!@#$%^&-)から8~64文字以内で指定してください",
        password: req.body.password,
      } as any);
      return;
    }

    const hash = crypto
      .createHash("sha256")
      .update(req.body.password)
      .digest("hex");
    const uid = ulid();
    await query(connection)(
      "INSERT INTO user(id, email, name, password) VALUES (?, ?, ?, ?)",
      [uid, req.body.email, req.body.name, hash]
    );

    const newSessionID = generateSessionID();
    const token = jwt.sign(
      {
        user_id: uid,
        jti: newSessionID,
      },
      JWT_SECRET_KEY
    );

    session.set(token, Symbol());

    res
      .status(201)
      .cookie("jwt", token, {
        httpOnly: true,
        path: "/",
        expires: new Date(Date.now() + 24 * 60 * 60 * 1000),
      })
      .send();
  } catch (e) {
    console.error(e);
    res.status(500).json({
      message: "サーバー側のエラーです",
    });
  } finally {
    connection?.release();
  }
});

app.post<undefined, ErrorResBody, Partial<Pick<User, "email" | "password">>>(
  "/login",
  async (req, res) => {
    let connection: mysql.PoolConnection | undefined;
    try {
      connection = await getConnection(db);
      const [user] = await query(connection)<User[]>(
        "SELECT * FROM user WHERE email = ?",
        [req.body.email]
      );
      if (user === undefined || typeof req.body.password !== "string") {
        res.status(400).json({
          message: "ユーザー名またはパスワードが不正です",
        });
        return;
      }

      const hash = crypto
        .createHash("sha256")
        .update(req.body.password)
        .digest("hex");

      if (hash !== user.password) {
        res.status(400).json({
          message: "ユーザー名またはパスワードが不正です",
        });
        return;
      }

      const newSessionID = generateSessionID();
      const token = jwt.sign(
        {
          user_id: user.id,
          jti: newSessionID,
        },
        JWT_SECRET_KEY
      );

      session.set(token, Symbol());

      res
        .status(201)
        .cookie("jwt", token, {
          httpOnly: true,
          path: "/",
          expires: new Date(Date.now() + 24 * 60 * 60 * 1000),
        })
        .send();
    } catch (e) {
      console.error(e);
      res.status(500).json({
        message: "サーバー側のエラーです",
      });
    } finally {
      connection?.release();
    }
  }
);

authRouter.post<
  undefined,
  {},
  Partial<Pick<User, "name" | "email" | "password">>,
  {},
  AuthLocals
>("/logout", async (_, res) => {
  const { jti } = res.locals;
  session.delete(jti);
  res.status(201).send();
});

authRouter.get<undefined, UserResponse | ErrorResBody, {}, {}, AuthLocals>(
  "/me",
  async (_, res) => {
    const { userId } = res.locals;
    let connection: mysql.PoolConnection | undefined;
    try {
      connection = await getConnection(db);
      const [me] = await query(connection)<User[]>(
        "SELECT * FROM user WHERE id = ?",
        [userId]
      );
      if (me === undefined) {
        res.status(404).json({
          message: "ユーザーが見つかりませんでした",
        });
        return;
      }

      res.json({
        id: me.id,
        email: me.email,
        name: me.name,
        icon: me.image_binary ? `/icon/${me.id}` : undefined,
      });
    } catch (e) {
      console.error(e);
      res.status(500).json({
        message: "サーバー側のエラーです",
      });
    } finally {
      connection?.release();
    }
  }
);

authRouter.get<
  { readonly userId: string },
  UserResponse | ErrorResBody,
  {},
  {},
  AuthLocals
>("/user/:userId", async (req, res) => {
  const { userId } = req.params;
  let connection: mysql.PoolConnection | undefined;
  try {
    connection = await getConnection(db);
    const [user] = await query(connection)<User[]>(
      "SELECT * FROM user WHERE id = ?",
      [userId]
    );
    if (user === undefined) {
      res.status(404).json({
        message: "ユーザーが見つかりませんでした",
      });
      return;
    }
    res.json({
      id: user.id,
      email: user.email,
      name: user.name,
      icon: user.image_binary ? `/icon/${user.id}` : undefined,
    });
  } catch (e) {
    console.error(e);
    res.status(500).json({
      message: "サーバー側のエラーです",
    });
  } finally {
    connection?.release();
  }
});

authRouter.get<
  { readonly userId: string },
  Buffer | ErrorResBody,
  {},
  {},
  AuthLocals
>("/user/icon/:userId", async (req, res) => {
  const { userId } = req.params;
  let connection: mysql.PoolConnection | undefined;
  try {
    connection = await getConnection(db);
    const [user] = await query(connection)<User[]>(
      "SELECT * FROM user WHERE id = ?",
      [userId]
    );
    if (user === undefined) {
      res.status(404).json({
        message: "ユーザーが見つかりませんでした",
      });
      return;
    }

    if (user.image_binary === undefined || user.image_binary === null) {
      res.status(404).json({
        message: "アイコンが登録されていません",
      });
      return;
    }

    const tempDir = await mkdtemp("/tmp/");
    const tempFile = path.join(tempDir, "icon");
    await writeFile(tempFile, user.image_binary);
    const mimeType = await detectMimeType(tempFile);
    if (mimeType === undefined) {
      throw new Error("invalid icon mime-type");
    }

    res.contentType(mimeType).send(user.image_binary);
  } catch (e) {
    console.error(e);
    res.status(500).json({
      message: "サーバー側のエラーです",
    });
  } finally {
    connection?.release();
  }
});

authRouter.put<
  undefined,
  ErrorResBody,
  Partial<Pick<User, "email" | "name" | "password">>,
  {},
  AuthLocals
>("/me", async (req, res) => {
  const { userId } = res.locals;
  if (!req.body.email?.match(new RegExp(EMAIL_REGEXP))) {
    res.status(400).json({
      message: "メールアドレスの形式が不正です",
    });
  }

  let connection: mysql.PoolConnection | undefined;
  try {
    connection = await getConnection(db);
    const [me] = await query(connection)<User[]>(
      "SELECT * FROM user WHERE id = ?",
      [userId]
    );
    if (me === undefined) {
      res.status(404).json({
        message: "ユーザーが見つかりませんでした",
      });
      return;
    }
    const { email, name, password } = req.body;
    if (email === undefined || email === "") {
      res.status(400).json({
        message: "メールアドレスを設定してください",
      });
      return;
    }
    if (name === undefined || name === "") {
      res.status(400).json({
        message: "ユーザー名を設定してください",
      });
      return;
    }
    if (password === undefined || password === "") {
      res.status(400).json({
        message: "パスワードを設定してください",
      });
      return;
    }
    if (!password.match(new RegExp(PASSWORD_REGEXP))) {
      res.status(400).json({
        message:
          "パスワードは数字、アルファベット、記号(!@#$%^&-)から8~64文字以内で指定してください",
        password: req.body.password,
      } as any);
      return;
    }

    const hash = crypto.createHash("sha256").update(password).digest("hex");
    await query(connection)(
      "UPDATE user SET email = ?, name = ?, password = ? WHERE id = ?",
      [email, name, hash, userId]
    );
    res.status(200).send();
  } catch (e) {
    console.error(e);
    res.status(500).json({
      message: "サーバー側のエラーです",
    });
  } finally {
    connection?.release();
  }
});

authRouter.put<
  {},
  ErrorResBody,
  Partial<Pick<User, "email" | "name" | "password">>,
  {},
  AuthLocals
>("/me/icon", upload.single("icon"), async (req, res) => {
  const { userId } = res.locals;

  let connection: mysql.PoolConnection | undefined;
  try {
    connection = await getConnection(db);
    const [me] = await query(connection)<User[]>(
      "SELECT * FROM user WHERE id = ?",
      [userId]
    );
    if (me === undefined) {
      res.status(404).json({
        message: "ユーザーが見つかりませんでした",
      });
      return;
    }

    const icon = req.file?.buffer;
    if (icon === undefined) {
      console.error("failed to get image from request: icon is undefined");
      return res.status(400).json({
        message: "アイコンがリクエストに含まれていません",
      });
    }

    let tempDir: string | undefined = undefined;
    let mimeType: MimeType | undefined = undefined;
    try {
      tempDir = await mkdtemp("/tmp");
      const tempFile = path.join(tempDir, "icon");
      await writeFile(tempFile, icon);

      mimeType = await detectMimeType(tempFile);
    } catch (e) {
      console.error(e);
      throw new Error("failed to check content type of icon");
    } finally {
      if (tempDir !== undefined) {
        await rm(tempDir, { recursive: true, force: true });
      }
    }

    if (mimeType === undefined) {
      res.status(400).json({
        message:
          "アイコンに指定できるのはjpeg, png, gifまたはbmpの画像ファイルのみです",
      });
      return;
    }

    await query(connection)("UPDATE user SET image_binary = ? WHERE id = ?", [
      icon,
      userId,
    ]);
    res.status(200).send();
  } catch (e) {
    console.error(e);
    res.status(500).json({
      message: "サーバー側のエラーです",
    });
  } finally {
    connection?.release();
  }
});

authRouter.get<
  undefined,
  { readonly users: UserResponse[] } | ErrorResBody,
  {},
  { readonly query: string },
  AuthLocals
>("/user", async (req, res) => {
  const { query: searchQuery = "" } = req.query;

  if (searchQuery === "") {
    res.status(400).json({
      message: "検索条件を指定してください",
    });
    return;
  }

  let connection: mysql.PoolConnection | undefined;
  try {
    connection = await getConnection(db);
    const users = await query(connection)<User[]>(
      // ユーザーが作成された順にソートして返す
      "SELECT * FROM user WHERE email LIKE ? OR name LIKE ? ORDER BY id",
      [`${searchQuery}%`, `${searchQuery}%`]
    );
    if (users.length === 0) {
      res.status(204).json({
        message: "ユーザーが見つかりませんでした",
      });
      return;
    }
    res.json({
      users: users.map((user) => ({
        id: user.id,
        email: user.email,
        name: user.name,
        icon: user.image_binary ? `/icon/${user.id}` : undefined,
      })),
    });
  } catch (e) {
    console.error(e);
    res.status(500).json({
      message: "サーバー側のエラーです",
    });
  } finally {
    connection?.release();
  }
});

authRouter.post<
  undefined,
  { readonly id: string } | ErrorResBody,
  Partial<
    Pick<Schedule, "title" | "description"> & {
      readonly start_at: number;
      readonly end_at: number;
      readonly meeting_room: string;
      readonly attendees: string[];
    }
  >
>("/schedule", async (req, res, next) => {
  const {
    title = "",
    description = "",
    start_at,
    end_at,
    meeting_room = "",
    attendees = [],
  } = req.body;

  if (title === "") {
    res.status(400).json({
      message: "タイトルを設定してください",
    });
    return;
  }

  if (!start_at || !end_at) {
    res.status(400).json({
      message: "時間の指定が不正です",
    });
    return;
  }

  if (start_at >= end_at) {
    res.status(400).json({
      message: "終了時間は開始時間よりも後に設定してください",
    });
    return;
  }

  if (attendees.length === 0) {
    res.status(400).json({
      message: "参加者を指定してください",
    });
    return;
  }

  let connection: mysql.PoolConnection | undefined;
  try {
    connection = await getConnection(db);
    const attendeeSet = new Set<string>();
    for (const attendeeId of attendees) {
      if (attendeeSet.has(attendeeId)) {
        res.status(400).json({
          message: "ユーザーIDに重複が存在しました",
        });
        return;
      }
      const [user] = await query(connection)<User[]>(
        "SELECT * FROM user WHERE id = ?",
        [attendeeId]
      );
      if (user === undefined) {
        res.status(400).json({
          message: "ユーザーを見つけることができませんでした",
        });
        return;
      }
      attendeeSet.add(attendeeId);
    }

    const id = ulid();
    try {
      await beginTransaction(connection);
      const startAt = new Date(start_at * 1000);
      const endAt = new Date(end_at * 1000);
      const scheduleAttendeesID = attendees.join(",");
      await query(connection)(
        "INSERT INTO schedule (id, title, description, schedule_attendee, start_at, end_at) VALUES (?, ?, ?, ?, ?, ?)",
        [id, title, description, scheduleAttendeesID, startAt, endAt]
      );
      if (meeting_room !== "") {
        const { statusCode, body } = await callOrecocoReserve("POST", "/room", {
          schedule_id: id,
          meeting_room_id: meeting_room,
          start_at,
          end_at,
        });
        if (statusCode !== 201) {
          console.error("failed to request orecoco-reserve", body);
          throw new ApplicationError(statusCode, JSON.stringify(body));
        }
      }
      await commit(connection);
    } catch (e) {
      await rollback(connection);
      throw e;
    }
    res.status(201).json({
      id,
    });
    return;
  } catch (e) {
    if (e instanceof ApplicationError) {
      console.error(e.message);
      res.status(e.statusCode).json({
        message: e.message,
      });
      return;
    }
    console.error(e);
    res.status(500).json({
      message: "サーバー側のエラーです",
    });
  } finally {
    connection?.release();
  }
});

authRouter.get<
  { readonly scheduleId: string },
  ScheduleResponse | ErrorResBody,
  {},
  {},
  AuthLocals
>("/schedule/:scheduleId", async (req, res) => {
  const { scheduleId } = req.params;

  let connection: mysql.PoolConnection | undefined;
  try {
    connection = await getConnection(db);
    const [schedule] = await query(connection)<Schedule[]>(
      "SELECT * FROM schedule WHERE id = ?",
      [scheduleId]
    );
    if (schedule === undefined) {
      res.status(400).json({
        message: "スケジュールが見つかりませんでした",
      });
      return;
    }

    const attendees: (Pick<User, "id" | "email" | "name"> & {
      readonly icon?: string;
    })[] = [];
    for (const attendeeId of schedule.schedule_attendee.split(",")) {
      const [user] = await query(connection)<User[]>(
        "SELECT * FROM user WHERE id = ?",
        [attendeeId]
      );
      if (user === undefined) {
        res.status(400).json({
          message: "参加者が見つかりませんでした",
        });
        return;
      }
      attendees.push({
        id: user.id,
        email: user.email,
        name: user.name,
        icon: user.image_binary ? `/icon/${user.id}` : undefined,
      });
    }
    attendees.sort((u1, u2) =>
      u1.email === u2.email ? 0 : u1.email < u2.email ? 1 : -1
    );

    const { statusCode, body } = await callOrecocoReserve("GET", "/room", {
      scheduleId: scheduleId,
    });

    let meeting_room: string = "";
    switch (statusCode) {
      case 200:
        statusCode;
        meeting_room = body.room_id;
        break;

      case 404:
        break;

      default:
        console.error("failed to request orecoco-reserve", body);
        res.status(statusCode).json(body);
        return;
    }

    res.status(200).json({
      id: scheduleId,
      title: schedule.title,
      description: schedule.description,
      start_at: schedule.start_at.getTime() / 1000,
      end_at: schedule.end_at.getTime() / 1000,
      attendees,
      meeting_room,
    });
  } catch (e) {
    console.error(e);
    res.status(500).json({
      message: "サーバー側のエラーです",
    });
  } finally {
    connection?.release();
  }
});

authRouter.put<
  { readonly scheduleId: string },
  ErrorResBody,
  Partial<
    Pick<Schedule, "title" | "description"> & {
      readonly start_at: number;
      readonly end_at: number;
      readonly meeting_room: string;
      readonly attendees: string[];
    }
  >,
  {},
  AuthLocals
>("/schedule/:scheduleId", async (req, res) => {
  const {
    title = "",
    description = "",
    start_at,
    end_at,
    meeting_room = "",
    attendees = [],
  } = req.body;
  const { scheduleId } = req.params;

  if (title === "") {
    res.status(400).json({
      message: "タイトルを設定してください",
    });
    return;
  }

  if (!start_at || !end_at) {
    res.status(400).json({
      message: "時間の指定が不正です",
    });
    return;
  }

  if (start_at >= end_at) {
    res.status(400).json({
      message: "終了時間は開始時間よりも後に設定してください",
    });
    return;
  }

  if (attendees.length === 0) {
    res.status(400).json({
      message: "参加者を指定してください",
    });
    return;
  }

  let connection: mysql.PoolConnection | undefined;
  try {
    connection = await getConnection(db);
    await beginTransaction(connection);

    const attendeeSet = new Set<string>();
    for (const attendeeId of attendees) {
      if (attendeeSet.has(attendeeId)) {
        throw new ApplicationError(400, "ユーザーIDに重複が存在しました");
      }
      const [user] = await query(connection)<User[]>(
        "SELECT * FROM user WHERE id = ?",
        [attendeeId]
      );
      if (user === undefined) {
        throw new ApplicationError(
          400,
          "ユーザーを見つけることができませんでした"
        );
      }
      attendeeSet.add(attendeeId);
    }

    const startAt = new Date(start_at * 1000);
    const endAt = new Date(end_at * 1000);
    const scheduleAttendeesID = attendees.join(",");
    await query(connection)(
      "UPDATE schedule SET title = ?, description = ?, schedule_attendee = ?, start_at = ?, end_at = ? WHERE id = ?",
      [title, description, scheduleAttendeesID, startAt, endAt, scheduleId]
    );
    if (meeting_room !== "") {
      const { statusCode, body } = await callOrecocoReserve("PUT", "/room", {
        schedule_id: scheduleId,
        meeting_room_id: meeting_room,
        start_at,
        end_at,
      });
      if (statusCode !== 200) {
        console.error("failed to request orecoco-reserve", body);
        throw new ApplicationError(statusCode, JSON.stringify(body));
      }
    }
    await commit(connection);
    res.status(200).send();
  } catch (e) {
    if (connection) await rollback(connection);
    if (e instanceof ApplicationError) {
      console.error(e.message);
      res.status(e.statusCode).json({
        message: e.message,
      });
    } else {
      console.error(e);
      res.status(500).json({
        message: "サーバー側のエラーです",
      });
    }
  } finally {
    connection?.release();
  }
});

authRouter.get<
  { readonly userId: string },
  | {
      readonly date: number;
      readonly schedules: ScheduleResponse[];
    }
  | ErrorResBody,
  {},
  { readonly date?: string },
  AuthLocals
>("/calendar/:userId", async (req, res) => {
  const { userId } = req.params;
  const { date: dateString = "" } = req.query;

  let connection: mysql.PoolConnection | undefined;
  try {
    connection = await getConnection(db);
    const [user] = await query(connection)<User[]>(
      "SELECT * FROM user WHERE id = ?",
      [userId]
    );
    if (user === undefined) {
      res.status(404).json({
        message: "ユーザーが見つかりませんでした",
      });
      return;
    }

    const date = new Number(dateString).valueOf();
    if (!Number.isInteger(date)) {
      res.status(400).json({
        message: "指定された日時は不正です",
      });
      return;
    }

    const participationSchedules = await query(connection)<Schedule[]>(
      "SELECT * FROM schedule WHERE schedule_attendee LIKE ?",
      [`%${userId}%`]
    );

    const dateRaw = new Date(date * 24 * 60 * 60 * 1000);
    const startOfTheDay = new Date(
      Date.UTC(
        dateRaw.getUTCFullYear(),
        dateRaw.getUTCMonth(),
        dateRaw.getUTCDate(),
        0,
        0,
        0,
        0
      )
    );
    const endOfTheDay = new Date(
      Date.UTC(
        dateRaw.getUTCFullYear(),
        dateRaw.getUTCMonth(),
        dateRaw.getUTCDate() + 1,
        0,
        0,
        0,
        0
      )
    );

    const participationScheduleIDs = participationSchedules
      .filter((s) => timeIn(s.start_at, s.end_at, startOfTheDay, endOfTheDay))
      .map((s) => s.id);

    if (participationScheduleIDs.length === 0) {
      res.json({
        date,
        schedules: [],
      });
      return;
    }

    const userSchedules = await query(connection)<Schedule[]>(
      "SELECT * FROM schedule WHERE id IN (?) ORDER BY start_at ASC, end_at DESC, id ASC",
      [participationScheduleIDs]
    );

    const schedules: ScheduleResponse[] = [];
    for (const s of userSchedules) {
      const attendees: UserResponse[] = [];
      for (const attendeeId of s.schedule_attendee.split(",")) {
        const [attendee] = await query(connection)<User[]>(
          "SELECT * FROM user WHERE id = ?",
          [attendeeId]
        );
        if (attendee === undefined) {
          res.status(400).json({
            message: "参加者が見つかりませんでした",
          });
          return;
        }

        attendees.push({
          id: attendee.id,
          email: attendee.email,
          name: attendee.name,
          icon: attendee.image_binary ? `/icon/${attendee.id}` : undefined,
        });
      }

      attendees.sort((u1, u2) =>
        u1.email === u2.email ? 0 : u1.email < u2.email ? 1 : -1
      );

      let meeting_room = "";
      const { statusCode, body } = await callOrecocoReserve("GET", "/room", {
        scheduleId: s.id,
      });
      switch (statusCode) {
        case 200:
          meeting_room = body.room_id;
          break;

        case 404:
          break;

        default:
          console.error("failed to request orecoco-reserve", body);
          res.status(statusCode).json(body);
          return;
      }

      schedules.push({
        id: s.id,
        title: s.title,
        description: s.description,
        start_at: s.start_at.getTime() / 1000,
        end_at: s.end_at.getTime() / 1000,
        meeting_room,
        attendees,
      });
    }

    res.status(200).json({
      date: date.valueOf(),
      schedules,
    });
  } catch (e) {
    console.error(e);
    res.status(500).json({
      message: "サーバー側のエラーです",
    });
  } finally {
    connection?.release();
  }
});

const detectMimeType = async (
  filename: string
): Promise<MimeType | undefined> => {
  let mimeType: MimeType | undefined = undefined;
  for (const [mime, formatStr] of Object.entries(ALLOWED_FILE_TYPE_AND_MIME)) {
    const { stdout } = await exec(`file "${filename}"`);
    if (stdout.match(new RegExp(formatStr))) {
      mimeType = mime as MimeType;
    }
  }

  return mimeType;
};

type HttpStatusCode = 200 | 201 | 400 | 404 | 500;
type OrecocoReserveResult<
  SuccessStatusCode extends HttpStatusCode,
  SuccessBody = undefined,
  ErrorBody = unknown
> =
  | {
      readonly statusCode: SuccessStatusCode;
      readonly body: SuccessBody;
    }
  | {
      readonly statusCode: Exclude<HttpStatusCode, SuccessStatusCode>;
      readonly body: ErrorBody;
    };

interface CallOrecocoReserveFunc {
  (
    method: "GET",
    endpoint: "/room",
    payload: { readonly scheduleId: string }
  ): Promise<
    OrecocoReserveResult<200, { readonly room_id: string }, ErrorResBody>
  >;

  (
    method: "POST",
    endpoint: "/room",
    payload: {
      readonly schedule_id: string;
      readonly meeting_room_id: string;
      readonly start_at: number;
      readonly end_at: number;
    }
  ): Promise<OrecocoReserveResult<201, undefined, ErrorResBody>>;

  (
    method: "PUT",
    endpoint: "/room",
    payload: {
      readonly schedule_id: string;
      readonly meeting_room_id: string;
      readonly start_at: number;
      readonly end_at: number;
    }
  ): Promise<OrecocoReserveResult<201, undefined, ErrorResBody>>;
}

const callOrecocoReserve: CallOrecocoReserveFunc = async (
  method,
  endpoint,
  payload
) => {
  const url = new URL(endpoint, ORECOCO_RESERVE_URL);
  if (method === "GET") {
    for (const [key, value] of Object.entries(payload)) {
      url.searchParams.set(key, value.toString());
    }
  }
  const response = await fetch(url, {
    method,
    headers: {
      Authorization: `Bearer ${orecocoReserveToken}`,
      "Content-Type": "application/json",
    },
    body: method !== "GET" ? JSON.stringify(payload) : undefined,
  });

  return {
    statusCode: response.status as any,
    body: await response.json(),
  };
};

app.use("/", authRouter);
app.listen(PORT, () => console.log(`Listening ${PORT}`));
