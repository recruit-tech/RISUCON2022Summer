import express from "express";
import morgan from "morgan";
import mysql, { type PoolConfig } from "mysql";
import { exec as _exec } from "node:child_process";
import { readFile } from "node:fs/promises";
import path from "node:path";
import { promisify } from "node:util";

const exec = promisify(_exec);

const PORT = process.env.ORECOCO_PORT ?? 3000;
const POOL_CONFIG: PoolConfig = {
  host: process.env.MYSQL_HOST ?? "127.0.0.1",
  port: Number(process.env.MYSQL_PORT ?? 3308),
  user: process.env.MYSQL_USER ?? "r-isucon",
  password: process.env.MYSQL_PASS ?? "r-isucon",
  database: process.env.MYSQL_DBNAME ?? "orecoco-reserve",
};

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

async function generateOrecocoToken() {
  const { stdout } = await exec(
    "cat /dev/urandom | base64 | fold -w 32 | head -n 1"
  );
  return stdout.trim();
}

interface ErrorResBody {
  readonly message: string;
}

interface MeetingRoom {
  readonly id: string;
  readonly room_id: string;
  readonly start_at: number;
  readonly end_at: number;
}

interface BaseRequestData {
  readonly schedule_id: string;
  readonly meeting_room_id: string;
  readonly start_at: number;
  readonly end_at: number;
}

const meetingRoomExist = async (name: string) => {
  const file = await readFile("./meeting_room.txt");
  let exist = false;
  for (const roomName of file.toString("utf8").split("\n")) {
    if (roomName === name) {
      exist = true;
    }
  }

  return exist;
};

async function main() {
  const token = await generateOrecocoToken();
  const app = express();
  const db = mysql.createPool(POOL_CONFIG);
  app.set("db", db);

  app.use(morgan("combined"));
  app.use(express.json());

  app.post("/initialize", async (_, res, next) => {
    try {
      const db_dir = path.resolve("..", "..", "sql");
      const exec_files = ["orecoco-reserve_0_Schema.sql"].map((file) =>
        path.join(db_dir, file)
      );
      for (const exec_file of exec_files) {
        await exec(
          `mysql -h ${POOL_CONFIG.host} -u ${POOL_CONFIG.user} -p${POOL_CONFIG.password} -P ${POOL_CONFIG.port} ${POOL_CONFIG.database} < ${exec_file}`
        );
      }
      res.json({ token });
    } catch (e) {
      next(e);
    }
  });

  app.post<undefined, {} | ErrorResBody, Partial<BaseRequestData>>(
    "/room",
    async (req, res) => {
      const {
        schedule_id = "",
        meeting_room_id = "",
        start_at = 0,
        end_at = 0,
      } = req.body;
      if (!(await meetingRoomExist(meeting_room_id))) {
        res.status(400).json({
          message: "会議室が見つかりません",
        });
        return;
      }

      const startAt = new Date(start_at * 1000);
      const endAt = new Date(end_at * 1000);

      let connection: mysql.PoolConnection | undefined;
      try {
        connection = await getConnection(db);
        try {
          await beginTransaction(connection);

          const conflictSchedules = await query(connection)<MeetingRoom[]>(
            "SELECT * FROM meeting_room WHERE room_id = ? AND NOT (start_at >= ? OR end_at <= ?) FOR UPDATE",
            [meeting_room_id, endAt, startAt]
          );

          if (conflictSchedules.length > 0) {
            throw new ApplicationError(409, "すでに予定が入ってます");
          }

          await query(connection)(
            "INSERT INTO meeting_room (id, room_id, start_at, end_at) VALUES (?, ?, ?, ?)",
            [schedule_id, meeting_room_id, startAt, endAt]
          );

          await commit(connection);
        } catch (e) {
          if (connection) await rollback(connection);
          if (e instanceof ApplicationError) {
            res.status(e.statusCode).json({
              message: e.message,
            });
            return;
          }
          throw e;
        }
      } catch {
        res.status(500).json({
          message: "サーバー側のエラーです",
        });
        return;
      } finally {
        connection?.release();
      }

      res.status(201).json({});
    }
  );

  app.get<
    undefined,
    Pick<MeetingRoom, "room_id"> | ErrorResBody,
    undefined,
    Partial<{ readonly scheduleId?: string }>
  >("/room", async (req, res) => {
    const { scheduleId } = req.query;
    if (!scheduleId) {
      res.status(400).json({
        message: "スケジュールを指定してください",
      });
      return;
    }

    let connection: mysql.PoolConnection | undefined;
    try {
      connection = await getConnection(db);
      const [meetingRoom] = await query(connection)<MeetingRoom[]>(
        "SELECT * FROM meeting_room WHERE id = ?",
        [scheduleId]
      );

      if (meetingRoom === undefined) {
        res.status(404).json({
          message: "指定したスケジュールには部屋が予約されていません",
        });
        return;
      }

      res.status(200).json({
        room_id: meetingRoom.room_id,
      });
      return;
    } catch {
      res.status(500).json({
        message: "サーバー側のエラーです",
      });
    } finally {
      connection?.release();
    }
  });

  app.put<undefined, {} | ErrorResBody, Partial<BaseRequestData>>(
    "/room",
    async (req, res) => {
      const {
        schedule_id = "",
        meeting_room_id = "",
        start_at = 0,
        end_at = 0,
      } = req.body;

      let connection: mysql.PoolConnection | undefined;
      try {
        connection = await getConnection(db);

        if (meeting_room_id === "") {
          await query(connection)("DELETE FROM meeting_room WHERE id = ?", [
            schedule_id,
          ]);
          res.status(200).json({});
          return;
        }

        if (!(await meetingRoomExist(meeting_room_id))) {
          res.status(400).json({
            message: "会議室が見つかりません",
          });
          return;
        }

        try {
          await beginTransaction(connection);

          const startAt = new Date(start_at * 1000);
          const endAt = new Date(end_at * 1000);

          const conflictSchedules = await query(connection)<MeetingRoom[]>(
            "SELECT * FROM meeting_room WHERE room_id = ? AND NOT (start_at >= ? OR end_at <= ?) AND id != ? FOR UPDATE",
            [meeting_room_id, startAt, endAt, schedule_id]
          );

          if (conflictSchedules.length > 0) {
            throw new ApplicationError(409, "すでに予定が入ってます");
          }

          await query(connection)("DELETE FROM meeting_room WHERE id = ?", [
            schedule_id,
          ]);

          await query(connection)(
            "INSERT INTO meeting_room (id, room_id, start_at, end_at) VALUES (?, ?, ?, ?)",
            [schedule_id, meeting_room_id, startAt, endAt]
          );

          await commit(connection);
        } catch (e) {
          if (connection) await rollback(connection);
          if (e instanceof ApplicationError) {
            res.status(e.statusCode).json({
              message: e.message,
            });
            return;
          }
          throw e;
        }

        res.status(200).json({});
        return;
      } catch {
        res.status(500).json({
          message: "サーバー側のエラーです",
        });
        return;
      } finally {
        connection?.release();
      }
    }
  );

  app.listen(PORT, () => console.log(`Listening ${PORT}`));
}

main();
