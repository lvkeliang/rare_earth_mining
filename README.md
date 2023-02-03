# rare_earth_mining
This is the RedRock's winter vacation examine of the first semester.

# Database structure

## The database consists of the following tables:

- articleClassification

- articleTags

- articles

- collections

- comments

- events_list

- likes

- users

## Below is the structure of each table:

### 1.users(Used to store user data)

| Field        | Type          | Null | Key  | Default          | Extra |
| ------------ | ------------- | ---- | ---- | ---------------- | ----- |
| ID           | int           | NO   | PRI  | auto_increment   |       |
| mail         | varchar(50)   | NO   | UNI  |                  |       |
| nickname     | varchar(60)   | NO   | UNI  | new user         |       |
| password     | varchar(15)   | NO   |      | 123456           |       |
| position     | varchar(30)   | NO   |      | 还没写呢~      |       |
| company      | varchar(30)   | NO   |      | 还没写呢~      |       |
| introduction | varchar(300)  | NO   |      | 还没写呢~      |       |
| articleNum   | int           | NO   |      | 0                |       |
| collectNum   | int           | NO   |      | 0                |       |
| likeNum      | int           | NO   |      | 0                |       |
| commentNum   | int           | NO   |      | 0                |       |
| headPortrait | int           | NO   |      | 0                |       |

### 2.likes(Used to store like records)

| Field   | Type         | Null | Key  | Default                | Extra          |
|---------|--------------|------|------|------------------------|----------------|
| ID      | int          | NO   | PRI  |                        | auto_increment |
| uID     | int          | NO   |      |                        |                |
| oID     | varchar(13)  | NO   |      |                        |                |
| time    | timestamp    | NO   |      | CURRENT_TIMESTAMP      | DEFAULT_GENERATED |


### 3.events_list(Used to record events)

| Field      | Type         | Null | Key | Default               | Extra |
| ---------- | ------------ | ---- | --- | --------------------- | ----- |
| event_name | varchar(30)  | NO   |     |                      |       |
| event_started | timestamp | NO   |     |                      |       |

### 4.comments(Used to store comments)

| Field        | Type          | Null | Key  | Default           | Extra          |
| ------------ | ------------- | ---- | ---- | ----------------- | -------------- |
| ID           | int           | NO   | PRI  |                   | auto_increment |
| uID          | int           | NO   |      |                   |                |
| oID          | varchar(13)   | NO   |      |                   |                |
| publishTime  | timestamp     | NO   |      | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
| likeNum      | int           | NO   |      | 0                 |                |
| commentNum   | int           | NO   |      | 0                 |                |
| layer        | int           | NO   |      | 1                 |                |
| content      | varchar(300)  | NO   |      | 还没写呢~      |                |

### 5.collections(Used to store bookmark records)

| Field | Type | Null | Key | Default | Extra |
|-------|------|------|-----|---------|-------|
| ID    | int  | NO   | PRI |         | auto_increment |
| uID   | int  | NO   |     |         |         |
| oID   | varchar(13) | NO   |     |         |         |
| time  | timestamp | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |

### 6.articleTags(Used to store all article tags)

| Field    | Type        | Null | Key | Default           | Extra |
|----------|-------------|------|-----|-------------------|-------|
| ID       | int         | NO   | PRI | NULL              | auto_increment |
| tag      | varchar(20) | NO   |     | NULL              |        |

### 7.articles(Used to store article information)

