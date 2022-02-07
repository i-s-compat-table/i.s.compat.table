<h1>
    <abbr title="Information Schema Compatibility Table">I.S.Compat.Table</abbr>
</h1>

**I**nformation **S**chema **Compat**ibility **Table**(s)

Compares the information_schema of some of [the major databases that implement][implementors] the information_schema standard:

<!-- https://insights.stackoverflow.com/survey/2021#section-most-popular-technologies-databases -->
- MySQL
- Mariadb
- Postgres
- CockroachDB
- Microsoft SQL Server
  <!-- materialize: no information_schema _documented_; it should be postgres, right? -->
  <!-- dolt? -->
  <!-- ksqldb? -->
  <!-- Presto -->
  <!-- Apache Hive -->
  <!-- IBM db2? -->
  <!-- oracle whatever? -->
  <!-- your database here! -->

## Methodology

I looked at the public docs of each of these projects to check if they documented each of the information_schema tables.
I went as far as scraping the documentation sites that had permissive licenses (postgres, MariaDB, MSSQL, CockroachDB).
I added all the information_schema tables with at least 2-3 implementors to the comparison table.

### Caveats

- Some links may be broken or mistaken (e.g. a mysql link pointing to the mariadb docs due to a typo in the html).
- My search might not have been complete.
- I interpreted absence of evidence of a table in a database's docs as evidence of absence of that table from the database.

**If you find evidence of incorrect information, please create a PR to fix it!**

## The standard

According to wikipedia, the information schema (`information_schema`) **is an ANSI-standard** set of read-only views that provide information about all of the **tables, views, columns**, and **procedures** in a database.
Specifically, the specification for the information_schema views are published in [ISO/IEC 9075][iso/iec-9057]. [This standard has several versions][version history].

This implies several crucial points:

1. The information_schema views are an amazing, standard way to discover metadata about a given database!
1. Different versions of the same database and different databases might implement different versions of the standard for information_schema
1. Since the standard is published by <abbr title="the International Standards Organization">ISO</abbr> reading it costs a nontrivial amount of money.
   Thus, volunteer developers _might_ choose to do something nice for themselves rather than shelling out so that they can implement the latest standards.

Naturally: most databases implement different views in `information_schema`.

