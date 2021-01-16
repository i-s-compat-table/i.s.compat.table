<h1>
    <abbr title="Information Schema Compatibility Table">i.s.compat.table</abbr>
</h1>

> Information Schema Compatibility Table(s) 

Compares the information_schema of some of [the major databases that implement][implementors] the information_schema standard:
- MySQL
- Mariadb
- Postgres
- CockroachDB
- Microsoft SQL Server
<!-- Presto -->
<!-- Apache Hive -->
<!-- IBM db2? -->
<!-- oracle whatever? -->
<!-- your database here! -->

## Methodology
I looked at the public docs of each of these projects to check if they documented each of the information_schema tables.
I added all the information_schema tables with at least 2-3 implementors to the comparison table.  

### Caveats
- Some links may be broken or mistaken (e.g. a mysql link pointing to the mariadb docs due to a typo in the html).
- My search might not have been complete.
- I interpreted absence of  evidence of a table in a database's docs as evidence of absence of that table from the database.

**If you find evidence of incorrect information, please create both a PR to fix it!**


## The standard

According to wikipedia, the information schema (`information_schema`) **is an ANSI-standard** set of read-only views that provide information about all of the **tables, views, columns**, and **procedures** in a database.
Specifically, the specification for the information_schema views are published in ISO/IEC 9075, section 4. [This standard has several versions][version history].

This implies several crucial points:
1. The information_schema views are an amazing, standard way to discover metadata about a given database!
1. Different versions of the same database and different databases might implement different versions of the standard for information_schema 
1. Since the standard is published by <abbr title="the International Standards Organization">ISO</abbr> reading it costs a nontrivial amount of money.
Thus, volunteer developers _might_ choose to do something nice for themselves rather than shelling out so that they can implement the latest standards.

Natually: most databases implement different views in `information_schema`.

