USE `eveauth`;

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
-- Dumping data for table eveauth.accounts: ~6 rows (approximately)
/*!40000 ALTER TABLE `accounts` DISABLE KEYS */;
INSERT INTO `accounts` (`id`, `userid`, `apikeyid`, `apivcode`, `apiaccessmask`, `active`) VALUES
	(1, 1, 1, 'a', 0, 1),
	(2, 2, 2, 'b', 0, 0),
	(3, 3, 3, 'c', 0, 1),
	(4, 3, 4, 'd', 268435455, 1),
	(5, 4, 5, 'e', 268435455, 0),
	(6, 4, 6, 'f', 268435455, 0);
/*!40000 ALTER TABLE `accounts` ENABLE KEYS */;

-- Dumping data for table eveauth.applications: ~2 rows (approximately)
/*!40000 ALTER TABLE `applications` DISABLE KEYS */;
INSERT INTO `applications` (`id`, `name`, `maintainer`, `secret`, `callback`, `active`) VALUES
	(1, 'Testapp', 'Testmaintainer', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa', 'http://localhost/callback', 1),
	(2, 'Apptest', 'Testmaintainer', 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb', 'http://example.com/callback', 0);
/*!40000 ALTER TABLE `applications` ENABLE KEYS */;

-- Dumping data for table eveauth.characters: ~6 rows (approximately)
/*!40000 ALTER TABLE `characters` DISABLE KEYS */;
INSERT INTO `characters` (`id`, `accountid`, `corporationid`, `name`, `evecharacterid`, `defaultcharacter`, `active`) VALUES
	(1, 1, 1, 'Test Character', 1, 1, 1),
	(2, 2, 2, 'Please Ignore', 2, 1, 1),
	(3, 3, 1, 'Herp', 3, 1, 1),
	(4, 3, 1, 'Derp', 4, 0, 1),
	(5, 4, 2, 'Spai', 5, 0, 0),
	(6, 4, 2, 'NoSpai', 6, 1, 0);
/*!40000 ALTER TABLE `characters` ENABLE KEYS */;

-- Dumping data for table eveauth.corporations: ~2 rows (approximately)
/*!40000 ALTER TABLE `corporations` DISABLE KEYS */;
INSERT INTO `corporations` (`id`, `name`, `ticker`, `evecorporationid`, `ceoid`, `apikeyid`, `apivcode`, `active`) VALUES
	(1, 'Test Corp Please Ignore', 'TEST', 1, 1, 1, 'a', 1),
	(2, 'Corp Test Ignore Please', 'CORP', 2, 2, NULL, NULL, 0);
/*!40000 ALTER TABLE `corporations` ENABLE KEYS */;

-- Dumping data for table eveauth.grouproles: ~4 rows (approximately)
/*!40000 ALTER TABLE `grouproles` DISABLE KEYS */;
INSERT INTO `grouproles` (`id`, `groupid`, `roleid`, `autoadded`, `granted`) VALUES
	(1, 1, 1, 1, 1),
	(2, 1, 3, 0, 1),
	(3, 2, 2, 0, 0),
	(4, 2, 4, 1, 0);
/*!40000 ALTER TABLE `grouproles` ENABLE KEYS */;

-- Dumping data for table eveauth.groups: ~2 rows (approximately)
/*!40000 ALTER TABLE `groups` DISABLE KEYS */;
INSERT INTO `groups` (`id`, `name`, `active`) VALUES
	(1, 'Test Group', 1),
	(2, 'Dank Access', 0);
/*!40000 ALTER TABLE `groups` ENABLE KEYS */;

-- Dumping data for table eveauth.loginattempts: ~0 rows (approximately)
/*!40000 ALTER TABLE `loginattempts` DISABLE KEYS */;
/*!40000 ALTER TABLE `loginattempts` ENABLE KEYS */;

-- Dumping data for table eveauth.roles: ~4 rows (approximately)
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` (`id`, `name`, `active`) VALUES
	(1, 'ping.all', 1),
	(2, 'destroy.world', 0),
	(3, 'logistics.read', 1),
	(4, 'logistics.write', 1);
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;

-- Dumping data for table eveauth.usergroups: ~6 rows (approximately)
/*!40000 ALTER TABLE `usergroups` DISABLE KEYS */;
INSERT INTO `usergroups` (`id`, `userid`, `groupid`, `active`) VALUES
	(1, 1, 1, 1),
	(2, 2, 1, 0),
	(3, 3, 1, 1),
	(4, 3, 2, 1),
	(5, 4, 1, 0),
	(6, 4, 2, 0);
/*!40000 ALTER TABLE `usergroups` ENABLE KEYS */;

-- Dumping data for table eveauth.userroles: ~2 rows (approximately)
/*!40000 ALTER TABLE `userroles` DISABLE KEYS */;
INSERT INTO `userroles` (`id`, `userid`, `roleid`, `autoadded`, `granted`) VALUES
	(1, 1, 1, 0, 0),
	(2, 3, 2, 1, 1);
/*!40000 ALTER TABLE `userroles` ENABLE KEYS */;

-- Dumping data for table eveauth.users: ~4 rows (approximately)
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` (`id`, `username`, `password`, `email`, `verifiedemail`, `active`) VALUES
	(1, 'test1', '$2a$10$veif8VUZt7lShFhJKD0wGeY1YjCwIuWjYL0vQzlTqu8wNaYQMqzbe', 'test1@example.com', 1, 1),
	(2, 'test2', '$2a$10$95z.WXfIreLKJ9px.3KgpOq4aXTG3DF7/5ehGYzUWALhpN6MMq/aK', 'test2@example.com', 0, 0),
	(3, 'test3', '$2a$10$7Yxm2scdTVpEJpvZAT7tbOFA.G9JfyxtiHbr989iocX6U37C3/j4q', 'test3@example.com', 0, 1),
	(4, 'test4', '$2a$10$WOWTgqaqLKbkb1uhYbtLnOuuYX4kXBC61GVAke7RkjiODoBpgGGzy', 'test4@example.com', 1, 0);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