<hr/>
<p align=center> <span style="font-family: monospace" title="flipping table"> (╯°□°）╯︵ ┻━┻</span></p>
<hr/>

<!-- ## Table of contents (spoiler: the contents are tables)
0. tables
1. `applicable_roles`
1. `character_sets`
1. `check_constraints`
1. `collation_character_set_applicability`
1. `collations`
1. `column_domain_usage`
1. `column_privileges`
1. `columns`
1. `constraint_column_usage`
1. `constraint_table_usage`
1. `domain_constraints`
1. `domains`
1. `enabled_roles`
1. `engines`
1. `events`
1. `key_column_data_store`
1. `parameters`
1. `partitions`
1. `referential_contraints`
1. `role_column_grants`
1. `role_routine_grants`
1. `role_table_grants`
1. `schema_privileges`
1. `schemata`
1. `sequences`
1. `spatial_ref_sys`
1. `statistics`
1. `table_constraints`
1. `table_privileges`
1. `triggers`
1. `user_privileges`
1. `views`
1. `view_routine_usage`
1. `view_table_usage` -->

## Tables

<!-- this should be an output of the underlying data -->
|                                         |    MySQL     |     Mariadb      | PostgreSQL |  cockroachdb  | SQL Server | [snowflake][sf] |
| :-------------------------------------: | :----------: | :--------------: | :--------: | :-----------: | :--------: | :-------------: |
|           `applicable_roles`            |  [>=8.0][my01]|       [all][ma01]| [all][pg01] |       X       |     X      |       all       |
|            `character_sets`             |    [all][my02]|       [all][ma02]| [all][pg02] |       X       |     X      |        X        |
|           `check_constraints`           |  [>=8.0][my03]| [>=10.2.22][ma03]| [all][pg03] | [>=19.2][010] |     X      |        X        |
| `collation_character_set_applicability` |    [all][my04]|       [all][ma04]| [all][pg04] |       X       |     X      |        X        |
|              `collations`               |    [all][my05]|       [all][ma05]| [all][pg05] |       X       |     X      |        X        |
|          `column_domain_usage`          |      X       |        X         | [all][pg06] |       X       | [all][018] |        X        |
|           `column_privileges`           |  [all][my07]  |       [all][ma07]| [all][pg07] |       X       | [all][022] |        X        |
|                `columns`                |  [all][my08]  |       [all][ma08]| [all][pg08] |  [all][026]   | [all][027] |       all       |
|        `constraint_column_usage`        |      X       |        X         | [all][pg09] |       X       | [all][029] |        X        |
|        `constraint_table_usage`         |      X       |        X         | [all][pg10] |       X       | [all][031] |        X        |
|          `domain_constraints`           |      X       |        X         | [all][pg11] |       X       | [all][033] |        X        |
|                `domains`                |      X       |        X         | [all][pg12] |       X       | [all][035] |        X        |
|             `enabled_roles`             | [>=8.0][036] |  [>=10.0.5][ma13]| [all][pg13] |       X       |     X      |       all       |
|                `engines`                |  [all][039]  |       [all][ma14]|     X      |       X       |     X      |        X        |
|                `events`                 |  [all][041]  |       [all][ma15]|     X      |       X       |     X      |        X        |
|           `key_column_usage`            |  [all][043]  |       [all][ma16]| [all][pg16] |  [all][046]   | [all][047] |        X        |
|              `parameters`               |  [all][048]  |     [>=5.5][ma17]| [all][pg17] |   [ X][051]   | [all][052] |        X        |
|              `partitions`               |  [all][053]  |       [all][ma18]| [all][pg18] |       X       |     X      |        X        |
|        `referential_contraints`         |  [all][056]  |       [all][ma19]| [all][pg19] |  [all][059]   | [all][060] |       all       |
|          `role_column_grants`           | [>=8.0][061] |        X         | [all][pg20] |       X       |     X      |        X        |
|          `role_routine_grants`          |  [all][063]  |        X         | [all][pg21] |       X       |     X      |        X        |
|           `role_table_grants`           | [>=8.0][065] |        X         | [all][pg22] |  [all][068]   |     X      |        X        |
|               `routines`                |       [all][]       |                  |            |               |            |
|           `schema_privileges`           |  [all][068]  |    [all][ma24]    |     X      |       X       |     X      |        X        |
|               `schemata`                |  [all][070]  |    [all][ma25]    | [all][pg25] |  [all][073]   | [all][074] |       all       |
|               `sequences`               |      X       |        X         | [all][pg26] | [>=2.0][076]  |     X      |       all       |
|              `statistics`               |  [all][077]  |    [all][ma27]    |     X      |       X       |     X      |        X        |
|           `table_constraints`           |  [all][079]  |    [all][ma28]    | [all][pg28] |  [all][082]   | [all][083] |       all       |
|           `table_privileges`            |  [all][084]  |    [all][ma29]    | [all][pg29] |  [all][087]   | [all][088] |       all       |
|                `tables`                 |              |                  |            |               |            |       all       |
|               `triggers`                |  [all][089]  |    [all][ma31]    | [all][pg31] |       X       |     X      |        X        |
|            `user_privileges`            |  [all][092]  |    [all][ma32]    |     X      |       X       |     X      |        X        |
|                 `views`                 |  [all][094]  |    [all][ma33]    | [all][pg33] |  [all][097]   | [all][098] |       all       |
|          `view_routine_usage`           |  [all][099]  |        X         | [all][pg34] |       X       |     X      |        X        |
|           `view_table_usage`            | [>=8.0][101] |        X         | [all][pg35] |       X       | [all][103] |        X        |

<!-- notes: I rounded the Cockroachdb version ranges to the minor version -->
<!-- from snowflake allegations: USAGE_PRIVILEGES? -->

https://www.postgresql.org/docs/13/infoschema-routines.html
## Prior art

MDN's fantastic compatibility tables efforts!

<!-- general links -->

[implementors]: https://en.wikipedia.org/wiki/Information_schema#Implementation
[iso/iec-9075]: https://www.iso.org/standard/63555.html
[version history]: https://en.wikipedia.org/wiki/SQL#Standardization_history

<!-- reference links -->
<!-- TODO: group by database? -->

[pg01]: https://www.postgresql.org/docs/13/infoschema-applicable-roles.html
[pg02]: https://www.postgresql.org/docs/13/infoschema-character-sets.html
[pg03]: https://www.postgresql.org/docs/13/infoschema-check-constraints.html
[pg04]: https://www.postgresql.org/docs/13/infoschema-collation-character-set-applicab.html
[pg05]: https://www.postgresql.org/docs/13/infoschema-collations.html
[pg06]: https://www.postgresql.org/docs/13/infoschema-column-domain-usage.html
[pg07]: https://www.postgresql.org/docs/13/infoschema-column-privileges.html
[pg08]: https://www.postgresql.org/docs/13/infoschema-columns.html
[pg09]: https://www.postgresql.org/docs/13/infoschema-constraint-column-usage.html
[pg10]: https://www.postgresql.org/docs/13/infoschema-constraint-table-usage.html
[pg16]: https://www.postgresql.org/docs/13/infoschema-key-column-usage.html
[pg17]: https://www.postgresql.org/docs/13/infoschema-parameters.html
[pg18]: https://www.postgresql.org/docs/13/infoschema-parameters.html
[pg19]: https://www.postgresql.org/docs/13/infoschema-referential-constraints.html
[pg20]: https://www.postgresql.org/docs/13/infoschema-role-column-grants.html
[pg35]: https://www.postgresql.org/docs/13/infoschema-view-table-usage.html
[pg34]: https://www.postgresql.org/docs/13/infoschema-view-routine-usage.html
[pg31]: https://www.postgresql.org/docs/13/infoschema-triggers.html
[pg33]: https://www.postgresql.org/docs/13/infoschema-views.html
[pg25]: https://www.postgresql.org/docs/13/infoschema-routines.html
[pg26]: https://www.postgresql.org/docs/13/infoschema-sequences.html
[pg21]: https://www.postgresql.org/docs/13/infoschema-role-routine-grants.html
[pg22]: https://www.cockroachlabs.com/docs/stable/information-schema.html#role_table_grants
[pg28]: https://www.postgresql.org/docs/13/infoschema-table-constraints.html
[pg29]: https://www.postgresql.org/docs/13/infoschema-table-privileges.html
[066]: https://www.postgresql.org/docs/13/infoschema-role-table-grants.html

[my01]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-applicable-roles-table.html
[my02]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-character-sets-table.html
[my03]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-check-constraints-table.html
[010]: https://www.cockroachlabs.com/docs/stable/information-schema.html#check_constraints
[my04]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-collation-character-set-applicability-table.html

[ma01]: https://mariadb.com/kb/en/information-schema-applicable_roles-table/
[ma02]: https://mariadb.com/kb/en/information-schema-character_sets-table/
[ma03]: https://mariadb.com/kb/en/information-schema-check_constraints-table/
[ma04]: https://mariadb.com/kb/en/information-schema-collation_character_set_applicability-table/
[ma05]: https://mariadb.com/kb/en/information-schema-collations-table/
[ma07]: https://mariadb.com/kb/en/information-schema-column_privileges-table/
[ma08]: https://mariadb.com/kb/en/information-schema-columns-table/
[ma13]: https://mariadb.com/kb/en/information-schema-enabled_roles-table/
[ma14]: https://mariadb.com/kb/en/information-schema-engines-table/
[ma15]: https://mariadb.com/kb/en/information-schema-events-table/
[ma16]: https://mariadb.com/kb/en/information-schema-key_column_usage-table/
[ma17]: https://mariadb.com/kb/en/information-schema-parameters-table/
[ma18]: https://mariadb.com/kb/en/information-schema-partitions-table/
[ma19]: https://mariadb.com/kb/en/information-schema-referential_constraints-table/


[my05]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-collations-table.html
[my07]: https://dev.mysql.com/doc/refman/5.7/en/information-schema-column-privileges-table.html
[my08]: https://dev.mysql.com/doc/refman/5.7/en/information-schema-columns-table.html
[036]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-enabled-roles-table.html
[039]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-engines-table.html
[041]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-events-table.html
[043]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-key-column-usage-table.html
[048]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-parameters-table.html
[053]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-partitions-table.html


[018]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/column-domain-usage-transact-sql?view=sql-server-ver15
[022]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/column-privileges-transact-sql?view=sql-server-ver15
[026]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/columns-transact-sql?view=sql-server-ver15
[029]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/constraint-column-usage-transact-sql?view=sql-server-ver15
[031]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/constraint-table-usage-transact-sql?view=sql-server-ver15
[033]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/domain-constraints-transact-sql?view=sql-server-ver15
[035]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/domains-transact-sql?view=sql-server-2pg6
[047]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/key-column-usage-transact-sql?view=sql-server-ver15
[052]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/parameters-transact-sql?view=sql-server-ver15
[060]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/referential-constraints-transact-sql?view=sql-server-ver15

[027]: https://www.cockroachlabs.com/docs/stable/information-schema.html#columns
[046]: https://www.cockroachlabs.com/docs/stable/information-schema.html#key_column_usage
[051]: https://www.cockroachlabs.com/docs/stable/information-schema.html
[059]: https://www.cockroachlabs.com/docs/stable/information-schema.html#referential_constraints

[056]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-referential-constraints-table.html
[061]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-role-column-grants-table.html
[063]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-role-routine-grants-table.html
[065]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-role-table-grants-table.html
[068]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-schema-privileges-table.html
[070]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-routines-table.html

[ma24]: https://mariadb.com/kb/en/information-schema-schema_privileges-table/
[ma25]: https://mariadb.com/kb/en/information-schema-routines-table/
[ma27]: https://mariadb.com/kb/en/information-schema-statistics-table/
[ma28]: https://mariadb.com/kb/en/information-schema-table_constraints-table/
[ma29]: https://mariadb.com/kb/en/information-schema-table_privileges-table/
[ma31]: https://mariadb.com/kb/en/information-schema-triggers-table/
[ma32]: https://mariadb.com/kb/en/information-schema-user_privileges-table/

[073]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/routines-transact-sql?view=sql-server-ver15
[074]: https://www.cockroachlabs.com/docs/stable/information-schema.html
[076]: https://www.cockroachlabs.com/docs/stable/information-schema.html#sequences
[077]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-statistics-table.html
[079]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-table-constraints-table.html
[082]: https://www.cockroachlabs.com/docs/stable/information-schema.html#table_constraints
[083]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/table-constraints-transact-sql?view=sql-server-ver15
[084]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-table-privileges-table.html
[087]: https://www.cockroachlabs.com/docs/stable/information-schema.html#table_privileges
[088]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/table-privileges-transact-sql?view=sql-server-ver15
[089]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-triggers-table.html
[092]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-user-privileges-table.html
[094]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-views-table.html
[099]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-view-routine-usage-table.html
[ma33]: https://mariadb.com/kb/en/information-schema-views-table/
[097]: https://www.cockroachlabs.com/docs/stable/information-schema.html#views
[098]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/views-transact-sql?view=sql-server-ver15
[101]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-view-table-usage-table.html
[103]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/view-table-usage-transact-sql?view=sql-server-ver15
[sf]: https://docs.snowflake.com/en/sql-reference/info-schema.html
