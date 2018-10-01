# DO NOT DELETE THIS FILE!

CREATE TABLE IF NOT EXISTS `accessStats` (
  `fullPath` text NOT NULL,
  `fileName` text NOT NULL,
  `accesses` bigint(20) NOT NULL DEFAULT '0',
  `lastAccess` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `requestLog` (
  `address` text NOT NULL,
  `userAgent` text NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `url` text NOT NULL,
  `code` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;