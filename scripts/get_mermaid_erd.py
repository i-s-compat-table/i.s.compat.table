import sqlite3
from pathlib import Path
from typing import Dict, List, NamedTuple, Optional

scripts_dir = Path(__file__).parent
repo_root = scripts_dir.parent
schema_def_file = Path(repo_root, "pkg/schema/db.sql")


class Column(NamedTuple):
    table_name: str
    column_name: str
    data_type: str
    nullable: bool


class Relation(NamedTuple):
    table_from: str
    column_from: str
    table_to: str
    column_to: str


class Table(NamedTuple):
    # name: str
    columns: List[Column]
    fks: List[Relation]  # outgoing only

    def get_col(self, column_name: str) -> Optional[Column]:
        for col in self.columns:
            if col.column_name == column_name:
                return col
        return None


# simulate a subset of information_schema.columns!
columns_sql = """
SELECT
    tbl.name AS table_name
  , col.name AS column_name
  , col.type AS data_type  
  , NOT col."notnull" AS nullable
FROM sqlite_master AS tbl
JOIN pragma_table_info(tbl.name) AS col
WHERE tbl.type = 'table';
"""

# from https://observablehq.com/@simonw/datasette-table-diagram-using-mermaid
relations_sql = """
select
  tbl.name as table_from,
  fk_info.[from] as column_from,
  fk_info.[table] as table_to,
  fk_info.[to] as column_to
from
  sqlite_master AS tbl
  join pragma_foreign_key_list(tbl.name) as fk_info
order by tbl.name
"""


def get_mermaid_erd() -> str:
    # create an empty in-memory copy of the schema
    conn = sqlite3.connect(":memory:")
    with open(schema_def_file) as f:
        schema_def = f.read()
    for statement in schema_def.split(";"):
        try:
            conn.execute(statement)
        except Exception as err:
            print(statement)
            raise err

    tables: Dict[str, Table] = {}  # map table name => table
    relations = [Relation(*row) for row in conn.execute(relations_sql)]
    columns = [Column(*row) for row in conn.execute(columns_sql)]

    for col in columns:
        table = tables.get(col.table_name, Table([], []))
        table.columns.append(col)
        tables[col.table_name] = table
    for rel in relations:
        table = tables.get(rel.table_from, Table([], []))
        table.fks.append(rel)
        tables[rel.table_from] = table

    erd = "erDiagram"
    for table_name in sorted(tables.keys()):
        table_from = tables[table_name]
        erd += "\n  " + table_name + " {"
        for col in table_from.columns:
            erd += f"\n    {col.data_type} {col.column_name}"

            if col.column_name == "id":
                erd += " PK"
            else:
                for fk in table_from.fks:
                    if col.column_name == fk.column_from:
                        erd += " FK"
            # this is where you'd add a comment to the column:
            # erd += f'"{comment}"'
        erd += "\n  }"

        for fk in table_from.fks:
            column_from = table_from.get_col(fk.column_from)
            table_to = tables[fk.table_to]
            column_to = table_to.get_col(fk.column_to)
            assert column_from is not None
            assert column_to is not None
            # https://mermaid-js.github.io/mermaid/#/entityRelationshipDiagram
            # assumes always fks always on id of one row in the foreign table
            from_cardinality = "}o" if column_from.nullable else "}|"
            to_cardinality = (
                "o|" if column_to.nullable and column_to.column_name != "id" else "||"
            )
            erd += f"\n  {table_name} {from_cardinality}--{to_cardinality} {fk.table_to} : {fk.column_from}"
    return erd


if __name__ == "__main__":
    print(get_mermaid_erd())
