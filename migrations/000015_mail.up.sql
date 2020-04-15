BEGIN;
CREATE TABLE mail
(
    id           SERIAL  NOT NULL PRIMARY KEY,
    sender_id    INT     NOT NULL REFERENCES characters (id),
    recipient_id INT     NOT NULL REFERENCES characters (id),
    subject      VARCHAR NOT NULL DEFAULT '',
    body         VARCHAR NOT NULL DEFAULT '',
    read         BOOL    NOT NULL DEFAULT FALSE
);

CREATE INDEX mail_recipient_index ON mail (recipient_id);
END;