<pre> (╯°□°）╯︵ ┻━┻</pre>

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

|                                       |MySQL       |Mariadb         |PostgreSQL|cockroachdb  |SQL Server|
|:-------------------------------------:|:----------:|:--------------:|:--------:|:-----------:|:--------:|
|`applicable_roles`                     |[>=8.0][001]|      [all][002]|[all][003]|      X      |    X     |
|`character_sets`                       |  [all][004]|      [all][005]|[all][006]|      X      |    X     |
|`check_constraints`                    |[>=8.0][007]|[>=10.2.22][008]|[all][009]|[>=19.2][010]|    X     |
|`collation_character_set_applicability`|  [all][011]|      [all][012]|[all][013]|      X      |    X     |
|`collations`                           |  [all][014]|      [all][015]|[all][016]|      X      |    X     |
|`column_domain_usage`                  |      X     |       X        |[all][017]|      X      |[all][018]|
|`column_privileges`                    |  [all][019]|      [all][020]|[all][021]|      X      |[all][022]|
|`columns`                              |  [all][023]|      [all][024]|[all][025]|   [all][026]|[all][027]|
|`constraint_column_usage`              |      X     |       X        |[all][028]|      X      |[all][029]|
|`constraint_table_usage`               |      X     |       X        |[all][030]|      X      |[all][031]|
|`domain_constraints`                   |      X     |       X        |[all][032]|      X      |[all][033]|
|`domains`                              |      X     |       X        |[all][034]|      X      |[all][035]|
|`enabled_roles`                        |[>=8.0][036]| [>=10.0.5][037]|[all][038]|      X      |    X     |
|`engines`                              |  [all][039]|      [all][040]|    X     |      X      |    X     |
|`events`                               |  [all][041]|      [all][042]|    X     |      X      |    X     |
|`key_column_usage`                     |  [all][043]|      [all][044]|[all][045]|   [all][046]|[all][047]|
|`parameters`                           |  [all][048]|    [>=5.5][049]|[all][050]|   [  X][051]|[all][052]|
|`partitions`                           |  [all][053]|      [all][054]|[all][055]|      X      |    X     |
|`referential_contraints`               |  [all][056]|      [all][057]|[all][058]|   [all][059]|[all][060]|
|`role_column_grants`                   |[>=8.0][061]|       X        |[all][062]|      X      |    X     |
|`role_routine_grants`                  |  [all][063]|       X        |[all][064]|      X      |    X     |
|`role_table_grants`                    |[>=8.0][065]|       X        |[all][067]|   [all][068]|    X     |
|`schema_privileges`                    |  [all][068]|      [all][069]|    X     |      X      |    X     |
|`schemata`                             |  [all][070]|      [all][071]|[all][072]|   [all][073]|[all][074]|
|`sequences`                            |      X     |       X        |[all][075]| [>=2.0][076]|    X     |
|`statistics`                           |  [all][077]|      [all][078]|    X     |      X      |    X     |
|`table_constraints`                    |  [all][079]|      [all][080]|[all][081]|   [all][082]|[all][083]|
|`table_privileges`                     |  [all][084]|      [all][085]|[all][086]|   [all][087]|[all][088]|
|`triggers`                             |  [all][089]|      [all][090]|[all][091]|      X      |    X     |
|`user_privileges`                      |  [all][092]|      [all][093]|    X     |      X      |    X     |
|`views`                                |  [all][093]|      [all][094]|[all][095]|   [all][096]|[all][097]|
|`view_routine_usage`                   |  [all][098]|       X        |[all][099]|      X      |    X     |
|`view_table_usage`                     |[>=8.0][100]|       X        |[all][101]|      X      |[all][102]|

<!-- notes: I rounded the Cockroachdb version ranges to the minor version -->
<!--  -->

## Prior art
MDN's fantastic compatibility tables efforts!

<!-- general links -->
[implementors]: https://en.wikipedia.org/wiki/Information_schema#Implementation
[iso/iec-9075]: https://www.iso.org/standard/63555.html
[version history]: https://en.wikipedia.org/wiki/SQL#Standardization_history

<!-- reference links -->
<!-- TODO: group by database? -->
[001]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-applicable-roles-table.html
[002]: https://mariadb.com/kb/en/information-schema-applicable_roles-table/
[003]: https://www.postgresql.org/docs/13/infoschema-applicable-roles.html
[004]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-character-sets-table.html
[005]: https://mariadb.com/kb/en/information-schema-character_sets-table/
[006]: https://www.postgresql.org/docs/13/infoschema-character-sets.html
[007]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-check-constraints-table.html
[008]: https://mariadb.com/kb/en/information-schema-check_constraints-table/
[009]: https://www.postgresql.org/docs/13/infoschema-check-constraints.html
[010]: https://www.cockroachlabs.com/docs/stable/information-schema.html#check_constraints
[011]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-collation-character-set-applicability-table.html
[012]: https://mariadb.com/kb/en/information-schema-collation_character_set_applicability-table/
[013]: https://www.postgresql.org/docs/13/infoschema-collation-character-set-applicab.html
[014]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-collations-table.html
[015]: https://mariadb.com/kb/en/information-schema-collations-table/
[016]: https://www.postgresql.org/docs/13/infoschema-collations.html
[017]: https://www.postgresql.org/docs/13/infoschema-column-domain-usage.html
[018]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/column-domain-usage-transact-sql?view=sql-server-ver15
[019]: https://dev.mysql.com/doc/refman/5.7/en/information-schema-column-privileges-table.html
[020]: https://mariadb.com/kb/en/information-schema-column_privileges-table/
[021]: https://www.postgresql.org/docs/13/infoschema-column-privileges.html
[022]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/column-privileges-transact-sql?view=sql-server-ver15
[023]: https://dev.mysql.com/doc/refman/5.7/en/information-schema-columns-table.html
[024]: https://mariadb.com/kb/en/information-schema-columns-table/
[025]: https://www.postgresql.org/docs/13/infoschema-columns.html
[026]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/columns-transact-sql?view=sql-server-ver15
[027]: https://www.cockroachlabs.com/docs/stable/information-schema.html#columns
[028]: https://www.postgresql.org/docs/13/infoschema-constraint-column-usage.html
[029]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/constraint-column-usage-transact-sql?view=sql-server-ver15
[030]: https://www.postgresql.org/docs/13/infoschema-constraint-table-usage.html
[031]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/constraint-table-usage-transact-sql?view=sql-server-ver15
[032]: https://www.postgresql.org/docs/13/infoschema-domain-constraints.html
[033]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/domain-constraints-transact-sql?view=sql-server-ver15
[034]: https://www.postgresql.org/docs/13/infoschema-domains.html
[035]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/domains-transact-sql?view=sql-server-2017
[036]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-enabled-roles-table.html
[037]: https://mariadb.com/kb/en/information-schema-enabled_roles-table/
[038]: https://www.postgresql.org/docs/13/infoschema-enabled-roles.html
[039]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-engines-table.html
[040]: https://mariadb.com/kb/en/information-schema-engines-table/
[041]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-events-table.html
[042]: https://mariadb.com/kb/en/information-schema-events-table/
[043]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-key-column-usage-table.html
[044]: https://mariadb.com/kb/en/information-schema-key_column_usage-table/
[045]: https://www.postgresql.org/docs/13/infoschema-key-column-usage.html
[046]: https://www.cockroachlabs.com/docs/stable/information-schema.html#key_column_usage
[047]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/key-column-usage-transact-sql?view=sql-server-ver15
[048]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-parameters-table.html
[049]: https://mariadb.com/kb/en/information-schema-parameters-table/
[050]: https://www.postgresql.org/docs/13/infoschema-parameters.html
[051]: https://www.cockroachlabs.com/docs/stable/information-schema.html
[052]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/parameters-transact-sql?view=sql-server-ver15
[053]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-partitions-table.html
[054]: https://mariadb.com/kb/en/information-schema-partitions-table/
[055]: https://www.postgresql.org/docs/13/infoschema-parameters.html
[056]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-referential-constraints-table.html
[057]: https://mariadb.com/kb/en/information-schema-referential_constraints-table/ 
[058]: https://www.postgresql.org/docs/13/infoschema-referential-constraints.html
[059]: https://www.cockroachlabs.com/docs/stable/information-schema.html#referential_constraints
[060]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/referential-constraints-transact-sql?view=sql-server-ver15
[061]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-role-column-grants-table.html
[062]: https://www.postgresql.org/docs/13/infoschema-role-column-grants.html
[063]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-role-routine-grants-table.html
[064]: https://www.postgresql.org/docs/13/infoschema-role-routine-grants.html
[065]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-role-table-grants-table.html
[066]: https://www.postgresql.org/docs/13/infoschema-role-table-grants.html
[067]: https://www.cockroachlabs.com/docs/stable/information-schema.html#role_table_grants
[068]: https://dev.mysql.com/doc/refman/5.6/en/information-schema-schema-privileges-table.html
[069]: https://mariadb.com/kb/en/information-schema-schema_privileges-table/
[070]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-routines-table.html
[071]: https://mariadb.com/kb/en/information-schema-routines-table/
[072]: https://www.postgresql.org/docs/13/infoschema-routines.html
[073]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/routines-transact-sql?view=sql-server-ver15
[074]: https://www.cockroachlabs.com/docs/stable/information-schema.html
[075]: https://www.postgresql.org/docs/13/infoschema-sequences.html
[076]: https://www.cockroachlabs.com/docs/stable/information-schema.html#sequences
[077]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-statistics-table.html
[078]: https://mariadb.com/kb/en/information-schema-statistics-table/
[079]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-table-constraints-table.html
[080]: https://mariadb.com/kb/en/information-schema-table_constraints-table/
[081]: https://www.postgresql.org/docs/13/infoschema-table-constraints.html
[082]: https://www.cockroachlabs.com/docs/stable/information-schema.html#table_constraints
[083]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/table-constraints-transact-sql?view=sql-server-ver15
[084]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-table-privileges-table.html
[085]: https://mariadb.com/kb/en/information-schema-table_privileges-table/
[086]: https://www.postgresql.org/docs/13/infoschema-table-privileges.html
[087]: https://www.cockroachlabs.com/docs/stable/information-schema.html#table_privileges
[088]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/table-privileges-transact-sql?view=sql-server-ver15
[089]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-triggers-table.html
[090]: https://mariadb.com/kb/en/information-schema-triggers-table/
[091]: https://www.postgresql.org/docs/13/infoschema-triggers.html
[092]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-user-privileges-table.html
[093]: https://mariadb.com/kb/en/information-schema-user_privileges-table/
[093]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-views-table.html
[094]: https://mariadb.com/kb/en/information-schema-views-table/
[095]: https://www.postgresql.org/docs/13/infoschema-views.html
[096]: https://www.cockroachlabs.com/docs/stable/information-schema.html#views
[097]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/views-transact-sql?view=sql-server-ver15
[098]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-view-routine-usage-table.html
[099]: https://www.postgresql.org/docs/13/infoschema-view-routine-usage.html
[100]: https://dev.mysql.com/doc/refman/8.0/en/information-schema-view-table-usage-table.html
[101]: https://www.postgresql.org/docs/13/infoschema-view-table-usage.html
[102]: https://docs.microsoft.com/en-us/sql/relational-databases/system-information-schema-views/view-table-usage-transact-sql?view=sql-server-ver15