| Field | Type          | Null | Key | Default                 | Extra |
|-------|---------------|------|-----|-------------------------|-------|
| ID    | int           | NO   | PRI |                        | auto_increment |
| uID   | int           | NO   |     |                        |       |
| title | varchar(100)  | NO   |     | 无标题                 |       |
| publishTime | timestamp | NO   |     | CURRENT_TIMESTAMP    | DEFAULT_GENERATED |
| updateTime  | timestamp | NO   |     | CURRENT_TIMESTAMP    | DEFAULT_GENERATED |
| viewerNum   | int       | NO   |     | 0                     |       |
| likeNum     | int       | NO   |     | 0                     |       |
| commentNum  | int       | NO   |     | 0                     |       |
| collectNum  | int       | NO   |     | 0                     |       |
| classification | varchar(20) | NO   |     | 未分类           |       |
| tags         | varchar(20) | NO   |     | 无                 |       |
| popularityValue | int | NO   |     | 0                     |       |
| state       | int       | NO   |     | 2                     |       |

The table has a trigger.

    create definer = ***@`%` trigger update_popularityValue
        before update
        on articles
        for each row
    BEGIN
        SET NEW.popularityValue = OLD.popularityValue + 1;
    end;

### 8.articleClassification(Used to store all article categories.)

| Field         | Type         | Null | Key  | Default | Extra          |
| ------------ | ------------ | ---- | ---- | ------- | ------------- |
| ID           | int          | NO   | PRI  |         | auto_increment|
| className    | varchar(10)  | NO   | UNI  |         |               |

## Below are the events of the database.

    create definer = ***@`%` event zero_popularityValue on schedule
        every '2' HOUR
            starts '2023-01-13 14:00:42'
        enable
        do
        BEGIN
            UPDATE articles SET popularityValue = 0;
            insert INTO events_list values('zero_popularityValue', now());
        end;

# API Readme

## Introduction

This API provides access to resources, such as data and services, for FE. Our API is designed to be RESTful and returns data in JSON format.

## Endpoints

The following endpoints are currently available:

### user

- POST `/user/register` : Sign up for an account
  - format : 
    - mail (must)
    - password (must)
    - nickname

- POST `/user/login` : Log in to your registered account
    - format :
        - mail (must)
        - password (must)

- GET `/user/information` : Gets the user's information
    - format :
        - uID (must)

- PUT `/user/profile` : Edit personal information
    - format :
        - null

### article

- GET `/article/brief` : Get brief information for multiple articles
    - format :
        - mode (must) : newest/popularity/publisher
        - pageNumber (must)
        - count (must)
        - firstaID (Not necessary, when mode is newest)
        - publisheruID (when mode is publisher)
        - classification
        - tags

- GET `/article/detail/{aID}` : Check out an article
    - format :
      - null

- POST `/article/postComment` : Post a comment to an article or comment
    - format :
        - oID (must)
        - content (must)

- POST `/article/like` : Like an article or comment
  - format :
      - oID (must)

- POST `/article/collect` : Bookmark an article
    - format :
        - oID (must)

- GET `/classification` : Get all article categories
    - format :
        - null

- GET `/tags` : Get all article tags
    - format :
        - null

### creator

- POST `/creator/publishArticle` : Publish an article
    - format :
        - title (must)
        - content (must)
        - classification (must)
        - tags (You should separate multiple tags with ",")

- POST `/creator/information` : View creator information for every day for 30 days
    - format :
      - null

- POST `/creator/myArticles` : View all your articles and their review status
  - format :
    - null

## Authentication
To access the API, you will need to include an API key as token in some request as a parameter to make sure you're signed in.

## Error Handling
If there is an error with a request, the API will return a non-10000 status code and an error message in JSON format.

## Response Format
The API returns data in JSON format.

## Examples

Get articles on the homepage by real-time popularity ranking:

    GET /articles/brief
    {
    "mode": "popularity",
    "pageNumber": 1,
    "count": 6
    "classification": "后端"
    "tags": "MySQL,知识"
    }

Publish an article:

    POST /creator/publishArticle
    {
    "title": "测试测试"
    "content": "{HTML code}"
    "classification": "后端"
    "tags": "MySQL,知识"
    }

## Conclusion

This API provides access to resources for FE. With its RESTful design and JSON data format, it is easy to use and integrate into your project. If you have any questions or concerns, please don't hesitate to contact us.