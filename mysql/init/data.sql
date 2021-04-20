
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
SET character_set_client = utf8mb4 ;
CREATE TABLE `user` (
`user_id` int(11) NOT NULL AUTO_INCREMENT,
`user_name` varchar(225) NOT NULL,
`avatar` varchar(225),
`is_delete` tinyint(1) NOT NULL DEFAULT '1',
`create_at` timestamp  default current_timestamp,
`update_at` timestamp default current_timestamp on update current_timestamp,
PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `auth`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
SET character_set_client = utf8mb4 ;
CREATE TABLE `auth` (
`auth_id` int(11) NOT NULL AUTO_INCREMENT,
`user_id` int(11) NOT NULL,
`login_id` varchar(225) NOT NULL UNIQUE,
`password` varchar(223) NOT NULL,
`mail_address` varchar(223) NOT NULL,
`create_at` timestamp  default current_timestamp,
`update_at` timestamp default current_timestamp on update current_timestamp,
FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE,
PRIMARY KEY (`auth_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `project`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
SET character_set_client = utf8mb4 ;
CREATE TABLE `project` (
`project_id` int(11) NOT NULL AUTO_INCREMENT,
`project_name` varchar(225) NOT NULL,
`description` varchar(225),
`project_avatar` varchar(225),
`create_at` timestamp  default current_timestamp,
`update_at` timestamp default current_timestamp on update current_timestamp,
PRIMARY KEY (`project_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `user_project`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
SET character_set_client = utf8mb4 ;
CREATE TABLE `user_project` (
`user_project_id` int(11) NOT NULL AUTO_INCREMENT,
`user_id` int(11) NOT NULL,
`project_id` int(11) NOT NULL,
`create_at` timestamp  default current_timestamp,
`update_at` timestamp default current_timestamp on update current_timestamp,
PRIMARY KEY (`user_project_id`),
FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE,
FOREIGN KEY (`project_id`) REFERENCES `project` (`project_id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
SET character_set_client = utf8mb4 ;
CREATE TABLE `status` (
`status_id` int(11) NOT NULL AUTO_INCREMENT,
`project_id` int(11) NOT NULL,
`progress` int(11) NOT NULL,
`status_name` varchar(225) NOT NULL,
`create_at` timestamp  default current_timestamp,
`update_at` timestamp default current_timestamp on update current_timestamp,
FOREIGN KEY (`project_id`) REFERENCES `project` (`project_id`) ON DELETE CASCADE,
PRIMARY KEY (`status_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `ticket`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
SET character_set_client = utf8mb4 ;
CREATE TABLE `ticket` (
`ticket_id` int(11) NOT NULL AUTO_INCREMENT,
`project_id` int(11) NOT NULL,
`title` varchar(225) NOT NULL,
`explanation` varchar(225),
`reporter` int(11),
`worker` int(11),
`create_at` timestamp  default current_timestamp,
`update_at` timestamp default current_timestamp on update current_timestamp,
FOREIGN KEY (`project_id`) REFERENCES `project` (`project_id`) ON DELETE CASCADE,
PRIMARY KEY (`ticket_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `ticket_status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
SET character_set_client = utf8mb4 ;
CREATE TABLE `ticket_status` (
`ticket_status_id` int(11) NOT NULL AUTO_INCREMENT,
`ticket_id` int(11) NOT NULL,
`status_id` int(11) NOT NULL,
`create_at` timestamp  default current_timestamp,
`update_at` timestamp default current_timestamp on update current_timestamp,
FOREIGN KEY (`ticket_id`) REFERENCES `ticket` (`ticket_id`) ON DELETE CASCADE,
FOREIGN KEY (`status_id`) REFERENCES `status` (`status_id`) ON DELETE CASCADE,
PRIMARY KEY (`ticket_status_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `ticket_img`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
SET character_set_client = utf8mb4 ;
CREATE TABLE `ticket_img` (
`ticket_img_id` int(11) NOT NULL AUTO_INCREMENT,
`ticket_id` int(11) NOT NULL,
`ticket_img_path` varchar(225) NOT NULL,
`create_at` timestamp  default current_timestamp,
`update_at` timestamp default current_timestamp on update current_timestamp,
FOREIGN KEY (`ticket_id`) REFERENCES `ticket` (`ticket_id`) ON DELETE CASCADE,
PRIMARY KEY (`ticket_img_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `comment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
SET character_set_client = utf8mb4 ;
CREATE TABLE `comment` (
`comment_id` int(11) NOT NULL AUTO_INCREMENT,
`user_id` int(11) NOT NULL,
`ticket_id` int(11) NOT NULL,
`comment` varchar(225) NOT NULL,
`create_at` timestamp  default current_timestamp,
`update_at` timestamp default current_timestamp on update current_timestamp,
FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE,
FOREIGN KEY (`ticket_id`) REFERENCES `ticket` (`ticket_id`) ON DELETE CASCADE,
PRIMARY KEY (`comment_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `comment_img`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
SET character_set_client = utf8mb4 ;
CREATE TABLE `comment_img` (
`comment_img_id` int(11) NOT NULL AUTO_INCREMENT,
`comment_id` int(11) NOT NULL ,
`comment_img_path` varchar(225) NOT NULL,
`create_at` timestamp  default current_timestamp,
`update_at` timestamp default current_timestamp on update current_timestamp,
FOREIGN KEY (`comment_id`) REFERENCES `comment` (`comment_id`) ON DELETE CASCADE,
PRIMARY KEY (`comment_img_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `sample`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
SET character_set_client = utf8mb4 ;
CREATE TABLE `sample` (
`sample_id` int(11) NOT NULL AUTO_INCREMENT,
`numa` int(11) NOT NULL ,
`numb` int(11) NOT NULL,
`numc` int(11) NOT NULL ,
`create_at` timestamp  default current_timestamp,
`update_at` timestamp default current_timestamp on update current_timestamp,
PRIMARY KEY (`sample_id`),
unique num_unique (numa, numb, numc)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

insert into user (user_name) values ('test');
insert into auth (user_id, login_id, password, mail_address) values
((select max(user_id) from user), 'test', '$2a$10$.QOu/aJH9U45RbfQyaySOeVScRmDJixLwtf1xiuBzwmZvpGX2W6qO', 'reporter@reporter.com');

insert into project(project_name, description, project_avatar) values
('project1', 'タスク管理アプリケーションAPI', 'http://test1'),
('project2', 'タスク管理アプリケーションFrontEnd', 'http://test2');

insert into status (project_id, progress, status_name) values
((select project_id from project where project_name = 'project1'), 1, 'TODO'),
((select project_id from project where project_name = 'project1'), 2, 'GOING'),
((select project_id from project where project_name = 'project1'), 3, 'DONE'),
((select project_id from project where project_name = 'project2'), 1, '未着手'),
((select project_id from project where project_name = 'project2'), 2, '着手'),
((select project_id from project where project_name = 'project2'), 3, '完了');

insert into user_project(user_id, project_id) values
((select max(user_id) from user), (select project_id from project where project_name = 'project1')),
((select max(user_id) from user), (select project_id from project where project_name = 'project2'));

insert into user (user_name, avatar) values ('reporter', 'user/02.jpg');
insert into auth (user_id, login_id, password, mail_address) values
((select max(user_id) from user), 'reporter', '2cb22086791fca9cc06c2239fd74f2f9', 'reporter@reporter.com');

insert into user_project(user_id, project_id) values
((select max(user_id) from user), (select project_id from project where project_name = 'project1')),
((select max(user_id) from user), (select project_id from project where project_name = 'project2'));

insert into user (user_name, avatar) values ('worker', 'user/03.jpg');
insert into auth (user_id, login_id, password, mail_address) values
((select max(user_id) from user), 'worker', 'c8d97880379aa8ba1efeacbcced85c0e', 'reporter@reporter.com');
insert into user_project(user_id, project_id) values
((select max(user_id) from user), (select project_id from project where project_name = 'project1')),
((select max(user_id) from user), (select project_id from project where project_name = 'project2'));
insert into ticket (project_id, title, explanation, reporter, worker) values
((select project_id from project where project_name = 'project1'), 'bug修正', 'testBody1', (select user_id from user where user_name = 'reporter'), (select user_id from user where user_name = 'worker')),
((select project_id from project where project_name = 'project1'), 'storage URL変更', 'testBody2', (select user_id from user where user_name = 'reporter'), (select user_id from user where user_name = 'worker')),
((select project_id from project where project_name = 'project1'), 'log機能実装', 'testBody3', (select user_id from user where user_name = 'reporter'), (select user_id from user where user_name = 'worker')),
((select project_id from project where project_name = 'project2'), 'header改修', 'testBody4', (select user_id from user where user_name = 'reporter'), (select user_id from user where user_name = 'worker')),
((select project_id from project where project_name = 'project2'), '環境変数導入', 'testBody5', (select user_id from user where user_name = 'reporter'), (select user_id from user where user_name = 'worker')),
((select project_id from project where project_name = 'project2'), 'プロジェクト一覧page改修', 'testBody6', (select user_id from user where user_name = 'reporter'), (select user_id from user where user_name = 'worker'));

insert into ticket_status (ticket_id, status_id) values
((select ticket_id from ticket where title = 'bug修正'), (select status_id from status where status_name = 'TODO')),
((select ticket_id from ticket where title = 'storage URL変更'), (select status_id from status where status_name = 'GOING')),
((select ticket_id from ticket where title = 'log機能実装'), (select status_id from status where status_name = 'DONE')),
((select ticket_id from ticket where title = 'header改修'), (select status_id from status where status_name = '未着手')),
((select ticket_id from ticket where title = '環境変数導入'), (select status_id from status where status_name = '着手')),
((select ticket_id from ticket where title = 'プロジェクト一覧page改修'), (select status_id from status where status_name = '完了'));

insert into user (user_name, avatar) values ('yasui', 'user/01.jpg');
insert into auth (user_id, login_id, password, mail_address) values
((select max(user_id) from user), 'yasui', '$2a$10$.r0EzJ4ahKKDl2fe9/UwkOZxMI2jXlmdtUWCdtOv7q4XkSrYJEMcy', 'sample@gmail.com');

insert into user_project(user_id, project_id) values
((select max(user_id) from user), (select project_id from project where project_name = 'project1')),
((select max(user_id) from user), (select project_id from project where project_name = 'project2'));