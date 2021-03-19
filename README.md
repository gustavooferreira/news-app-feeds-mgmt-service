# News App Exercise - Feeds Management Service

Please check out documentation for the whole system [here](https://github.com/gustavooferreira/news-app-docs).

This repository structure follows [this convention](https://github.com/golang-standards/project-layout).

---

# Build

To build this project run:

```
make build
```

The `api-server` binary will be placed inside the `bin/` folder.

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

## Free tip

> If you run `make` without any targets, it will display all options available on the makefile followed by a short description.

# Design

## MySQL tables

Table `feeds`:

| url                                               | provider | category | enabled |
| ------------------------------------------------- | :------: | :------: | :-----: |
| http://feeds.bbci.co.uk/news/uk/rss.xml           |    1     |    1     |  true   |
| http://feeds.bbci.co.uk/news/technology/rss.xml   |    1     |    2     |  true   |
| http://feeds.skynews.com/feeds/rss/uk.xml         |    2     |    1     |  true   |
| http://feeds.skynews.com/feeds/rss/technology.xml |    2     |    2     |  true   |

Table `providers`:

| id  | name     |
| :-: | -------- |
|  1  | BBC News |
|  2  | Sky News |

Table `categories`:

| id  | name       |
| :-: | ---------- |
|  1  | UK         |
|  2  | Technology |
