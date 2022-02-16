<h1 align="center">
  <abbr title="Information Schema Compatibility Table">I.S.Compat.Table</abbr>
</h1>

**I**nformation **S**chema **Compat**ibility **Table**(s)

Compares the information_schema of some of [the major databases that implement][implementors] the information_schema standard:

## About `information_schema`

`information_schema` **is an ANSI-standard** set of read-only views that provide information about all of the **tables, views, columns**, and **procedures** in a database.
Specifically, the specification for the information_schema views are published in [ISO/IEC 9075][iso-9075]. [This standard has several versions][version history].

This implies several crucial points:

1. The information_schema views are an amazing, standard way to discover metadata about a given database!
1. Different versions of the same database and different databases might implement different versions of the standard for information_schema
1. Since the standard is published by <abbr title="the International Standards Organization">ISO</abbr> reading it costs a nontrivial amount of money.
   Thus, volunteer developers _might_ choose to do something nice for themselves rather than shelling out so that they can implement the latest standards.

Naturally: most databases implement different views in `information_schema`.

<hr/>
<p align=center> <span style="font-family: monospace" title="flipping table"> (╯°□°）╯︵ ┻━┻</span></p>
<hr/>

## Motivation

I'd like to use `information_schema` more

## Inspirations

- [MDN's fantastic compatibility tables efforts](https://github.com/mdn/browser-compat-data)
- [Can I Use?](https://caniuse.com/ciu/about)
- [dbdb.io](https://dbdb.io)

<!-- https://simonwillison.net/2020/Oct/9/git-scraping/ -->

## Methodology

I scrape at the public documentation where the documentation licenses allow.
I also run databases without restrictive EULAs and observe those databases' `information_schema` tables directly.

I prioritize the most popular databases that implement an `information_schema` according to [2021 Stack Overflow Developer Survey](https://insights.stackoverflow.com/survey/2021#section-most-popular-technologies-databases)

| database name | % of respondents use | documentation scraped | `information_schema` queried directly |
| ------------- | -------------------: | --------------------- | ------------------------------------- |
| `postgres`    | 44%                  | YES                   | YES                                   |
| `mssql`       | 29%                  | YES                   | YES                                   |
| `mariadb`     | 17%                  | YES                   | YES                                   |
<!-- | `mysql`       | 48%            | NO                    | YES                           | -->
<!-- | `oracle`      | 13%            | NO                    | NO                            | -->
<!-- | `db2`         | 2%             | NO                    | YES                           | -->

  <!-- materialize: no information_schema _documented_; it should be postgres, right? -->
  <!-- dolt? -->
  <!-- ksqldb? -->
  <!-- Presto -->
  <!-- Apache Hive -->
  <!-- IBM db2? -->
  <!-- oracle whatever? -->
  <!-- your database here! -->

**If note a missing database that implements `information_schema` or evidence of incorrect information, please create a pull request with a fix!**

<!-- general links -->

[implementors]: https://en.wikipedia.org/wiki/Information_schema#Implementation
[iso-9075]: https://www.iso.org/standard/63555.html
[version history]: https://en.wikipedia.org/wiki/SQL#Standardization_history
