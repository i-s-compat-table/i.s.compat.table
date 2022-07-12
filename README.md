<h1 align="center">
  <abbr title="Information Schema Compatibility Table">I.S.Compat.Table</abbr>
</h1>
<p align=center><b>I</b>nformation <b>S</b>chema <b>Compat</b>ibility <b>Table</b>(s)</p>

Compares the information_schema of some of [the major databases that implement][implementors] the information_schema standard.

## About `information_schema`

`information_schema` is an ANSI-standard set of read-only views that provide information about all of the tables, views, columns, and procedures in a database.
Specifically, the specification for the information_schema views are published in [ISO/IEC 9075][iso-9075]. [This standard has several versions][version history].

These facts implies several crucial points:

1. The information_schema views are an amazing, standard way to discover metadata about a given database!
1. Different versions of the same database and different databases might implement different versions of the standard for information_schema
1. Since the standard is published by <abbr title="the International Standards Organization">ISO</abbr> reading it costs a nontrivial amount of money.
   Thus, volunteer developers _might_ choose to do something nice for themselves rather than shelling out so that they can implement the latest standards.

Naturally, most databases that implement `information_schema` a subset of the standard's views, add extra database-specific views, and stuff otherwise-standard views with database-specific columns.
This makes `information_schema` a highly-nonstandard standard. Thinking of `information_schema` as a convention might be more accurate.

<hr/>
<p align=center> <span style="font-family: monospace" title="flipping a table (pun intended)"> (╯°□°）╯︵ ┻━┻</span></p>
<hr/>

## Motivation

I'd like to use `information_schema` more. Before I do that, however, I'd like to know what views and columns are in the standard or better yet, what views and columns are actually in each database's `information_schema`.

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
| ------------- | :------------------: | :-------------------: | :-----------------------------------: |
| `mysql`       |         48%          |          NO           |                  YES                  |
| `postgres`    |         44%          |          YES          |                  YES                  |
| `mssql`       |         29%          |          YES          |                  NO                   |
| `mariadb`     |         17%          |          YES          |                  YES                  |
| `cockroachdb` |                      |          YES          |                  TODO                 |

 <!-- |`oracle` |                  13% | NO                    | NO                                    | -->
 <!-- |`db2`    |                   2% | NO                    | NO                                    | -->

  <!-- `tidb` -->
  <!-- `presto` -->
  <!-- `materializedb`: no information_schema _documented_; it should be postgres, right? -->
  <!-- dolt? -->
  <!-- ksqldb? -->

  <!-- Apache Hive -->
  <!-- your database here! -->

  <!-- commercial databases -->
  <!-- `snowflakedb`? -->
  <!-- `db2`? -->
  <!-- `oracle` via oracle cloud's free tier? -->

**If note a missing database that implements `information_schema` or evidence of incorrect information, please create a pull request with a fix!**
See [`./CONTRIBUTING.md`](./CONTRIBUTING.md) for more details.

<!-- general links -->

[implementors]: https://en.wikipedia.org/wiki/Information_schema#Implementation
[iso-9075]: https://www.iso.org/standard/63555.html
[version history]: https://en.wikipedia.org/wiki/SQL#Standardization_history
