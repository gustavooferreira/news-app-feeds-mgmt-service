# News App Exercise - Feeds Management Service

Please check out documentation for the whole system [here](https://github.com/gustavooferreira/news-app-docs).

This repository structure follows [this convention](https://github.com/golang-standards/project-layout).

---

# Build

To build a binary, run:

```
make build
```

The `api-server` binary will be placed inside the `bin/` folder.

To build the docker image, run:

```
make build-docker
```

The docker image is called `news-app/feeds-api-server`.

---

# Tests

To run tests:

```
make test
```

To get coverage:

```
make coverage
```

## Tip

> If you run `make` without any targets, it will display all options available on the makefile followed by a short description.

# Design

Discuss any design decision here ...

## MySQL tables

Table `feeds`:

| url                                               | provider | category | enabled |
| ------------------------------------------------- | :------: | :------: | :-----: |
| http://feeds.bbci.co.uk/news/uk/rss.xml           |    1     |    1     |  true   |
| http://feeds.bbci.co.uk/news/technology/rss.xml   |    1     |    2     |  true   |
| http://feeds.skynews.com/feeds/rss/uk.xml         |    2     |    1     |  true   |
| http://feeds.skynews.com/feeds/rss/technology.xml |    2     |    2     |  true   |

url is a Primary key.

Table `providers`:

| id  | name     |
| :-: | -------- |
|  1  | BBC News |
|  2  | Sky News |

id is a Primary key.

Table `categories`:

| id  | name       |
| :-: | ---------- |
|  1  | UK         |
|  2  | Technology |

id is a Primary key.
