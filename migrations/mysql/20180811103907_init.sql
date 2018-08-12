-- +mig Up

CREATE TABLE `email` (
  `id`         INT(11)                           NOT NULL AUTO_INCREMENT,
  `to`         VARCHAR(255)                      NOT NULL,
  `from`       VARCHAR(255)                      NOT NULL,
  `title`      VARCHAR(255)                      NOT NULL,
  `mime_type`  VARCHAR(4)                        NOT NULL DEFAULT 'text',
  `text`       TEXT COLLATE utf8_unicode_ci      NOT NULL,
  `status`     ENUM ('created', 'sent', 'error') NOT NULL DEFAULT 'created',
  `error`      VARCHAR(512)                               DEFAULT NULL,
  `dt_created` INT(11)                           NOT NULL,
  `dt_updated` INT(11)                                    DEFAULT NULL,
  `dt_sent`    INT(11)                                    DEFAULT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_unicode_ci;

-- +mig Down

DROP TABLE `email`;

