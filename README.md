# health

A very basic web app for collecting personal health data and storing it
in a PostgreSQL database.

## Quick Start

* Create a database:

  ```bash
  $ createdb health
  ```

* Create tables:

  ```bash
  $ cat sql/users.sql | psql health
  $ cat sql/health_ref_variables.sql | psql health
  ```

* Insert some variables into `ref_variable`:

  ```sql
  INSERT INTO ref_variable (variable, description, units, sequence)
  VALUES ('temperature', 'Temperature', 'degrees Fahrenheit', 1),
         ('pulse', 'Pulse', 'beats per minute', 2),
         ('oxygen saturation', 'Oxygen saturation', 'percent', 3),
         ('body weight', 'Body weight, naked', 'pounds', 5);
  ```

* Update CORS `AllowedOrigins` on line 447 of `backend/health.go` with
  your domain.

* Build backend:

  ```bash
  $ cd backend
  $ go build health.go
  ```

* Configure `backend/health.env`, and copy to `backend/.env`

* Update `backend/health.service`, correcting the `User`,
  `WorkingDirectory` and `ExecStart` parameters to match your account
  and the location of the code.

* Start backend:

  ```bash
  $ cd /etc/systemd/system
  $ sudo ln -s ~/health/backend/health.service .
  $ sudo systemctl daemon-reload
  $ sudo systemctl enable health.service
  $ sudo systemctl start health.service
  ```

* Update `frontend/health/src/components/Data.vue`, correcting
  `LISTEN_PORT` or the server hostname on lines 84 and 85.

* Update `frontend/health/src/store/index.js`, correcting `LISTEN_PORT`
  or the server hostname on lines 10 and 11.

* Install and build (or `run serve`) the frontend:

  ```bash
  $ cd frontend/health
  $ npm install
  $ npm run build (or run serve)
  ```

## Acknowledgements

Much of the authentication code using
`github.com/antonlindstrom/pgstore` was coded by Sam Vanderwaal,
[ABR—Environmental Research & Services](https://www.abrinc.com/) and
used with permission.

[modeline]: # ( vim: set ft=markdown fenc=utf-8 tw=72 ts=2 sw=2 sts=2 spell spl=en: )
