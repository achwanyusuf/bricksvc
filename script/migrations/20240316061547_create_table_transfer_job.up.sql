CREATE SEQUENCE transfer_job_id_seq;
CREATE TYPE transferstatus AS ENUM ('failed', 'pending', 'success');

CREATE TABLE IF NOT EXISTS transfer_jobs (
  id integer primary key DEFAULT nextval('transfer_job_id_seq'),
  job_id varchar(30) NOT NULL,
  api_key text NOT NULL,
  payload jsonb NOT NULL,
  status transferstatus NOT NULL DEFAULT 'pending',
  created_by integer default 0 NOT NULL,
  created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_by integer default 0 NOT NULL,
  updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_by integer,
  deleted_at timestamp WITH TIME ZONE
);

ALTER SEQUENCE transfer_job_id_seq OWNED BY transfer_jobs.id;