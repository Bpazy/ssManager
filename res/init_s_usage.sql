CREATE TABLE s_usage
(
  port      int      NOT NULL,
  date      datetime NOT NULL,
  downUsage bigint DEFAULT 0 NOT NULL,
  upUsage   bigint DEFAULT 0 NOT NULL
);
