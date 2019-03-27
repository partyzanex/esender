-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE `email` (
  `id`         BIGINT(20)                        NOT NULL AUTO_INCREMENT,
  `recipients` TEXT                              NOT NULL,
  `cc`         TEXT                              NOT NULL,
  `bcc`        TEXT                              NOT NULL,
  `sender`     VARCHAR(500)                      NOT NULL,
  `subject`    VARCHAR(500)                      NOT NULL,
  `mime_type`  ENUM ('html', 'text')             NOT NULL,
  `body`       LONGTEXT COLLATE utf8_unicode_ci  NOT NULL,
  `status`     ENUM ('created', 'sent', 'error') NOT NULL DEFAULT 'created',
  `error`      VARCHAR(1000)                              DEFAULT NULL,
  `dt_created` DATETIME                          NOT NULL DEFAULT now(),
  `dt_updated` DATETIME                                   DEFAULT NULL,
  `dt_sent`    DATETIME                                   DEFAULT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_unicode_ci;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE `email`;
